package logs_test

import (
	"errors"
	"github.com/aacfactory/logs"
	"testing"
	"time"
)

func TestLog(t *testing.T) {

	log := logs.New(logs.LogOption{
		Name:             "FOO",
		Formatter:        logs.LogConsoleFormatter,
		ActiveLevel:      logs.LogDebugLevel,
		Colorable:        true,
		EnableCaller:     true,
		EnableStacktrace: false,
	})

	wLog := logs.With(log, logs.F("foo", "bar"), logs.F("baz", "xxx"))

	wLog.Info("with")

	log.Debug("Debug")
	log.Debugf("Debug %s", "debug")
	log.Debugw("Debug", "k", "v", "t", time.Now())

	log.Info("Info", "debug")
	log.Infof("Info %s", "debug")
	log.Infow("Info", "k", "v", "t", time.Now())

	log.Warn("Warn")
	log.Warnf("Warn %s", "debug")
	log.Warnw("Warn", "k", "v", "t", time.Now())

	log.Error("Error", "debug")
	log.Errorf("Error %s", "debug")
	log.Errorw("Error", "k", "v", "t", time.Now())
	logs.With(log, logs.Error(errors.New("kv"))).Error("foo")
	_ = log.Sync()

}
