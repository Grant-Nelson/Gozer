package common

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
	"runtime/debug"
	"strings"
)

// ThrowError throws an error via a panic.
func ThrowError(msg ...interface{}) {
	panic(fmt.Errorf("%s", fmt.Sprint(msg...)))
}

// Logger handles logging messages for the transpiler.
type Logger struct {

	// output is the writer to print to or nil to print to the console.
	output io.Writer

	// errorCount is the number of errors which have been logged.
	errorCount int

	// showStacks indicates if error should also print stacktraces.
	showStacks bool
}

// NewLogger creates a new logger.
func NewLogger() *Logger {
	return &Logger{
		output:     nil,
		errorCount: 0,
		showStacks: false,
	}
}

// SetOutput sets the output for the log.
// If the output is set to nil then the log is printed to the console.
func (log *Logger) SetOutput(output io.Writer) {
	log.output = output
}

// ShowStackTraces indicates if stacktraces should be shown with errors.
func (log *Logger) ShowStackTraces(showStacks bool) {
	log.showStacks = showStacks
}

// StackTrace gets the current stack strace.
func (log *Logger) StackTrace(minLevel int, maxLevel int, indent string) string {
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")
	maxCut := maxLevel*2 + 1
	if len(lines) > maxCut {
		lines = lines[:maxCut]
	}
	minCut := minLevel*2 + 1
	if len(lines) > minCut {
		lines = lines[minCut:]
	}
	return indent + strings.Join(lines, "\n"+indent)
}

// Write prints the given message to the console or output.
func (log *Logger) Write(msg ...interface{}) {
	str := fmt.Sprint(fmt.Sprint(msg...), "\n")
	if log.output != nil {
		io.WriteString(log.output, str)
	} else {
		fmt.Print(str)
	}
}

// ErrorCount gets the number of errors which have been logged.
func (log *Logger) ErrorCount() int {
	return log.errorCount
}

// Error will log an error to the current output or console.
func (log *Logger) Error(msg ...interface{}) {
	log.Write(msg...)
	log.errorCount++
	if log.showStacks {
		log.Write(log.StackTrace(3, 10, "   "))
	}
}

// WriteAst prints the tree for the given AST object.
func (log *Logger) WriteAst(fileSet *token.FileSet, data interface{}) {
	buf := bytes.Buffer{}
	ast.Fprint(&buf, fileSet, data,
		func(name string, value reflect.Value) bool { return true })
	log.Write(buf.String())
}
