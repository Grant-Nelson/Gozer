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
			`string msg = "Hello World!"`,
			`fmt.Print(msg)`))
}

func TestGoReader9(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`msg, name := "Hello %s!", "World"`,
			`fmt.Printf(msg, name)`),
		Lines(
			`string msg = "Hello %s!"`,
			`string name = "World"`,
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
			`int a = 1`,
			`int b = 4`,
			`fmt.Print("Value: ", (a + b))`))
}

func TestGoReader12(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a, b := 1, 4`,
			`fmt.Print("Value: ", a*0x10 - 1)`),
		Lines(
			`int a = 1`,
			`int b = 4`,
			`fmt.Print("Value: ", ((a * 0x10) - 1))`))
}

func TestGoReader13(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a, b := 1, 4`,
			`fmt.Print("Value: ", (a*(0x10 - 1))>>b)`),
		Lines(
			`int a = 1`,
			`int b = 4`,
			`fmt.Print("Value: ", ((a * (0x10 - 1)) >> b))`))
}

func TestGoReader14(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a := 14`,
			`fmt.Print("Value 1: ", a)`,
			`a = -8`,
			`fmt.Print("Value 2: ", -a)`),
		Lines(
			`int a = 14`,
			`fmt.Print("Value 1: ", a)`,
			`a = -8`,
			`fmt.Print("Value 2: ", -a)`))
}

func TestGoReader15(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a, b := 10, 12`,
			`fmt.Println("Value: ", +a - -b)`),
		Lines(
			`int a = 10`,
			`int b = 12`,
			`fmt.Println("Value: ", (+a - -b))`))
}

func TestGoReader16(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a, b := 10, 12`,
			`fmt.Println("A: ", a, ", B: ", b)`,
			`a, b = b, a`, // Swap
			`fmt.Println("A: ", a, ", B: ", b)`),
		Lines(
			`int a = 10`,
			`int b = 12`,
			`fmt.Println("A: ", a, ", B: ", b)`,
			`int temp0 = b`,
			`int temp1 = a`,
			`a = temp0`,
			`b = temp1`,
			`fmt.Println("A: ", a, ", B: ", b)`))
}

func TestGoReader17(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`n, err := fmt.Printf("Two = %d\n", 2)`,
			`fmt.Println("n: ", n, ", err: ", err)`),
		Lines(
			`int n`,
			`error err`,
			`n, err = fmt.Printf("Two = %d\n", 2)`,
			`fmt.Println("n: ", n, ", err: ", err)`))
}

func TestGoReader18(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a := 10 `,
			`if a > 4 {`,
			`  fmt.Println("greater: ", a)`,
			`}`),
		Lines(
			`int a = 10`,
			`if (a > 4) {`,
			`  fmt.Println("greater: ", a)`,
			`}`))
}

func TestGoReader19(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a := 10 `,
			`if a > 4 {`,
			`  fmt.Println("greater: ", a)`,
			`} else {`,
			`  fmt.Println("less: ", a)`,
			`}`),
		Lines(
			`int a = 10`,
			`if (a > 4) {`,
			`  fmt.Println("greater: ", a)`,
			`} else {`,
			`  fmt.Println("less: ", a)`,
			`}`))
}

func TestGoReader20(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a := 2`,
			`if b, c := 10, 4; b > c + a {`,
			`  fmt.Println("greater: ", a, " > ", c + a)`,
			`} else {`,
			`  fmt.Println("less: ", a, " > ", c + a)`,
			`}`),
		Lines(
			`int a = 2`,
			`{`,
			`  int b = 10`,
			`  int c = 4`,
			`  if (b > (c + a)) {`,
			`    fmt.Println("greater: ", a, " > ", (c + a))`,
			`  } else {`,
			`    fmt.Println("less: ", a, " > ", (c + a))`,
			`  }`,
			`}`))
}

func TestGoReader21(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a := 5`,
			`if a > 10 {`,
			`  fmt.Println(a, " > 10")`,
			`} else if a > 8 {`,
			`  fmt.Println(a, " > 8")`,
			`} else if a > 0 {`,
			`  fmt.Println(a, " > 0")`,
			`} else {`,
			`  fmt.Println(a, " <= 0")`,
			`}`),
		Lines(
			`int a = 5`,
			`if (a > 10) {`,
			`  fmt.Println(a, " > 10")`,
			`} else if (a > 8) {`,
			`  fmt.Println(a, " > 8")`,
			`} else if (a > 0) {`,
			`  fmt.Println(a, " > 0")`,
			`} else {`,
			`  fmt.Println(a, " <= 0")`,
			`}`))
}

func TestGoReader22(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a := 5`,
			`if a > 0 {`,
			`  if a > 10 {`,
			`    fmt.Println(a, " > 10")`,
			`  } else {`,
			`    fmt.Println(a, " > 0")`,
			`  }`,
			`} else {`,
			`  if a < -10 {`,
			`    fmt.Println(a, " < -10")`,
			`  } else {`,
			`    fmt.Println(a, " <= 0")`,
			`  }`,
			`}`),
		Lines(
			`int a = 5`,
			`if (a > 0) {`,
			`  if (a > 10) {`,
			`    fmt.Println(a, " > 10")`,
			`  } else {`,
			`    fmt.Println(a, " > 0")`,
			`  }`,
			`} else {`,
			`  if (a < -10) {`,
			`    fmt.Println(a, " < -10")`,
			`  } else {`,
			`    fmt.Println(a, " <= 0")`,
			`  }`,
			`}`))
}

func TestGoReader23(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`for i := 0; i < 10; i++ {`,
			`  fmt.Println("index = ", i)`,
			`}`),
		Lines(
			`for(int i = 0; (i < 10); i++) {`,
			`  fmt.Println("index = ", i)`,
			`}`))
}

func TestGoReader24(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`for i := 9; i >= 0; i-- {`,
			`  fmt.Println("index = ", i)`,
			`}`),
		Lines(
			`for(int i = 9; (i >= 0); i--) {`,
			`  fmt.Println("index = ", i)`,
			`}`))
}

func TestGoReader25(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`for i, j := 9, 0; j < 10; i, j = i-1, i+1 {`,
			`  fmt.Println("up = ", i, ", down = ", j)`,
			`}`),
		Lines(
			`{`,
			`  int i = 9`,
			`  int j = 0`,
			`  for(; (j < 10); {int temp0 = (i - 1); int temp1 = (i + 1); i = temp0; j = temp1}) {`,
			`    fmt.Println("up = ", i, ", down = ", j)`,
			`  }`,
			`}`))
}

func TestGoReader26(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`i := 0`,
			`for i < 10 {`,
			`  fmt.Println("index = ", i)`,
			`  i++`,
			`}`),
		Lines(
			`int i = 0`,
			`for(; (i < 10); ) {`,
			`  fmt.Println("index = ", i)`,
			`  i++`,
			`}`))
}

func TestGoReader27(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`for i := 0; true; i++ {`,
			`  if i < 5 {`,
			`    fmt.Println("small: ", i)`,
			`    continue`,
			`  }`,
			`  fmt.Println("large: ", i)`,
			`  if i > 10 {`,
			`    break`,
			`  }`,
			`}`),
		Lines(
			`for(int i = 0; true; i++) {`,
			`  if (i < 5) {`,
			`    fmt.Println("small: ", i)`,
			`    continue`,
			`  }`,
			`  fmt.Println("large: ", i)`,
			`  if (i > 10) {`,
			`    break`,
			`  }`,
			`}`))
}

func TestGoReader28(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`for i, first := 0, true; i < 10; i, first = i+1, false {`,
			`  if first {`,
			`    fmt.Println("First: ", i)`,
			`    continue`,
			`  }`,
			`  fmt.Println("Rest: ", i)`,
			`}`),
		Lines(
			`{`,
			`  int i = 0`,
			`  bool first = true`,
			`  for(; (i < 10); {int temp0 = (i + 1); bool temp1 = false; i = temp0; first = temp1}) {`,
			`    if first {`,
			`      fmt.Println("First: ", i)`,
			`      continue`,
			`    }`,
			`    fmt.Println("Rest: ", i)`,
			`  }`,
			`}`))
}
