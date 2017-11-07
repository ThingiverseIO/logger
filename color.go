package logger

import "fmt"

type Color int

// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
const (
	BLACK Color = iota
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
)

func (c Color) AsForeground() Color {
	return c + 30
}

func (c Color) AsBackground() Color {
	return c + 40
}

func (c Color) AsHiForeground() Color {
	return c + 90
}

func (c Color) AsHiBackground() Color {
	return c + 100
}

func Paint(s string, foreground Color, foregroundHi bool, background Color, backgroundHi bool) string {
	var fg, bg Color
	if foregroundHi {
		fg = foreground.AsHiForeground()
	} else {
		fg = foreground.AsForeground()
	}
	if backgroundHi {
		bg = background.AsHiBackground()
	} else {
		bg = background.AsBackground()
	}
	return fmt.Sprintf("\x1b[%d;%dm %s \x1b[0m", fg, bg, s)
}
