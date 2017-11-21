package logger

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

const defaultTimeFormat = "15:04:05.000"

//go:generate evt2gogen -t Message

type Message struct {
	Timestamp time.Time
	Prefix    string
	Level     Level
	Msg       string
	Nr        uint64

	tmpl *template.Template
}

func NewMessage(nr uint64, prefix string, lvl Level, tmpl *template.Template, format string, val ...interface{}) (m Message) {
	m = Message{
		Timestamp: time.Now(),
		Prefix:    prefix,
		Level:     lvl,
		Msg:       fmt.Sprintf(format, val...),
		Nr:        nr,
		tmpl:      tmpl,
	}
	return
}

func (m Message) Format() string {
	var buf bytes.Buffer
	m.tmpl.Execute(&buf, m)
	return Paint(
		fmt.Sprintf(buf.String(), m.Msg),
		m.Level.Foreground,
		m.Level.ForegroundHi,
		m.Level.Background,
		m.Level.BackgroundHi,
	)
}

func (m Message) FormattedTime() string {
	return m.Timestamp.Format(defaultTimeFormat)
}
