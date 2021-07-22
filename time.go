package logs

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	s := t.Format(time.RFC3339Nano)
	s1 := s[:strings.LastIndexByte(s, '.')]
	s2 := s[strings.LastIndexByte(s, '.')+1 : strings.LastIndexByte(s, '+')]
	for i := len(s2); i < 9; i++ {
		s2 = s2 + "0"
	}
	s3 := s[strings.LastIndexByte(s, '+')+1:]
	enc.AppendString(fmt.Sprintf("%s.%s+%s", s1, s2, s3))
}
