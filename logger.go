package logger

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"sync"
	"time"
)

type Level string

const (
	INFO    Level = "INFO"
	INIT    Level = "INIT"
	ERROR   Level = "ERROR"
	FATAL   Level = "FATAL"
	DEBUG   Level = "DEBUG"
	SUCCESS Level = "SUCCESS"
)

type LevelProperties struct {
	Foreground   Color
	ForegroundHi bool
	Background   Color
	BackgroundHi bool
}

var DefaultLevelProperties = map[Level]LevelProperties{
	INFO:    LevelProperties{Foreground: GREEN},
	INIT:    LevelProperties{Foreground: BLUE, ForegroundHi: true},
	ERROR:   LevelProperties{Foreground: WHITE, Background: RED},
	FATAL:   LevelProperties{Foreground: WHITE, Background: RED, BackgroundHi: true},
	DEBUG:   LevelProperties{Background: YELLOW, BackgroundHi: true},
	SUCCESS: LevelProperties{Foreground: CYAN},
}

type Message struct {
	Timestamp string
	Prefix    string
	Level     Level
	Msg       string
}

func NewMessage(prefix, timeFormat string, lvl Level, msg string, val ...interface{}) (m Message) {
	m = Message{
		Timestamp: time.Now().Format(timeFormat),
		Prefix:    prefix,
		Level:     lvl,
		Msg:       fmt.Sprintf(msg, val...),
	}
	return
}

func (m Message) ExecuteTemplate(tmpl *template.Template) string {
	var buf bytes.Buffer
	tmpl.Execute(&buf, m)
	return fmt.Sprintf(buf.String(), m.Msg)
}

var (
	defaultBackend  Backend
	defaultTemplate *template.Template
)

func SetDefaultTemplate(tmpl string) {
	defaultTemplate, _ = template.New("tvio_logger_default").Parse(tmpl)
}

func GetDefaultTemplate() *template.Template {
	return defaultTemplate
}

const defaultTimeFormat = "15:04:05.000"

type Backend interface {
	Write(msg string)
}

type MultiWriteBackend struct {
	m      *sync.Mutex
	Writer []io.Writer
}

func (b *MultiWriteBackend) Write(msg string) {
	b.m.Lock()
	defer b.m.Unlock()
	for _, w := range b.Writer {
		fmt.Fprintln(w, msg)
	}
}

func NewMultiWriteBackend(w ...io.Writer) (b *MultiWriteBackend) {
	b = &MultiWriteBackend{
		m:      &sync.Mutex{},
		Writer: w,
	}
	return
}

func SetDefaultBackend(b Backend) {
	defaultBackend = b
}

func GetDefaultBackend() Backend {
	return defaultBackend
}

func init() {
	SetDefaultBackend(NewMultiWriteBackend(os.Stdout))
	SetDefaultTemplate("{{.Timestamp}} \u2771\u2771 {{.Prefix}} \u2771\u2771 {{.Level}} \u2771\u2771\t%s")
}

type Logger struct {
	Backend    Backend
	Levels     map[Level]LevelProperties
	Template   *template.Template
	TimeFormat string
	Prefix     string
}

func New(prefix string) (logger *Logger) {

	logger = &Logger{
		Backend:    defaultBackend,
		Levels:     DefaultLevelProperties,
		Template:   defaultTemplate,
		TimeFormat: defaultTimeFormat,
		Prefix:     prefix,
	}
	return
}

func (l *Logger) Init(message ...interface{}) {
	msg := fmt.Sprint(message...)
	l.print(INIT, msg)
}

func (l *Logger) Initf(format string, val ...interface{}) {
	l.print(INIT, format, val...)
}

func (l *Logger) Info(message ...interface{}) {
	msg := fmt.Sprint(message...)
	l.print(INFO, msg)
}

func (l *Logger) Infof(format string, val ...interface{}) {
	l.print(INFO, format, val...)
}

func (l *Logger) Debug(message ...interface{}) {
	msg := fmt.Sprint(message...)
	l.print(DEBUG, msg)
}

func (l *Logger) Debugf(format string, val ...interface{}) {
	l.print(DEBUG, format, val...)
}

func (l *Logger) Error(message ...interface{}) {
	msg := fmt.Sprint(message...)
	l.print(ERROR, msg)
}

func (l *Logger) Errorf(format string, val ...interface{}) {
	l.print(ERROR, format, val...)
}

func (l *Logger) Success(message ...interface{}) {
	msg := fmt.Sprint(message...)
	l.print(SUCCESS, msg)
}

func (l *Logger) Successf(format string, val ...interface{}) {
	l.print(SUCCESS, format, val...)
}

func (l *Logger) PrintSuccess() {
	l.Success("Success")
}

func (l *Logger) print(level Level, format string, val ...interface{}) {
	properties, ok := l.Levels[level]
	if !ok {
		properties = DefaultLevelProperties[INFO]
	}
	message := NewMessage(l.Prefix, l.TimeFormat, level, format, val...).ExecuteTemplate(l.Template)
	message = Paint(message, properties.Foreground, properties.ForegroundHi, properties.Background, properties.BackgroundHi)
	l.Backend.Write(message)
}
