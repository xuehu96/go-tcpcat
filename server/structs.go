package server

import (
	"context"
	"net"
	"sync"
)

// Server TCP服务器实例, 用New()创建
type Server struct {
	bufLen int
	// socket
	ln      net.Listener
	clients map[string]*Client

	// 狗子
	hooks Hook

	// 功能表
	fns map[string]CodeFnType

	mu       sync.RWMutex
	exitChan chan struct{}
}

// Client TCP服务器Accept的客户端
type Client struct {
	id     int64
	key    string
	closed bool
	buf    []byte
	s      *Server
	conn   net.Conn
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}
