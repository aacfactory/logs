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
