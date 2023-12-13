package logs_test

import (
	"context"
	"errors"
	"github.com/aacfactory/logs"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	log, logErr := logs.New(
		logs.WithLevel(logs.DebugLevel),
		logs.WithConsoleWriterFormatter(logs.ColorTextFormatter),
	)
	if logErr != nil {
		t.Error(logErr)
		return
	}

	log.Debug().Caller().With("f1", "f1").With("f2", 2).Message("debug")
	log.Info().Caller().With("time", time.Now()).Message("info")
	log.Warn().Caller().Message("warn")
	log.Error().Caller().With("time", time.Now()).Cause(errors.New("some error")).Message("error")

	logErr = log.Shutdown(context.Background())
	if logErr != nil {
		t.Error(logErr)
		return
	}
}

func TestLogger_Json(t *testing.T) {
	log, logErr := logs.New(
		logs.WithLevel(logs.DebugLevel),
		logs.WithConsoleWriterFormatter(logs.JsonFormatter),
	)
	if logErr != nil {
		t.Error(logErr)
		return
	}

	log.Debug().Caller().With("f1", "f1").With("f2", 2).Message("debug")
	log.Info().With("time", time.Now()).Message("info")
	log.Warn().Caller().Message("warn")
	log.Error().Caller().Cause(errors.New("some error")).Message("error")

	logErr = log.Shutdown(context.Background())
	if logErr != nil {
		t.Error(logErr)
		return
	}

}

func BenchmarkLogger_Info(b *testing.B) {
	log, logErr := logs.New()
	if logErr != nil {
		b.Error(logErr)
		return
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info().Message("info")
	}
	logErr = log.Shutdown(context.Background())
	if logErr != nil {
		b.Error(logErr)
		return
	}
}
