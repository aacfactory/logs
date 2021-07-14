# 概述
LOGS 是一个AAC FACTORY LOG的规范。
# 获取
```go
go get github.com/aacfactory/logs
```
# 使用
环境变量

| 键            | 值类型 | 默认值     | 说明                                  |
| ------------- | ------ | ---------- | ------------------------------------- |
| LOG_FMT       | 枚举   | console    | 输出格式，值有：console json          |
| LOG_LEVEL     | 枚举   | info       | 日志等级，值有：debug info warn error |
| LOG_NAME      | 文本   | AACFACTORY | 日志名称，一般用于标识应用名称        |
| LOG_COLORABLE | 文本   |            | 是否开启颜色输出，非空即为开启。      |

```go
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
```