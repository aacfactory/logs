# 概述

LOGS 是一个 AAC FACTORY 的 LOG 规范。

## 获取

```go
go get github.com/aacfactory/logs
```

## 使用

```go
log, err := logs.New(
logs.Name("name"),
logs.Writer(os.Stdout),
logs.WithLevel(logs.DebugLevel),
logs.WithFormatter(logs.ConsoleFormatter),
)
if err != nil {
    // handle error
    return
}
log.Debug().Caller().Message("foo")
log.Info().Caller().Message("foo")
log.Warn().Caller().Message("foo")
log.Error().Caller().Message("foo")
```

## 引用感谢

* [rs/zerolog](https://github.com/rs/zerolog)