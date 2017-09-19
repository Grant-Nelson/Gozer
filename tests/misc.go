package tests

import (
	"strings"
	"testing"

	"github.com/grant-nelson/Gozer/common"
)

// SnippetType is the type of the code snippet.
type SnippetType int

const (
	// GoSnippet is a bit of Go code
	GoSnippet SnippetType = iota

	// BlockSnippet is a bit of construct blocks' output
	BlockSnippet

	// DartSnippet is a bit of Dart code
	DartSnippet
)

// SuiteHandler is the handler called for each value of the suite.
type SuiteHandler func(t *testing.T, snippets map[SnippetType]string)

// Lines concatinates several lines into a single string.
func Lines(code ...string) string {
	return strings.Join(code, "\n")
}

// GoMainMethodBody constructs a main file with only a main method,
// the body of the main method is the provided code.s
func GoMainMethodBody(code ...string) string {
	return Lines(
		`package main`,
		`import "fmt"`,
		``,
		`func main() {`,
		`  `+common.Indent(Lines(code...), "  "),
		`}`)
}

// BlockMainMethodBody constructs a main file with only a main method,
// the body of the main method is the provided code.
func BlockMainMethodBody(code ...string) string {
	return Lines(
		`import test{`,
		`  import builtin`,
		`  import fmt`,
		`  void main() {`,
		`    `+common.Indent(Lines(code...), "    "),
		`  }`,
		`}`)
}

// DartMainMethodBody constructs a main file with only a main method,
// the body of the main method is the provided code.
func DartMainMethodBody(code ...string) string {
	return Lines(
		`import test{`,
		`  import builtin`,
		`  import fmt`,
		`  void main() {`,
		`    `+common.Indent(Lines(code...), "    "),
		`  }`,
		`}`)
}
