package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/takumakei/go-enzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	enzap.ReplaceGlobals()
	defer zap.L().Sync()

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"debug", "info", "warn", "error", "fatal"}
	}

	for _, arg := range args {
		var level zapcore.Level
		if err := level.Set(arg); err != nil {
			zap.L().Error("level.Set", zap.String("arg", arg), zap.Error(err))
			continue
		}
		zap.L().Debug("level.Set", zap.String("arg", arg))

		switch level {
		case zapcore.DebugLevel:
			zap.L().Debug("debug")
		case zapcore.InfoLevel:
			zap.L().Info("info")
		case zapcore.WarnLevel:
			zap.L().Warn("warn")
		case zapcore.ErrorLevel:
			zap.L().Error("error")
		case zapcore.DPanicLevel:
			zap.L().DPanic("dpanic")
		case zapcore.PanicLevel:
			zap.L().Panic("panic")
		case zapcore.FatalLevel:
			zap.L().Fatal("fatal")
		}
	}
}

func usage() {
	fmt.Printf("usage: %s [<level>...]\n", os.Args[0])
}
