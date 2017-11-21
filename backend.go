package logger

import (
	"fmt"
	"io"
)

type Backend interface {
	Write(msg Message)
}

type WriterBackend struct {
	writer      io.Writer
	writeStream *MessageStreamController
}

func NewWriterBackend(writer io.Writer) (b *WriterBackend) {
	b = &WriterBackend{
		writer:      writer,
		writeStream: NewMessageStreamController(),
	}
	b.writeStream.Stream().Listen(b.onWrite)
	return
}

func (b *WriterBackend) Write(msg Message) {
	b.writeStream.Add(msg)
}

func (b *WriterBackend) onWrite(msg Message) {
	fmt.Fprintln(b.writer, msg.Format())
}
