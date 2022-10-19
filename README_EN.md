# go-tcpcat

[中文](./README.md)

TCP is often applied to IoT data server or game servers, and is usually a private format protocol;  
Go tcpcat can realize TCP transparent transmission and message processing;  
Message processing supports`Golang`Hook to inject private protocol parsers into the framework,`Lua `script processing,`Python`GRPC call processing,`HTTP`callback processing,`Redis`cache data, etc.

### 1.Hook Function

|Hook | Call time | Usage example|
| ---------- | :----------------------- |------------------------------------|
|OnListen | Callback after successful Listen | Program starts New resource|
|OnAccept | New client connection call | Determine the source of the client and limit the number of clients|
|OnReadData | Received client data call | Judge the data format, remove the header and footer, etc|
|OnFNCode | Parse function codes from data | [*Function codes and user-defined function handlers*](# 2. User-defined data processing function) |
|Onsenddata | call after sending data to the client | judge whether sending is successful, resend or notify the caller|
|Onclose | called when the client is active or disconnected | notify or mark the client offline after disconnection|
|Onstop | called after listening is closed | clean up resources|

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


### 2. User-defined data processing function

In private protocols, packet processing methods are generally called ** function codes **, similar to HTTP routing


Function Code Parsing Function 'OnFNCode' (Function code is parsed from data)
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
Functions for processing corresponding function codes
```go
s.AddFn("A", func(c *server.Client, code string, buf []byte, len int) {
	// TODO
})
s.AddFn("B", func(c *server.Client, code string, buf []byte, len int) {
    // TODO
})
```


## 3. Ping-Pong TCPServer example

After the client connects, it sends data to the server. The server processes the data in the following format
- "Ping" replies "pong"
- "Time" replies to the current server time
- The "exit" server is actively disconnected
- The client sends no reply in other formats


```go
package main

import (
	"context"
	"github.com/xuehu96/go-tcpcat/pkg/logger"
	"github.com/xuehu96/go-tcpcat/server"
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
		server.WithLogger(logger.DebugLogger()),
	)

	// 添加功能码对应的处理函数 类似于HTTP的路由
	s.AddFn("p", func(c *server.Client, code string, data []byte) {
		c.ReplyData([]byte("pong"))
	})
	s.AddFn("t", func(c *server.Client, code string, data []byte) {
		currentTime := time.Now()
		c.ReplyData([]byte(currentTime.Format("2006-01-02 15:04:05.000000000")))
	})
	s.AddFn("x", func(c *server.Client, code string, data []byte) {
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
