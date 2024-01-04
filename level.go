package logs

import (
	"github.com/fatih/color"
	"io"
	"unsafe"
)

const (
	DebugLevel Level = iota + 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

var (
	unknownLevelBytes = []byte("UNK")
	debugLevelBytes   = []byte("DBG")
	infoLevelBytes    = []byte("INF")
	warnLevelBytes    = []byte("WRN")
	errorLevelBytes   = []byte("ERR")
)

type Level int

func (level Level) MarshalJSON() ([]byte, error) {
	return append(append([]byte{'"'}, level.Bytes()...), '"'), nil
}

func (level Level) Bytes() []byte {
	switch level {
	case DebugLevel:
		return debugLevelBytes
	case InfoLevel:
		return infoLevelBytes
	case WarnLevel:
		return warnLevelBytes
	case ErrorLevel:
		return errorLevelBytes
	default:
		return unknownLevelBytes
	}
}

func (level Level) String() string {
	return string(level.Bytes())
}

func (level Level) ColorableLevelWriterTo() (w io.WriterTo) {
	var c *color.Color
	switch level {
	case DebugLevel:
		c = color.New(color.FgHiGreen, color.Bold)
		break
	case InfoLevel:
		c = color.New(color.FgHiCyan, color.Bold)
		break
	case WarnLevel:
		c = color.New(color.FgHiYellow, color.Bold)
		break
	case ErrorLevel:
		c = color.New(color.FgHiRed, color.Bold)
		break
	}
	c.EnableColor()
	w = &ColorableLevelWriterTo{
		color: c,
		level: level,
	}
	return
}

type ColorableLevelWriterTo struct {
	color *color.Color
	level Level
}

func (lc ColorableLevelWriterTo) WriteTo(writer io.Writer) (int64, error) {
	s := lc.color.Sprint(lc.level.String())
	p := unsafe.Slice(unsafe.StringData(s), len(s))
	n, err := writer.Write(p)
	return int64(n), err
}
