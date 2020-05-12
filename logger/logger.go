package logger

import (
	"encoding/json"
	"sync"

	"go.uber.org/zap"
)

type Logger struct {
	Zap *zap.SugaredLogger
}

var (
	ZapLogger *Logger
	once      sync.Once
)

func GetLogger() *Logger {

	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": [/tmp/trader"],
	  "errorOutputPaths": ["stderr"],
	  "initialFields": {"foo": "bar"},
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	once.Do(func() {
		var cfg zap.Config

		if err := json.Unmarshal(rawJSON, &cfg); err != nil {
			panic(err)
		}
		logger, err := cfg.Build()
		if err != nil {
			panic("Init logger failed")
		}
		ZapLogger = &Logger{Zap: logger.Sugar()}

	})
	return ZapLogger
}
