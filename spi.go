package logs

const (
	ConsoleFormatter = Formatter("console")
	JsonFormatter    = Formatter("json")

	DebugLevel = Level(0)
	InfoLevel  = Level(1)
	WarnLevel  = Level(2)
	ErrorLevel = Level(3)

	defaultLogName = "AACFACTORY"

	defaultTimeFormatter = "[3:04:05 PM]"
)

type Level int

type Formatter string

type Logger interface {
	With(key string, value interface{}) Logger
	DebugEnabled() (ok bool)
	Debug() (event Event)
	InfoEnabled() (ok bool)
	Info() (event Event)
	WarnEnabled() (ok bool)
	Warn() (event Event)
	ErrorEnabled() (ok bool)
	Error() (event Event)
}

type Event interface {
	Message(message string)
	Cause(err error) Event
	Caller() Event
	With(key string, value interface{}) Event
}
