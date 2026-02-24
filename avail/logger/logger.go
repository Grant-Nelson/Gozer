package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const indentSpacer = "┆ "

type Logger struct {
	indent int
	prefix string
	out    io.Writer
	lock   *sync.Mutex
}

func Nop() *Logger {
	return (*Logger)(nil)
}

func New(verbose bool, out io.Writer) *Logger {
	if !verbose {
		return Nop()
	}
	if out == nil {
		out = os.Stderr
	}
	return &Logger{
		indent: 0,
		prefix: ``,
		out:    out,
		lock:   &sync.Mutex{},
	}
}

func (imp *Logger) gainLock() func() {
	imp.lock.Lock()
	return imp.lock.Unlock
}

func (imp *Logger) setIndent(delta int) {
	imp.indent = max(imp.indent+delta, 0)
	imp.prefix = strings.Repeat(indentSpacer, imp.indent)
}

func (imp *Logger) output(text string) {
	count := len(text)
	if count <= 0 {
		fmt.Fprint(imp.out, imp.prefix, "\n")
		return
	}
	for ln := range strings.Lines(text) {
		fmt.Fprint(imp.out, imp.prefix, ln)
	}
	if text[count-1] != '\n' {
		fmt.Fprint(imp.out, "\n")
	}
}

func (imp *Logger) Print(v ...any) {
	if imp != nil {
		defer imp.gainLock()()
		imp.output(fmt.Sprint(v...))
	}
}

func (imp *Logger) Printf(format string, v ...any) {
	if imp != nil {
		defer imp.gainLock()()
		imp.output(fmt.Sprintf(format, v...))
	}
}

func (imp *Logger) Println(v ...any) {
	if imp != nil {
		defer imp.gainLock()()
		imp.output(fmt.Sprintln(v...))
	}
}

// Indent will indent all the following logs.
// This returns the Dedent method to allow, `defer log.Indent()()`.
func (imp *Logger) Indent() func() {
	if imp != nil {
		defer imp.gainLock()()
		imp.setIndent(1)
	}
	return imp.Dedent
}

// Dedent will reduce the indent for all the following logs.
func (imp *Logger) Dedent() {
	if imp != nil {
		defer imp.gainLock()()
		imp.setIndent(-1)
	}
}

// LogGroup logs the start of a group and indicate the end of the
// group when the returned function is called.
// Example: `defer log.LogGroup("Reading %q", filename)()`
func (imp *Logger) LogGroup(format string, args ...any) func() {
	if imp == nil {
		return func() {}
	}

	defer imp.gainLock()()
	text := fmt.Sprintf(format, args...)
	imp.output(fmt.Sprint(text, "...\n"))
	imp.setIndent(1)
	start := time.Now()

	return func() {
		defer imp.gainLock()()
		since := time.Since(start)
		imp.setIndent(-1)
		imp.output(fmt.Sprint(text, `... Done (`, since, ")\n"))
	}
}
