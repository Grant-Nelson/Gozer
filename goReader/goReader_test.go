package transpiler

import "testing"

func TestGoReader1(t *testing.T) {
	test := NewTestGoReader(t)
	test.gr.AddCode("test/main.go",
		`import "fmt"`,
		``,
		`func main() {`,
		`  fmt.Print("Hello World!")`,
		`}`)
	test.CheckErrors(1,
		`Failed to add source, test/main.go: test/main.go:1:1: expected 'package', found 'import'`)
}

func TestGoReader2(t *testing.T) {
	test := NewTestGoReader(t)
	test.AddCode("test/main.go",
		`package main`,
		`import "fmt"`,
		``,
		`func main() {`,
		`  fmt.Print("Hello World!")`,
		`}`)
	test.CheckPackages("builtin", "fmt", "io", "test")
	test.CheckImports("test", "builtin", "fmt")
	test.Transpile()
}
