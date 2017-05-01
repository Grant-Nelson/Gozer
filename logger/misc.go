package common

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"runtime/debug"
	"strings"
)

// MessageProcessor is an interface for processing and modifying a message.
type MessageProcessor interface {

	// Process handles and/or modifies the given message
	// and returns the same or new message.
	// To remove/stop the message return nil.
	Process(msg *Message) *Message
}

// ThrowError throws an error via a panic.
func ThrowError(msg ...interface{}) {
	panic(fmt.Errorf("%s", fmt.Sprint(msg...)))
}

// StackTrace gets the current stack strace.
func StackTrace(scope int, steps int, Indent string) string {
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")
	scopeCut := scope*2 + 1
	topCut := scopeCut + steps*2
	if len(lines) > topCut {
		lines = lines[scopeCut:topCut]
	}
	return Indent + strings.Join(lines, "\n"+Indent)
}

// AstString gets the string of the tree for the given AST object.
func AstString(fileSet *token.FileSet, data interface{}) string {
	buf := bytes.Buffer{}
	ast.Fprint(&buf, fileSet, data,
		func(name string, value reflect.Value) bool { return true })
	return buf.String()
}
