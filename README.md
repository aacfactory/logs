# 概述
LOGS 是一个AAC FACTORY LOG的规范。
# 获取
```go
go get github.com/aacfactory/logs
```
# 使用
```go
log := logs.New(logs.LogOption{
    Name:        "FOO",
    Formatter:   logs.LogConsoleFormatter,
    ActiveLevel: logs.LogDebugLevel,
    Colorable:   true,
})

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

_ = log.Sync()
```