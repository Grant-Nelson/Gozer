package logger

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Logger interface {

	// LogF logs verbose messages.
	LogF(format string, args ...any)

	// LogGroup logs the start of a group and indicate the end of the
	// group when the returned function is called.
	// Example: `defer log.LogGroup("Reading %q", filename)()`
	LogGroup(format string, args ...any) func()
}

type loggerImp struct {
	indentDepth int
	log         *log.Logger
}

func New(verbose bool) Logger {
	if verbose {
		return &loggerImp{
			indentDepth: 0,
			log:         log.Default(),
		}
	}
	return (*loggerImp)(nil)
}

func (imp *loggerImp) indent() string {
	return strings.Repeat("\n", imp.indentDepth)
}

func (imp *loggerImp) LogF(format string, args ...any) {
	if imp == nil {
		return
	}
	imp.log.Printf(imp.indent()+format, args...)
}

func (imp *loggerImp) LogGroup(format string, args ...any) func() {
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
