package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const bufferSize = 1024 * 32

// Options logger options
type Options struct {
	// log save filename
	Filename string
	// MB
	MaxSize int
	// number of files
	MaxBackups int
	// day
	MaxDays int
	// if true output code line
	Caller bool
	// file log level [debug info warn error dpanic panic fatal]
	FileLevel string
	// console file level [debug info warn error dpanic panic fatal]
	ConsoleLevel string
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// Helper zap log helper
type Helper struct {
	noCopy noCopy
	// zap
	*zap.Logger

	fileLevel    zap.AtomicLevel
	consoleLevel zap.AtomicLevel
}

var emptyAtomicLevel = zap.NewAtomicLevel()

// Attach Logger
func (l *Helper) Attach(src *Helper) {
	l.Logger = src.Logger
	l.fileLevel = src.fileLevel
	l.consoleLevel = src.consoleLevel
}

// Detach Logger
func (l *Helper) Detach() {
	l.Logger = nil
	l.fileLevel = emptyAtomicLevel
	l.consoleLevel = emptyAtomicLevel
}

// FileLevel return file level
func (l *Helper) FileLevel() zap.AtomicLevel {
	return l.fileLevel
}

// ConsoleLevel return console level
func (l *Helper) ConsoleLevel() zap.AtomicLevel {
	return l.consoleLevel
}

// NewHelper create a logger
func NewHelper(options *Options, zapOptions ...zap.Option) *Helper {
	var cores []zapcore.Core
	fileLevel := zap.NewAtomicLevel()
	consoleLevel := zap.NewAtomicLevel()
	// file
	fileLevel = zap.NewAtomicLevel()
	if options.FileLevel == "" {
		fileLevel.SetLevel(zap.FatalLevel)
	} else if e := fileLevel.UnmarshalText([]byte(options.FileLevel)); e != nil {
		fileLevel.SetLevel(zap.FatalLevel)
	}
	cores = append(cores, zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		newLoggerCache(zapcore.AddSync(&lumberjack.Logger{
			Filename:   options.Filename,
			MaxSize:    options.MaxSize, // megabytes
			MaxBackups: options.MaxBackups,
			MaxAge:     options.MaxDays, // days
		}), bufferSize),
		fileLevel,
	))

	// console
	if options.ConsoleLevel == "" {
		consoleLevel.SetLevel(zap.FatalLevel)
	} else if e := consoleLevel.UnmarshalText([]byte(options.ConsoleLevel)); e != nil {
		consoleLevel.SetLevel(zap.FatalLevel)
	}
	cores = append(cores, zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		newLoggerCache(monitor, bufferSize),
		consoleLevel,
	))

	if options.Caller {
		zapOptions = append(zapOptions, zap.AddCaller())
	}

	return &Helper{
		Logger: zap.New(
			zapcore.NewTee(cores...),
			zapOptions...,
		),
		fileLevel:    fileLevel,
		consoleLevel: consoleLevel,
	}
}
