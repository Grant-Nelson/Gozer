package coder

import (
	"fmt"
	"os"
	"strings"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type Coder struct {
	buf          *strings.Builder
	filename     string
	line         int
	column       int
	curIndent    string
	indentSymbol string
	indentCache  string
}

func New(filename, indentSymbol string) *Coder {
	return &Coder{
		buf:          &strings.Builder{},
		filename:     filename,
		line:         1,
		column:       0,
		curIndent:    ``,
		indentSymbol: indentSymbol,
		indentCache:  ``,
	}
}

func (c *Coder) Filename() string { return c.filename }
func (c *Coder) Line() int        { return c.line }
func (c *Coder) Column() int      { return c.column }
func (c *Coder) String() string   { return c.buf.String() }
func (c *Coder) Len() int         { return c.buf.Len() }

func (c *Coder) Indent() {
	if size := len(c.curIndent); size < len(c.indentCache) {
		c.curIndent = c.indentCache[:size+len(c.indentSymbol)]
	} else {
		c.indentCache += c.indentSymbol
		c.curIndent = c.indentCache
	}
}

func (c *Coder) Dedent() bool {
	if size := len(c.curIndent); size > 0 {
		c.curIndent = c.indentCache[:size-len(c.indentSymbol)]
		return true
	}
	return false
}

func (c *Coder) writeString(s string) {
	if _, err := c.buf.WriteString(s); err != nil {
		panic(faults.New(`unexpected error from the coder buffer`, err).
			With(`file`, c.filename).
			With(`line`, c.line).
			With(`column`, c.column).
			With(`offset`, c.buf.Len()))
	}
}

func (c *Coder) writeText(s string) {
	for ln := range strings.Lines(s) {
		if c.column <= 0 {
			c.writeString(c.curIndent)
			c.column = len(c.curIndent)
		}
		c.writeString(ln)
		c.column += len(ln)
		if strings.HasSuffix(ln, "\n") {
			c.line++
			c.column = 0
		}
	}
}

// Write formats using the default formats for its operands and writes to the buffer.
// Spaces are added between operands when neither is a string.
func (c *Coder) Write(a ...any) { c.writeText(fmt.Sprint(a...)) }

// Printf formats according to a format specifier and writes to the buffer.
func (c *Coder) WriteF(format string, a ...any) { c.writeText(fmt.Sprintf(format, a...)) }

// WriteLn formats using the default formats for its operands and writes to the buffer.
// Spaces are always added between operands and a newline is appended.
func (c *Coder) WriteLn(a ...any) { c.writeText(fmt.Sprintln(a...)) }

// Save writes the buffer out to the filename that was given.
func (c *Coder) Save() error {
	return os.WriteFile(c.filename, []byte(c.String()), 0o0644)
}
