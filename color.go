package logs

import "fmt"

// Foreground colors.
const (
	red      color = 31
	green    color = 32
	yellow   color = 33
	purple   color = 34
	pink     color = 35
	blue     color = 36
	grey     color = 37
	redBg    color = 41
	greenBg  color = 42
	yellowBg color = 43
	purpleBg color = 44
	pinkBg   color = 45
	blueBg   color = 46
	greyBg   color = 47
)

// Color represents a text color.
type color uint8

func colorable(s string, c color) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}
