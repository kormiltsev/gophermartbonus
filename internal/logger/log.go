// package logger is for using external zap logger.
package logger

import (
	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Blogs is logger we used.
type Blogs struct {
	Logger *zap.Logger
}

// NewLog creates new logger. Using file logs.log.
func NewLog() *Blogs {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["./logs.log"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000")
	cfg.EncoderConfig.StacktraceKey = "" // to hide stacktrace info

	logger := zap.Must(cfg.Build())
	defer logger.Sync()

	return &Blogs{logger}
}
