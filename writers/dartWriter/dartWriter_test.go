package dartWriter

import (
	"testing"

	"github.com/grant-nelson/Gozer/tests"
)

// TestRunTestSuite runs the suite of tests
func TestRunTestSuite(t *testing.T) {
	tests.RunTestSuite(t, func(t2 *testing.T, code map[tests.SnippetType]string) {
		test := NewTestGoReader(t2)
		test.AddCode("test/main.go", code[tests.GoSnippet])
		test.Transpile()
		test.CheckPackage("test", code[tests.BlockSnippet])
	})
}

// TestBasics001 check of basic main method definition, method selection
// from a package, method call, and literals.
func TestBasics001(t *testing.T) {
	test := NewTestGoReader(t)
	test.AddCode("test/main.go",
		`package main`,
		`import "fmt"`,
		``,
		`// This is a test of a printing a normal string.`,
		`func main() {`,
		`  fmt.Print("Hello World!")`,
		`}`)
	test.CheckPackages("builtin", "fmt", "io", "test")
	test.CheckImports("test", "builtin", "fmt")
	test.Transpile()
	test.CheckFunctions("test", "main")
	test.CheckFunction("test", "main",
		`void main() {`,
		`  fmt.Print("Hello World!")`,
		`}`)
}

// TestErrorHandling001 check that Go compilation is being checked by parser.
func TestErrorHandling001(t *testing.T) {
	test := NewTestGoReader(t)
	test.gr.AddCode("test/main.go",
		`import "fmt"`,
		``,
		`func main() {`,
		`  fmt.Print("Hello World!")`,
		`}`)
	test.CheckErrors(1, tests.Lines(
		`Error: Failed to add source, test/main.go: test/main.go:1:1: expected 'package', found 'import':`,
		`  FilePath: test/main.go`,
		`  Stage:    AddCode`))
}

// TestErrorHandling002 checks the make method call with the wrong number of parameters.
func TestErrorHandling002(t *testing.T) {
	MainMethodBodyError(t,
		tests.Lines(
			`arr := make([]int, 0, 4, 5)`,
			`fmt.Println("Count = ", len(arr))`),
		1, tests.Lines(
			`Error: Error occurred while processing bodies: Error occurred while parsing a block: Make call must have 1 to 3 arguments but got 4:`,
			`  Mathod: main`,
			`  Path:   test/main.go:4:13`,
			`  Stage:  Processing pending function body`))
}

// TestErrorHandling003 checks the reading a type name.
func TestErrorHandling003(t *testing.T) {
	MainMethodBodyError(t,
		tests.Lines(
			`arr := make(badType, 4)`,
			`fmt.Println("Count = ", len(arr))`),
		1, tests.Lines(
			`Error: Error occurred while processing bodies: Error occurred while parsing a block: Unhandled type name: badType:`,
			`  Mathod: main`,
			`  Path:   test/main.go:5:15`,
			`  Stage:  Processing pending function body`))
}
