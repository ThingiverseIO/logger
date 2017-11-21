package logger

const (
	INFO    string = "INFO"
	INIT           = "INIT"
	ERROR          = "ERROR"
	FATAL          = "FATAL"
	DEBUG          = "DEBUG"
	SUCCESS        = "SUCCESS"
	WARNING        = "WARNING"
)

var DefaultLevels = map[string]Level{
	INFO:    Level{Name: INFO, Foreground: GREEN},
	INIT:    Level{Name: INIT, Foreground: BLUE, ForegroundHi: true},
	WARNING: Level{Name: WARNING, Background: YELLOW, BackgroundHi: true},
	ERROR:   Level{Name: ERROR, Foreground: WHITE, Background: RED},
	FATAL:   Level{Name: FATAL, Foreground: WHITE, Background: RED, BackgroundHi: true},
	SUCCESS: Level{Name: SUCCESS, Foreground: CYAN},
	DEBUG:   Level{Name: DEBUG, Background: YELLOW, BackgroundHi: true},
}

type Level struct {
	Name string

	Foreground   Color
	ForegroundHi bool
	Background   Color
	BackgroundHi bool

	Display bool
}
