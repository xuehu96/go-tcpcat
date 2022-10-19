package server

import (
	"context"
	"github.com/pkg/errors"
	"github.com/xuehu96/go-tcpcat/pkg/snowflake"
	"go.uber.org/zap"
	"net"
	"sync"
)

var log *zap.Logger

// Serve TCP服务器开始干活
func (s *Server) Serve() {
	defer log.Sync()
	snowflake.Init("2022-07-07", 1)

	go s.hooks.OnListen(s)

	// 主线程accept
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Info("go-tcpcat server stop listening")
				return
			}
			zap.L().Error("Accept failed", zap.Error(err))
			continue
		}
		c := NewClient(s, conn, snowflake.GetId())
		go c.Start()
	}
}

func (s *Server) Stop(ctx context.Context) {
	// 1. 循环关client
	for _, c := range s.clients {
		c.Close()
	}

	//close(s.exitChan)

	// 最后关监听，监听在主线程，程序就会退出
	s.ln.Close()

	// hook onStop
	s.hooks.OnStop()
}

// AddFn 添加功能函数
func (s *Server) AddFn(code string, f CodeFnType) {
	s.fns[code] = f
}

// New 返回server实例
func New(opts ...Options) *Server {
	s := defaultServer()
	log = zap.NewNop()
	for _, fn := range opts {
		fn(s)
	}
	return s
}

func defaultServer() *Server {
	s := &Server{
		bufLen:  1024,
		ln:      nil,
		clients: make(map[string]*Client),
		hooks: Hook{
			OnListen:   defaultOnListen,
			OnAccept:   defaultOnAccept,
			OnReadData: defaultOnReadData,
			OnFnCode:   defaultGetFnCode,
			OnSendData: defaultOnSendData,
			OnClose:    defaultOnClose,
			OnStop:     defaultOnStop,
		},
		fns:      make(map[string]CodeFnType),
		mu:       sync.RWMutex{},
		exitChan: make(chan struct{}),
	}
	return s
}
