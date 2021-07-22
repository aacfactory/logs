package logs

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

const (
	LogConsoleFormatter = LogFormatter("console")
	LogJsonFormatter    = LogFormatter("json")

	LogDebugLevel = LogLevel("debug")
	LogInfoLevel  = LogLevel("info")
	LogWarnLevel  = LogLevel("warn")
	LogErrorLevel = LogLevel("error")

	defaultLogName = "AACFACTORY"
)

type LogLevel string

type LogFormatter string

type LogOption struct {
	Name             string       `json:"name,omitempty"`
	Formatter        LogFormatter `json:"formatter,omitempty"`
	ActiveLevel      LogLevel     `json:"activeLevel,omitempty"`
	Colorable        bool         `json:"colorable,omitempty"`
	EnableCaller     bool         `json:"enableCaller,omitempty"`
	EnableStacktrace bool         `json:"enableStacktrace,omitempty"`
}

func New(option LogOption) Logs {

	name := strings.TrimSpace(option.Name)
	if name == "" {
		name = defaultLogName
	}

	formatter := option.Formatter
	if formatter == "" {
		formatter = LogConsoleFormatter
	}
	activeLevel := option.ActiveLevel
	if activeLevel == "" {
		activeLevel = LogInfoLevel
	}

	zapLevel := zapLogLevel(activeLevel)
	var callerEncoder zapcore.CallerEncoder
	if option.EnableCaller {
		if formatter == LogJsonFormatter {
			callerEncoder = zapcore.ShortCallerEncoder
			option.Colorable = false
		} else {
			callerEncoder = zapLogFullCallerEncoder
		}
	}

	var encodeLevel zapcore.LevelEncoder
	if !option.Colorable {
		encodeLevel = zapcore.CapitalLevelEncoder
	} else {
		encodeLevel = zapLogCapitalColorLevelEncoder
	}
	encodingConfig := zapcore.EncoderConfig{
		TimeKey:        "_T",
		LevelKey:       "_L",
		NameKey:        "_N",
		CallerKey:      "_C",
		MessageKey:     "_M",
		StacktraceKey:  "_S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   callerEncoder,
	}

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       false,
		Encoding:          string(formatter),
		EncoderConfig:     encodingConfig,
		DisableStacktrace: !option.EnableStacktrace,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		InitialFields:     map[string]interface{}{"app": name},
	}

	log, createEr := config.Build()
	if createEr != nil {
		panic(fmt.Errorf("logs create failed, %v", createEr))
	}

	return log.Sugar()
}

func With(log Logs, fields ...Field) (_log Logs) {

	sLog, ok := log.(*zap.SugaredLogger)
	if !ok {
		panic(fmt.Errorf("logs with fields failed, it is not *zap.SugaredLogger"))
		return
	}

	if fields == nil || len(fields) == 0 {
		_log = sLog.Desugar().Sugar()
		return
	}

	kvs := make([]zap.Field, 0, 1)
	for _, field := range fields {
		if field.Key == "error" {
			codeErr, isCodeErr := field.Value.(codeError)
			if isCodeErr {
				kvs = append(kvs, zap.Any(field.Key, codeErr))
				continue
			}
			kvs = append(kvs, zap.Any(field.Key, field.Value))
		}
	}

	_log = sLog.Desugar().With(kvs...).Sugar()
	return
}
