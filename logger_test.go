package logs_test

import (
	"context"
	"errors"
	"github.com/aacfactory/logs"
	"os"
	"runtime"
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

func TestLogger_Shutdown(t *testing.T) {
	log, logErr := logs.New(
		logs.WithShutdownTimeout(3*time.Second),
		logs.DisableConsoleWriter(),
		logs.WithWriter(discardWriter()),
	)
	if logErr != nil {
		t.Error(logErr)
		return
	}
	go func(log logs.Logger) {
		for {
			log.Info().Message("info")
		}
	}(log)
	time.Sleep(3 * time.Second)
	err := log.Shutdown(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkLogger_Info(b *testing.B) {
	log, logErr := logs.New(
		logs.WithTimeoutDiscardLevel(logs.InfoLevel),
		logs.DisableConsoleWriter(),
		logs.WithWriter(discardWriter()),
		logs.WithConsumes(runtime.NumCPU()),
		logs.WithShutdownTimeout(3*time.Second),
	)
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

func discardWriter() *DiscardWriter {
	dev, _ := os.Open(os.DevNull)
	return &DiscardWriter{
		dev: dev,
	}
}

type DiscardWriter struct {
	dev *os.File
}

func (w *DiscardWriter) Write(entry logs.Entry) {
	p, _ := entry.MarshalJSON()
	_, _ = w.dev.Write(p)
}

func (w *DiscardWriter) Close() (err error) {
	return
}
