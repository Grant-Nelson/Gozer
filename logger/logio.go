package common

import (
	"fmt"
	"go/token"
	"io"
)

var _ Logger = (*LogIO)(nil)

// LogIO handles logging messages for the transpiler.
type LogIO struct {

	// output is the writer to print to or nil to print to the console.
	output io.Writer

	// errorCount is the number of errors which have been logged.
	errorCount int

	// showStacks indicates if error should also print stacktraces.
	showStacks bool

	// showDebug indicates debugging logs should be printed.
	showDebug bool
}

// NewLogIO creates a new logger to a buffer.
func NewLogIO() *LogIO {
	return &LogIO{
		output:     nil,
		errorCount: 0,
		showStacks: false,
		showDebug:  false,
	}
}

// SetOutput sets the output for the log.
// If the output is set to nil then the log is printed to the console.
func (log *LogIO) SetOutput(output io.Writer) {
	log.output = output
}

// ShowStackTraces indicates if stacktraces should be shown with errors.
func (log *LogIO) ShowStackTraces(showStacks bool) {
	log.showStacks = showStacks
}

// ShowDebug indicates if debugging messages should be shown.
func (log *LogIO) ShowDebug(showDebug bool) {
	log.showDebug = showDebug
}

// ErrorCount gets the number of errors which have been logged.
func (log *LogIO) ErrorCount() int {
	return log.errorCount
}

// Write prints the given message to the console or output.
func (log *LogIO) Write(msg ...interface{}) {
	str := fmt.Sprint(fmt.Sprint(msg...), "\n")
	if log.output != nil {
		io.WriteString(log.output, str)
	} else {
		fmt.Print(str)
	}
}

// Debug will log debugging information to the current output or console.
func (log *LogIO) Debug(msg ...interface{}) {
	if log.showDebug {
		log.Write(msg...)
	}
}

// Error will log an error to the current output or console.
func (log *LogIO) Error(msg ...interface{}) {
	log.Write(msg...)
	log.errorCount++
	if log.showStacks {
		log.Write(StackTrace(3, 10, "   "))
	}
}

// WriteAst prints the tree for the given AST object.
func (log *LogIO) WriteAst(fileSet *token.FileSet, data interface{}) {
	log.Write(AstString(fileSet, data))
}
