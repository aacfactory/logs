package logs

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"github.com/valyala/bytebufferpool"
	"io"
	"os"
	"time"
)

const (
	TextFormatter ConsoleWriterFormatter = iota
	ColorTextFormatter
	JsonFormatter
)

type ConsoleWriterFormatter int

const (
	StdOut ConsoleWriterOutType = iota
	StdErr
	StdMix
)

type ConsoleWriterOutType int

func NewConsoleWriter(formatter ConsoleWriterFormatter, out ConsoleWriterOutType) (w Writer) {
	var outWriter io.Writer
	var errWriter io.Writer
	switch out {
	case StdErr:
		outWriter = os.Stderr
		errWriter = os.Stderr
		break
	case StdMix:
		outWriter = os.Stdout
		errWriter = os.Stderr
		break
	default:
		outWriter = os.Stdout
		errWriter = os.Stdout
		break
	}
	var encoder EntryEncoder
	switch formatter {
	case ColorTextFormatter:
		encoder = &TextEntryEncoder{
			colorable:       true,
			debugWriterTo:   DebugLevel.ColorableLevelWriterTo(),
			infoWriterTo:    InfoLevel.ColorableLevelWriterTo(),
			warnWriterTo:    WarnLevel.ColorableLevelWriterTo(),
			errorWriterTo:   ErrorLevel.ColorableLevelWriterTo(),
			timeColor:       color.New(color.FgHiBlack),
			callerColor:     color.New(color.FgHiBlue, color.Underline),
			callerFnColor:   color.New(color.FgHiBlack),
			fieldKeyColor:   color.New(color.FgHiCyan),
			fieldValueColor: color.New(color.FgHiBlack),
			causeColor:      color.New(color.BgHiRed),
		}
		outWriter = colorable.NewColorable(outWriter.(*os.File))
		errWriter = colorable.NewColorable(errWriter.(*os.File))
		break
	case JsonFormatter:
		encoder = &JsonEntryEncoder{}
		break
	default:
		encoder = &TextEntryEncoder{}
		break
	}

	w = &ConsoleWriter{
		encoder: encoder,
		out:     outWriter,
		err:     errWriter,
	}
	return
}

type ConsoleWriter struct {
	encoder EntryEncoder
	out     io.Writer
	err     io.Writer
}

func (writer *ConsoleWriter) Write(entry Entry) {
	p := writer.encoder.Encode(entry)
	if entry.Level == ErrorLevel {
		_, _ = writer.err.Write(p)
	} else {
		_, _ = writer.out.Write(p)
	}
	return
}

func (writer *ConsoleWriter) Close() (err error) {
	return
}

type EntryEncoder interface {
	Encode(entry Entry) (p []byte)
}

type TextEntryEncoder struct {
	colorable       bool
	debugWriterTo   io.WriterTo
	infoWriterTo    io.WriterTo
	warnWriterTo    io.WriterTo
	errorWriterTo   io.WriterTo
	timeColor       *color.Color
	callerColor     *color.Color
	callerFnColor   *color.Color
	fieldKeyColor   *color.Color
	fieldValueColor *color.Color
	causeColor      *color.Color
}

func (encoder *TextEntryEncoder) Encode(entry Entry) (p []byte) {
	buf := bytebufferpool.Get()
	_, _ = buf.Write(lqb)
	if encoder.colorable {
		switch entry.Level {
		case DebugLevel:
			_, _ = encoder.debugWriterTo.WriteTo(buf)
			break
		case InfoLevel:
			_, _ = encoder.infoWriterTo.WriteTo(buf)
			break
		case WarnLevel:
			_, _ = encoder.warnWriterTo.WriteTo(buf)
			break
		case ErrorLevel:
			_, _ = encoder.errorWriterTo.WriteTo(buf)
			break
		default:
			_, _ = buf.Write(entry.Level.Bytes())
			break
		}
	} else {
		_, _ = buf.Write(entry.Level.Bytes())
	}
	_, _ = buf.Write(rqb)
	_, _ = buf.Write(space)
	_, _ = buf.Write(lqb)
	if encoder.colorable {
		_, _ = encoder.timeColor.Fprintf(buf, "%s", entry.Occur.Format(time.DateTime))
	} else {
		_, _ = buf.WriteString(entry.Occur.Format(time.DateTime))
	}
	_, _ = buf.Write(rqb)
	_, _ = buf.Write(space)
	if entry.Caller.File != "" {
		if encoder.colorable {
			_, _ = encoder.callerColor.Fprintf(buf, "%s:%d", entry.Caller.File, entry.Caller.Line)
			_, _ = buf.Write(space)
			_, _ = buf.Write(lab)
			_, _ = buf.Write(space)
			_, _ = encoder.callerFnColor.Fprintf(buf, "%s", entry.Caller.Fn)
		} else {
			_, _ = buf.WriteString(fmt.Sprintf("%s:%d > %s", entry.Caller.File, entry.Caller.Line, entry.Caller.Fn))
		}
	}
	_, _ = buf.Write(newline)
	_, _ = buf.WriteString(entry.Message)
	if entry.Cause != nil {
		_, _ = buf.Write(newline)
		if encoder.colorable {
			_, _ = encoder.causeColor.Fprintf(buf, ">>>>>>>>>>>>> ERROR <<<<<<<<<<<<<")
		} else {
			_, _ = buf.WriteString(">>>>>>>>>>>>> ERROR <<<<<<<<<<<<<")
		}
		_, _ = buf.Write(newline)
		_, _ = buf.WriteString(fmt.Sprintf("%+v", entry.Cause))
	}
	if len(entry.Fields) > 0 {
		_, _ = buf.Write(newline)
		for i, field := range entry.Fields {
			if i > 0 {
				_, _ = buf.Write(space)
			}
			if encoder.colorable {
				_, _ = encoder.fieldKeyColor.Fprintf(buf, "%s", field.Key)
				_, _ = buf.Write(equal)
				_, _ = encoder.fieldValueColor.Fprintf(buf, "%s", string(field.ValueBytes()))
			} else {
				_, _ = buf.WriteString(fmt.Sprintf("%s=%s", field.Key, string(field.ValueBytes())))
			}
		}
	}
	_, _ = buf.Write(newline)
	p = buf.Bytes()
	bytebufferpool.Put(buf)
	return
}

type JsonEntryEncoder struct {
}

func (encoder *JsonEntryEncoder) Encode(entry Entry) (p []byte) {
	p, _ = entry.MarshalJSON()
	p = append(p, newline...)
	return
}
