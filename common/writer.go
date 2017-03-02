package common

import (
	"bytes"
	"fmt"
)

// Writer is a formatted output for writing to a file or buffer.
type Writer struct {
	out *bytes.Buffer
}

// NewWriter creates a new writer.
func NewWriter() *Writer {
	return &Writer{
		out: &bytes.Buffer{},
	}
}

// Append writes the other buffer to this buffer.
func (w *Writer) Append(other *Writer) {
	if other != nil {
		w.out.Write(other.out.Bytes())
	}
}

// Write writes text to the buffer.
func (w *Writer) Write(data ...interface{}) {
	w.out.WriteString(fmt.Sprint(data...))
}

// Writeln writes text with a newline at the end of the write.
func (w *Writer) Writeln(data ...interface{}) {
	data = append(data, "\n")
	w.Write(data...)
}
