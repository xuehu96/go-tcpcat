package server

import (
	"fmt"
	"go.uber.org/zap"
	"net"
)

type Options func(s *Server)

func WithHook(h Hook) Options {
	return func(srv *Server) {
		if h.OnListen == nil {
			h.OnListen = defaultOnListen
		}
		if h.OnAccept == nil {
			h.OnAccept = defaultOnAccept
		}
		if h.OnReadData == nil {
			h.OnReadData = defaultOnReadData
		}
		if h.OnFnCode == nil {
			h.OnFnCode = defaultGetFnCode
		}
		if h.OnSendData == nil {
			h.OnSendData = defaultOnSendData
		}
		if h.OnClose == nil {
			h.OnClose = defaultOnClose
		}
		if h.OnStop == nil {
			h.OnStop = defaultOnStop
		}
		srv.hooks = h
	}
}

func WithIPandPort(ip string, port int) Options {
	return func(s *Server) {
		l, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
		s.ln = l
	}
}
func WithPort(port int) Options {
	return func(s *Server) {
		l, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
		s.ln = l
	}
}

func WithListener(l net.Listener) Options {
	return func(s *Server) {
		s.ln = l
	}
}

func WithLogger(logger *zap.Logger) Options {
	return func(s *Server) {
		log = logger
	}
}
