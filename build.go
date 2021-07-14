package logs

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

const (
	EnvFormatterKey = "LOG_FMT"
	EnvLevelKey     = "LOG_LEVEL"
	EnvNameKey      = "LOG_NAME"
	EnvColorable    = "LOG_COLORABLE"

	defaultLogFormatter = "console"
	defaultLogLevel     = "info"
	defaultLogName      = "AACFACTORY"
)

func build() (err error) {

	log, buildErr := config().Build()

	if buildErr != nil {
		err = fmt.Errorf("logs build failed, %v", buildErr)
		return
	}

	zapLog = log.Sugar()

	return
}

func config() (config zap.Config) {
	logFmt, _ := os.LookupEnv(EnvFormatterKey)
	if logFmt == "" {
		logFmt = defaultLogFormatter
	}
	logLevel, _ := os.LookupEnv(EnvLevelKey)
	if logLevel == "" {
		logLevel = defaultLogLevel
	}
	logName, _ := os.LookupEnv(EnvNameKey)
	if logName == "" {
		logName = defaultLogName
	}
	logColorable, _ := os.LookupEnv(EnvColorable)

	zapLevel := zapLogLevel(logLevel)
	var callerEncoder zapcore.CallerEncoder
	if logFmt == "json" {
		callerEncoder = zapcore.ShortCallerEncoder
	} else {
		callerEncoder = zapLogFullCallerEncoder
	}

	var encodeLevel zapcore.LevelEncoder
	if len(strings.TrimSpace(logColorable)) == 0 {
		encodeLevel = zapcore.CapitalLevelEncoder
	} else {
		encodeLevel = zapLogCapitalColorLevelEncoder
	}

	encodingConfig := zapcore.EncoderConfig{
		TimeKey:        "datetime",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   callerEncoder,
	}

	config = zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       false,
		Encoding:          logFmt,
		EncoderConfig:     encodingConfig,
		DisableStacktrace: true,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stdout"},
		InitialFields:     map[string]interface{}{"app": logName},
	}

	return
}
