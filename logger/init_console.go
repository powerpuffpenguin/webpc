package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitConsole(level string) {
	var zapOptions []zap.Option
	lv := zap.NewAtomicLevel()
	if level == `` {
		lv.SetLevel(zap.InfoLevel)
	} else if e := lv.UnmarshalText([]byte(level)); e != nil {
		lv.SetLevel(zap.InfoLevel)
	}
	Logger.Attach(&Helper{
		Logger: zap.New(
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
				os.Stdout,
				lv,
			),
			zapOptions...,
		),
		fileLevel:    lv,
		consoleLevel: lv,
	})
}
