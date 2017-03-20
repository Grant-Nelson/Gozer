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
	test.CheckFunction("test", "main",
		`func main() {`,
		`  fmt.Print("Hello World!")`,
		`}`)
}

func TestGoReader3(t *testing.T) {
	MainMethodBodyTest(t,
		Lines("fmt.Print(`Hello World!`)"),
		Lines(`fmt.Print("Hello World!")`))
}

func TestGoReader4(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			"fmt.Print(`Hello",
			"World!`)"),
		Lines(`fmt.Print("Hello\n  World!")`))
}

func TestGoReader5(t *testing.T) {
	MainMethodBodyTest(t,
		Lines("fmt.Print(`\\Hello☕!`)"),
		Lines(`fmt.Print("\\Hello\u2615!")`))
}

func TestGoReader7(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`print("Hello World!\n")`,
			`println("Hello World!\n")`,
			`fmt.Printf("Hello %s\n", "World!")`,
			`fmt.Println("Hello World!")`),
		Lines(
			`print("Hello World!\n")`,
			`println("Hello World!\n")`,
			`fmt.Printf("Hello %s\n", "World!")`,
			`fmt.Println("Hello World!")`))
}

func TestGoReader8(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`msg := "Hello World!"`,
			`fmt.Print(msg)`),
		Lines(
			`msg = "Hello World!"`,
			`fmt.Print(msg)`))
}

func TestGoReader9(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`msg, name := "Hello %s!", "World"`,
			`fmt.Printf(msg, name)`),
		Lines(
			`msg = "Hello %s!"`,
			`name = "World"`,
			`fmt.Printf(msg, name)`))
}

func TestGoReader10(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(`fmt.Print("Hello "+"World"+"!")`),
		Lines(`fmt.Print((("Hello " + "World") + "!"))`))
}

func TestGoReader11(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a, b := 1, 4`,
			`fmt.Print("Value: ", a + b)`),
		Lines(
			`a = 1`,
			`b = 4`,
			`fmt.Print("Value: ", (a + b))`))
}

func TestGoReader12(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a, b := 1, 4`,
			`fmt.Print("Value: ", a*0x10 - 1)`),
		Lines(
			`a = 1`,
			`b = 4`,
			`fmt.Print("Value: ", ((a * 0x10) - 1))`))
}

func TestGoReader13(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a, b := 1, 4`,
			`fmt.Print("Value: ", (a*(0x10 - 1))>>b)`),
		Lines(
			`a = 1`,
			`b = 4`,
			`fmt.Print("Value: ", ((a * (0x10 - 1)) >> b))`))
}
