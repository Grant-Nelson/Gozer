package coder

import (
	"fmt"
	"os"
	"strings"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type Coder struct {
	buf           *strings.Builder
	filepath      string
	line          int
	column        int
	curIndent     string
	indentSymbol  string
	indentCache   string
	commentSymbol string
}

func New(filepath, indentSymbol, commentSymbol string) *Coder {
	return &Coder{
		buf:           &strings.Builder{},
		filepath:      filepath,
		line:          1,
		column:        1,
		curIndent:     ``,
		indentSymbol:  indentSymbol,
		indentCache:   ``,
		commentSymbol: commentSymbol,
	}
}

func (c *Coder) FilePath() string { return c.filepath }
func (c *Coder) Line() int        { return c.line }
func (c *Coder) Column() int      { return c.column }
func (c *Coder) String() string   { return c.buf.String() }
func (c *Coder) Offset() int      { return c.buf.Len() }

// StartOfLine indicates if the last character that was written was a newline.
// This also returns true when nothing has been written.
func (c *Coder) StartOfLine() bool { return c.column <= 1 }

func (c *Coder) Indent() {
	if size := len(c.curIndent); size < len(c.indentCache) {
		c.curIndent = c.indentCache[:size+len(c.indentSymbol)]
	} else {
		c.indentCache += c.indentSymbol
		c.curIndent = c.indentCache
	}
}

func (c *Coder) Dedent() {
	if size := len(c.curIndent); size > 0 {
		c.curIndent = c.indentCache[:size-len(c.indentSymbol)]
	}
}

func (c *Coder) writeString(s string) {
	if len(s) <= 0 {
		return
	}
	if _, err := c.buf.WriteString(s); err != nil {
		panic(faults.New(`unexpected error from the coder buffer`, err).
			With(`filepath`, c.FilePath()).
			With(`line`, c.Line()).
			With(`column`, c.Column()).
			With(`offset`, c.Offset()))
	}
}

func (c *Coder) writeText(prefix, s string) {
	for ln := range strings.Lines(s) {
		if c.StartOfLine() {
			c.writeString(c.curIndent)
			c.writeString(prefix)
			c.column = len(c.curIndent) + 1
		}
		c.writeString(ln)
		c.column += len(ln)
		if strings.HasSuffix(ln, "\n") {
			c.line++
			c.column = 1
		}
	}
}

// Write writes using the default formats for its operands and writes to the buffer.
// Spaces are added between operands when neither is a string.
func (c *Coder) Write(a ...any) { c.writeText(``, fmt.Sprint(a...)) }

// WriteF writes according to a format specifier and writes to the buffer.
func (c *Coder) WriteF(format string, args ...any) { c.writeText(``, fmt.Sprintf(format, args...)) }

// WriteLn writes using the default formats for its operands and writes to the buffer.
// Spaces are always added between operands and a newline is appended.
func (c *Coder) WriteLn(args ...any) { c.writeText(``, fmt.Sprintln(args...)) }

// WriteComment works like WriteLn except it will insert the comment prefix
// at the start of each new line after the indent.
func (c *Coder) WriteComment(a ...any) { c.writeText(c.commentSymbol, fmt.Sprintln(a...)) }

// Save writes the buffer out to the filename that was given.
func (c *Coder) Save() error { return os.WriteFile(c.FilePath(), []byte(c.String()), 0o0644) }
