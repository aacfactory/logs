package logs

import (
	"github.com/rs/zerolog"
	"io"
	"strings"
)

func buildConsoleWriter(out io.Writer, colored bool, timeFormatter string) (w io.Writer) {

	w = zerolog.NewConsoleWriter(
		func(cw *zerolog.ConsoleWriter) {
			cw.Out = out
		},
		func(cw *zerolog.ConsoleWriter) {
			cw.NoColor = !colored
		},
		func(cw *zerolog.ConsoleWriter) {
			cw.TimeFormat = timeFormatter
		},
		func(cw *zerolog.ConsoleWriter) {
			cw.FormatLevel = consoleLevelFormatter(colored)
		},
	)

	return
}

func consoleCallerFormatter(value interface{}) (result string) {
	if value == nil {
		return
	}
	caller, ok := value.(string)
	if !ok {
		return
	}
	if strings.IndexByte(caller, '/') == 0 || strings.IndexByte(caller, ':') == 1 {
		idx := strings.Index(caller, "/src/")
		if idx > 0 {
			caller = caller[idx+5:]
		}
	}
	result = caller
	return
}

func consoleLevelFormatter(c bool) (fn zerolog.Formatter) {
	fn = func(i interface{}) (s string) {
		value, ok := i.(string)
		if !ok {
			if c {
				s = colorable("???", red)
			} else {
				s = "???"
			}
			return
		}
		switch value {
		case "debug":
			if c {
				s = colorable("DEBUG", grey)
			} else {
				s = "DEBUG"
			}
		case "info":
			if c {
				s = colorable("INFO ", blue)
			} else {
				s = "INFO "
			}
		case "warn":
			if c {
				s = colorable("WARN ", yellow)
			} else {
				s = "WARN "
			}
		case "error":
			if c {
				s = colorable("ERROR", pink)
			} else {
				s = "ERROR"
			}
		default:
			if c {
				s = colorable("???", red)
			} else {
				s = "???"
			}
		}
		return
	}
	return
}
