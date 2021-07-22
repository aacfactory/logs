package logs

type Field struct {
	Key   string
	Value interface{}
}

func F(key string, value interface{}) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

func Error(err error) Field {
	return Field{
		Key:   "error",
		Value: err,
	}
}

type codeError interface {
	Id() string
	Code() string
	FailureCode() int
	Message() string
	Stacktrace() (fn string, file string, line int)
}
