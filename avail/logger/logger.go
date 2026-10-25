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

// gainLock will mutex lock the logger and return the unlock function
// so that `defer imp.gainLock()()` can be used to lock the function.
//
// The lock is intended to protect the logger from race conditions, however,
// the logger has a shared indentation so calling from different go-routines (threads)
// will still produce wrong intents. Create multiple loggers if the threads
// are outputting to two different outputs. If they are outputting to the same
// out, then jumping back and forth between different indents for different
// threads is not great, so there is no attempt to have non-shared indents.
// If the consumer is using multiple threads, then it is up to them to figure out
// how to make the logs look reasonable and not the responsibility of the logger.
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
		check(fmt.Fprint(imp.out, imp.prefix, "\n"))
		imp.flush()
		return
	}

	for ln := range strings.Lines(text) {
		check(fmt.Fprint(imp.out, imp.prefix, ln))
	}
	if text[count-1] != '\n' {
		check(fmt.Fprint(imp.out, "\n"))
	}
	imp.flush()
}

func check(_ int, err error) {
	if err != nil {
		panic(err)
	}
}

// flush will try to flush the output if the writer has `Sync`.
func (imp *Logger) flush() {
	if s, ok := imp.out.(interface{ Sync() error }); ok {
		if err := s.Sync(); err != nil {
			panic(err)
		}
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
// This returns the `Dedent` method to allow, `defer log.Indent()()`.
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
// The returned function is only intended to be called once.
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
