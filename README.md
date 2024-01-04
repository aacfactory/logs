# 概述

LOGS 是一个 AAC FACTORY 的 LOG 规范。

## 获取

```go
go get github.com/aacfactory/logs
```

## 使用

```go
log, logErr := logs.New(
    logs.WithLevel(logs.DebugLevel),
)
// handle logErr
// print
log.Debug().Caller().With("f1", "f1").With("f2", 2).Message("debug")
log.Info().With("time", time.Now()).Message("info")
log.Warn().Caller().Message("warn")
log.Error().Caller().Cause(errors.New("some error")).Message("error")
// shutdown
shutdownErr := log.Shutdown(context.Background())
```


## 压测
```shell
BenchmarkLogger_Info-20 1602939   744.5 ns/op   243 B/op  5 allocs/op
```