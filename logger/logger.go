package logger

import (
	"fmt"
	"sync"

	toml "github.com/pelletier/go-toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Zap *zap.SugaredLogger
}

var (
	ZapLogger *Logger
	once      sync.Once
)

func NewLogger(toml *toml.Tree) *Logger {

	var (
		logLevel  = toml.Get("logger.level").(string)
		logPath   = toml.Get("logger.path").(string)
		logPrefix = toml.Get("logger.prefix").(string)
	)

	once.Do(func() {
		var zapLogLevel zapcore.Level
		switch logLevel {
		case "debug":
			zapLogLevel = zapcore.DebugLevel
		case "info":
			zapLogLevel = zapcore.InfoLevel
		case "warn":
			zapLogLevel = zapcore.WarnLevel
		case "error":
			zapLogLevel = zapcore.ErrorLevel
		case "panic":
			zapLogLevel = zapcore.PanicLevel
		case "fatal":
			zapLogLevel = zapcore.FatalLevel
		}

		l, err := zap.Config{
			Development:      false,
			Encoding:         "console",
			OutputPaths:      []string{"stdout", logPath + logPrefix},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:   "Message",
				TimeKey:      "Time",
				EncodeTime:   zapcore.RFC3339TimeEncoder,
				LevelKey:     "Level",
				EncodeLevel:  zapcore.CapitalLevelEncoder,
				CallerKey:    "caller",
				EncodeCaller: zapcore.ShortCallerEncoder,
			},
			Level: zap.NewAtomicLevelAt(zapLogLevel),
		}.Build(zap.AddCallerSkip(1))

		if err != nil {
			panic(fmt.Sprintf("Init logger fail, %s", err.Error()))
		}

		ZapLogger = &Logger{
			Zap: l.Sugar(),
		}

	})
	return ZapLogger
}

func GetLogger() *Logger {
	return ZapLogger
}
