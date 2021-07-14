package logs

import "fmt"

// Foreground colors.
const (
	red     Color = 35
	yellow  Color = 33
	blue    Color = 36
	magenta Color = 37
)

// Color represents a text color.
type Color uint8

// Add adds the coloring to the given string.
func (c Color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}
