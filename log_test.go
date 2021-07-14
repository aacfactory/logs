package logs_test

import (
	"github.com/aacfactory/logs"
	"testing"
	"time"
)

func TestLog(t *testing.T) {

	wLog, withErr := logs.With(logs.Log(), logs.F("foo", "bar"), logs.F("baz", "xxx"))
	if withErr != nil {
		t.Error(withErr)
		return
	}
	wLog.Info("with")

	logs.Log().Debug("Debug")
	logs.Log().Debugf("Debug %s", "debug")
	logs.Log().Debugw("Debug", "k", "v", "t", time.Now())

	logs.Log().Info("Info", "debug")
	logs.Log().Infof("Info %s", "debug")
	logs.Log().Infow("Info", "k", "v", "t", time.Now())

	logs.Log().Warn("Warn")
	logs.Log().Warnf("Warn %s", "debug")
	logs.Log().Warnw("Warn", "k", "v", "t", time.Now())

	logs.Log().Error("Error", "debug")
	logs.Log().Errorf("Error %s", "debug")
	logs.Log().Errorw("Error", "k", "v", "t", time.Now())

	logs.Sync()

}
