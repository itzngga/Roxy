package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.Logger
var cmdLogger *zap.Logger

func NewLogger(lvlstr string) *zap.Logger {
	if logger != nil {
		return logger
	}
	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(pe)
	consoleEncoder := zapcore.NewConsoleEncoder(pe)
	if lvlstr == "ingfo" {
		lvlstr = "info"
	}
	level, err := zapcore.ParseLevel(lvlstr)

	if err != nil {
		panic(err)
	}

	logFile, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	l := zap.New(core)
	logger = l
	return logger
}

func CmdLogger(lvlstr string) *zap.Logger {
	if cmdLogger != nil {
		return cmdLogger
	}
	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(pe)
	consoleEncoder := zapcore.NewConsoleEncoder(pe)
	if lvlstr == "ingfo" {
		lvlstr = "info"
	}
	level, err := zapcore.ParseLevel(lvlstr)

	if err != nil {
		panic(err)
	}

	logFile, _ := os.OpenFile("cmd.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	l := zap.New(core)
	cmdLogger = l
	return cmdLogger
}
