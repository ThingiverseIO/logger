package logger

import (
	"fmt"
	"html/template"
	"os"
	"sync/atomic"
)

var (
	defaultBackend  Backend
	defaultTemplate *template.Template
)


func init() {
	defaultTemplate, _ = template.New("tvio_logger_default").Parse(
		"{{.FormattedTime}} \u2771\u2771 {{.Prefix}} \u2771\u2771 {{.Level.Name}} \u2771\u2771\t%s")
	defaultBackend = NewWriterBackend(os.Stdout)
}

type Logger struct {
	Prefix string

	backend  Backend
	debug    bool
	levels   map[string]Level
	logNr    uint64
	template *template.Template
}

func New(prefix string) (logger *Logger) {

	logger = &Logger{
		Prefix:   prefix,
		backend:  defaultBackend,
		template: defaultTemplate,
		levels:   DefaultLevels,
	}
	return

}

func (l *Logger) SetBackend(backend Backend) *Logger {
	l.backend = backend
	return l
}

func (l *Logger) SetDebug(debug bool) *Logger {
	l.debug = debug
	return l
}

func (l *Logger) SetTemplate(tmpl string) *Logger {
	l.template, _ = template.New("tvio_logger_default").Parse(tmpl)
	return l
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
	if !l.debug {
		return
	}
	msg := fmt.Sprint(message...)
	l.print(DEBUG, msg)
}

func (l *Logger) Debugf(format string, val ...interface{}) {
	if !l.debug {
		return
	}
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

func (l *Logger) print(level, format string, val ...interface{}) {
	lvl, ok := l.levels[level]
	if !ok {
		lvl = l.levels[INFO]
	}
	message := NewMessage(atomic.AddUint64(&l.logNr, 1), l.Prefix, lvl, l.template, format, val...)
	l.backend.Write(message)
}
