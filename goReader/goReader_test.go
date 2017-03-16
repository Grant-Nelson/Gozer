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
		`// This is a test of a printing a normal string.`,
		`func main() {`,
		`  fmt.Print("Hello World!")`,
		`}`)
	test.CheckPackages("builtin", "fmt", "io", "test")
	test.CheckImports("test", "builtin", "fmt")
	test.Transpile()
	test.CheckFunctions("test", "main")
	test.CheckFunctionBody("test", "main",
		`func main() {`,
		`  fmt.Print("Hello World!")`,
		`}`)
}

func TestGoReader3(t *testing.T) {
	test := NewTestGoReader(t)
	test.AddCode("test/main.go",
		`package main`,
		`import "fmt"`,
		``,
		`// A print statement using a back ticks`,
		`// instead of the double quote.`,
		`func main() {`,
		"  fmt.Print(`Hello World!`)",
		`}`)
	test.Transpile()
	test.CheckFunctionBody("test", "main",
		`func main() {`,
		`  fmt.Print("Hello World!")`,
		`}`)
}

func TestGoReader4(t *testing.T) {
	test := NewTestGoReader(t)
	test.AddCode("test/main.go",
		`package main`,
		`import "fmt"`,
		``,
		`/**`,
		` * Test of a multi-line tick string.`,
		` */`,
		`func main() {`,
		"  fmt.Print(`Hello",
		"World!`)",
		`}`)
	test.Transpile()
	test.CheckFunctionBody("test", "main",
		`func main() {`,
		`  fmt.Print("Hello\nWorld!")`,
		`}`)
}

func TestGoReader5(t *testing.T) {
	test := NewTestGoReader(t)
	test.AddCode("test/main.go",
		`package main`,
		`import "fmt"`,
		`func main() {`,
		"  fmt.Print(`\\Hello☕!`)",
		`}`)
	test.Transpile()
	test.CheckFunctionBody("test", "main",
		`func main() {`,
		`  fmt.Print("\\Hello\u2615!")`,
		`}`)
}

func TestGoReader7(t *testing.T) {
	test := NewTestGoReader(t)
	test.AddCode("test/main.go",
		`package main`,
		`import "fmt"`,
		``,
		`func main() {`,
		`  print("Hello World!\n")`,
		`  println("Hello World!\n")`,
		`  fmt.Printf("Hello %s\n", "World!")`,
		`  fmt.Println("Hello World!")`,
		`}`)
	test.Transpile()
	test.CheckFunctionBody("test", "main",
		`func main() {`,
		`  print("Hello World!\n")`,
		`  println("Hello World!\n")`,
		`  fmt.Printf("Hello %s\n", "World!")`,
		`  fmt.Println("Hello World!")`,
		`}`)
}

func TestGoReader8(t *testing.T) {
	test := NewTestGoReader(t)
	test.AddCode("test/main.go",
		`package main`,
		`import "fmt"`,
		`func main() {`,
		`  msg := "Hello World!"`,
		`  fmt.Print(msg)`,
		`}`)
	test.Transpile()
	test.CheckFunctionBody("test", "main",
		`func main() {`,
		`  fmt.Print("Hello World!")`,
		`}`)
}
