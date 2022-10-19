package main

import (
	"context"
	"github.com/xuehu96/go-tcpcat/pkg/logger"
	"github.com/xuehu96/go-tcpcat/server"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// 空白TCP服务器，可以在此基础上增加功能函数和狗子(Hook)函数

func main() {
	// 创建TCPListener
	ln, err := net.Listen("tcp", ":9677")
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	// 创建TCP服务实例
	s := server.New(
		server.WithListener(ln),
		server.WithLogger(logger.DebugLogger()),
	)

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
