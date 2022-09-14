package server

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"os"
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

func WithDebugLogger() Options {
	return func(s *Server) {
		// 1. writeSyncer循环文件
		lumberJackLogger := &lumberjack.Logger{
			Filename:   "debug.log",
			LocalTime:  true,
			MaxSize:    16, //M
			MaxBackups: 0,  //个
			MaxAge:     1,  //天
			Compress:   false,
		}
		writeSyncer := zapcore.AddSync(lumberJackLogger)

		// 5. 自定义格式
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder := zapcore.NewJSONEncoder(encoderConfig)

		// 6.core
		var core zapcore.Core
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zap.DebugLevel),
		)

		// 7.log
		log = zap.New(core)
	}
}

func WithProductionLogger() Options {
	return func(s *Server) {
		// 1. writeSyncer循环文件
		os.Mkdir("log", os.ModePerm)
		lumberJackLogger := &lumberjack.Logger{
			Filename:   "./log/xgtcp.log",
			LocalTime:  true,
			MaxSize:    128, //M
			MaxBackups: 10,  //个
			MaxAge:     365, //天
			Compress:   false,
		}
		writeSyncer := zapcore.AddSync(lumberJackLogger)

		// 5. 自定义格式
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder := zapcore.NewJSONEncoder(encoderConfig)

		// 6.core
		var core zapcore.Core
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
		// 7.log
		log = zap.New(core)
	}
}
