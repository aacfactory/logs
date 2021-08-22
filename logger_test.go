package logs_test

import (
	"github.com/aacfactory/logs"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	log, err := logs.New(
		logs.Name("name"),
		logs.Writer(os.Stdout),
		logs.WithLevel(logs.DebugLevel),
		logs.WithFormatter(logs.ConsoleFormatter),
	)
	if err != nil {
		t.Error(err)
		return
	}
	log.Debug().Caller().Message("foo")
	log.Info().Caller().Message("foo")
	log.Warn().Caller().Message("foo")
	log.Error().Caller().Message("foo")

}
