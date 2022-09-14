# xgtcp 

TCP服务器常用于物联网、游戏服务器；  
HTTP协议有URL，MQTT协议有topic，而TCP协议通常使用私有协议，xgtcp为开发框架，通过hooks的方式，将私有协议的解析程序注入框架。


### 1.狗子函数(Hooks)
狗子函数就是回调函数，指定事件发生后，框架会回调提前注入的处理函数，
默认的处理函数是通过日志打印，设置狗子函数后，会覆盖默认的处理函数

程序设置了以下狗子函数：

| Hook       | 调用时机                 | 用途示例                                                     |
| ---------- | :----------------------- | ------------------------------------------------------------ |
| OnListen   | Listen成功后回调         | 程序启动 New资源                                             |
| OnAccept   | 新客户端连接调用         | 判断客户端来源，限制客户端数量                               |
| OnReadData | 接收到客户端数据调用     | 判断数据格式，去头尾等                                       |
| OnFnCode   | 从数据中解析出功能码     | [*功能码和自定义功能处理函数*](#2.功能码和自定义功能处理函数) |
| OnSendData | 向客户端发送数据后调用   | 判断是否发送成功，重发或通知调用方                           |
| OnClose    | 客户端主动或被断开后调用 | 断开后通知或标记客户端离线                                   |
| OnStop     | Listen关闭后调用         | 清理资源                                                     |

狗子使用方法：
```go
var fnc1 server.OnXXXXXX = func(param) ret{ 
	// TODO
}
var fnc2 server.OnYYYYYY = func(param) ret{
	// TODO
}
         ...

hooks := server.Hook{
	OnXXXXXX: fnc,
	OnYYYYYY: fnc,
	...
}
s := server.New(
    server.WithHook(hooks),
	...
)
```

### 2.功能码和自定义功能处理函数
私有协议中区分数据包处理方式的一般叫**功能码**，类似于HTTP的路由

功能码解析函数`OnFnCode`（从data中解析出功能码）
```go
var fnc server.OnFnCode = func(buf []byte) string {
	s := string(buf)
	if strings.Contains(s, "ping") {
		return "A"
	}
	if s[0] == 'B' {
		return "B"
	}
	return ""
}

hooks := server.Hook{
	OnFnCode: fnc,
}
```
处理对应功能码的函数（类似于路由）
```go
s.AddFn("A", func(c *server.Client, code string, buf []byte, len int) {
	// TODO
})
s.AddFn("B", func(c *server.Client, code string, buf []byte, len int) {
    // TODO
})
```



## 3. Ping-Pong TCP服务器example

客户端连接后 向服务器发送数据，服务器按以下格式处理
- "ping" 回复 "pong"
- "time" 回复当前服务器时间
- "exit" 服务器主动断开
-  客户端发送其他格式 不回复

```go
package main

import (
	"context"
	"github.com/xuehu96/xgtcp/server"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	// 创建TCPListener
	ln, err := net.Listen("tcp", ":9677")
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	// 自定义如何从数据包中获取获取功能码
	var fnc server.OnFnCode = func(buf []byte) string {
		s := string(buf)
		if strings.Contains(s, "ping") {
			return "p"
		}
		if strings.Contains(s, "time") {
			return "t"
		}
		if strings.Contains(s, "exit") {
			return "x"
		}
		return ""
	}

	hooks := server.Hook{
		OnFnCode: fnc,
	}
	// 创建TCP服务实例
	s := server.New(
		server.WithHook(hooks),
		server.WithListener(ln),
		server.WithDebugLogger(),
	)

	// 添加功能码对应的处理函数 类似于HTTP的路由
	s.AddFn("p", func(c *server.Client, code string, buf []byte, len int) {
		c.ReplyData([]byte("pong"))
	})
	s.AddFn("t", func(c *server.Client, code string, buf []byte, len int) {
		currentTime := time.Now()
		c.ReplyData([]byte(currentTime.Format("2006-01-02 15:04:05.000000000")))
	})
	s.AddFn("x", func(c *server.Client, code string, buf []byte, len int) {
		c.Close()
	})

	// Ctrl-C 结束
	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
		<-signalCh
		s.Stop(context.Background())
	}()

	// TCP服务器开始干活
	s.Serve()
}

```

