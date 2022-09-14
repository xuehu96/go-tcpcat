package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func DebugLogger() *zap.Logger {
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
	log := zap.New(core)
	return log
}

func ProductionLogger() *zap.Logger {
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
	log := zap.New(core)
	return log
}
