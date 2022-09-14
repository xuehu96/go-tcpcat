package main

import (
	"context"
	"github.com/xuehu96/xgtcp/pkg/logger"
	"github.com/xuehu96/xgtcp/server"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// ping-pong TCP服务器example
// 客户端连接后 向服务器发送数据，服务器按以下格式处理
// "ping" 回复 "pong"
// "time" 回复当前服务器时间
// "exit" 服务器主动断开
// 客户端发送其他格式 不回复

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
