package logger

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Logger struct {
	indentDepth int
	log         *log.Logger
}

func New(verbose bool) *Logger {
	if verbose {
		return &Logger{
			indentDepth: 0,
			log:         log.Default(),
		}
	}
	return (*Logger)(nil)
}

func (imp *Logger) indent() string {
	return strings.Repeat("\n", imp.indentDepth)
}

// LogF logs verbose messages.
func (imp *Logger) LogF(format string, args ...any) {
	if imp == nil {
		return
	}
	imp.log.Printf(imp.indent()+format, args...)
}

// LogGroup logs the start of a group and indicate the end of the
// group when the returned function is called.
// Example: `defer log.LogGroup("Reading %q", filename)()`
func (imp *Logger) LogGroup(format string, args ...any) func() {
	if imp == nil {
		return func() {}
	}

	text := fmt.Sprintf(format, args...)
	imp.LogF(text + "...")
	imp.indentDepth++
	start := time.Now()

	return func() {
		since := time.Since(start)
		imp.indentDepth--
		imp.LogF(text+"... Done (%v)", since)
	}
}
