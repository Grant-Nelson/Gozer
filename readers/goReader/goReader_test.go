package transpiler

import "testing"

// Check that Go compilation is being checked by parser.
func TestGoReader_ErrorHandling_001(t *testing.T) {
	test := NewTestGoReader(t)
	test.gr.AddCode("test/main.go",
		`import "fmt"`,
		``,
		`func main() {`,
		`  fmt.Print("Hello World!")`,
		`}`)
	test.CheckErrors(1, Lines(
		`Error: Failed to add source, test/main.go: test/main.go:1:1: expected 'package', found 'import':`,
		`  FilePath: test/main.go`,
		`  Stage:    AddCode`))
}

// Checks the make method call with the wrong number of parameters.
func TestGoReader_ErrorHandling_002(t *testing.T) {
	MainMethodBodyError(t,
		Lines(
			`arr := make([]int, 0, 4, 5)`,
			`fmt.Println("Count = ", len(arr))`),
		1, Lines(
			`Error: Error occurred while processing bodies: Error occurred while parsing a block: Make call must have 1 to 3 arguments but got 4:`,
			`  Mathod: main`,
			`  Path:   test/main.go:4:13`,
			`  Stage:  Processing pending function body`))
}

// Checks the reading a type name.
func TestGoReader_ErrorHandling_003(t *testing.T) {
	MainMethodBodyError(t,
		Lines(
			`arr := make(badType, 4)`,
			`fmt.Println("Count = ", len(arr))`),
		1, Lines(
			`Error: Error occurred while processing bodies: Error occurred while parsing a block: Unhandled type name: badType:`,
			`  Mathod: main`,
			`  Path:   test/main.go:5:15`,
			`  Stage:  Processing pending function body`))
}

// Check of basic main method definition, method selection
// from a package, method call, and literals.
func TestGoReader_Basics_001(t *testing.T) {
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

// Checks the print methods found in the built-in and fmt packages.
func TestGoReader_Basic_002(t *testing.T) {
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

// Checks string literal with back-tick quotes.
func TestGoReader_Literals_001(t *testing.T) {
	MainMethodBodyTest(t,
		Lines("fmt.Print(`Hello World!`)"),
		Lines(`fmt.Print("Hello World!")`))
}

// Checks string literal with back-tick quotes with a newline.
func TestGoReader_Literals_002(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			"fmt.Print(`Hello",
			"World!`)"),
		Lines(`fmt.Print("Hello\n  World!")`))
}

// Checks string literal with unicodes and escape caracters in back-tick quotes.
func TestGoReader_Literals_003(t *testing.T) {
	MainMethodBodyTest(t,
		Lines("fmt.Print(`\\Hello☕!`)"),
		Lines(`fmt.Print("\\Hello\u2615!")`))
}

// Checks string literal with double quotes with single
// quote and escaped double quotes.
func TestGoReader_Literals_004(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(`fmt.Print("'ello \"World\"!")`),
		Lines(`fmt.Print("'ello \"World\"!")`))
}

// Checks string assignment and type determination for the variable definition.
func TestGoReader_Assignment_001(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`msg := "Hello World!"`,
			`fmt.Print(msg)`),
		Lines(
			`string msg = "Hello World!"`,
			`fmt.Print(msg)`))
}

// Checks multiple assignments of string literals.
func TestGoReader_Assignment_002(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`msg, name := "Hello %s!", "World"`,
			`fmt.Printf(msg, name)`),
		Lines(
			`string msg = "Hello %s!"`,
			`string name = "World"`,
			`fmt.Printf(msg, name)`))
}

// Checks multiple assignments of string literals.
func TestGoReader_Assignment_003(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`msg, val := "Answer: %d!", 42`,
			`fmt.Printf(msg, val)`),
		Lines(
			`string msg = "Answer: %d!"`,
			`int val = 42`,
			`fmt.Printf(msg, val)`))
}

// Checks multiple addignment without declaration including a one line swap.
// Also tests generating temporary identifiers.
func TestGoReader_Assignment_004(t *testing.T) {
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
			`int gozerTemp0 = b`,
			`int gozerTemp1 = a`,
			`a = gozerTemp0`,
			`b = gozerTemp1`,
			`fmt.Println("A: ", a, ", B: ", b)`))
}

// Checks assigment from the result of multiple return function.
func TestGoReader_Assignment_005(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`n, err := fmt.Printf("Two = %d\n", 2)`,
			`fmt.Println("n: ", n, ", err: ", err)`),
		Lines(
			`printResult gozerTemp0 = fmt.Printf("Two = %d\n", 2)`,
			`int n = gozerTemp0.n`,
			`error err = gozerTemp0.err`,
			`fmt.Println("n: ", n, ", err: ", err)`))
}

// TODO: Assign to underscore

// Checks concatination of string literals.
func TestGoReader_BinaryOp_001(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(`fmt.Print("Hello "+"World"+"!")`),
		Lines(`fmt.Print((("Hello " + "World") + "!"))`))
}

// Checks multiple assignment and binary operation of integers.
func TestGoReader_BinaryOp_002(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a, b := 1, 4`,
			`fmt.Print("Value: ", a + b)`),
		Lines(
			`int a = 1`,
			`int b = 4`,
			`fmt.Print("Value: ", (a + b))`))
}

// Chesks multiplication of integers, subtration, and order of opertaions.
func TestGoReader_BinaryOp_003(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a, b := 1, 4`,
			`fmt.Print("Value: ", a*0x10 - 1)`),
		Lines(
			`int a = 1`,
			`int b = 4`,
			`fmt.Print("Value: ", ((a * 0x10) - 1))`))
}

// Checks some more binary operations and parentheses for ordering operations.
func TestGoReader_BinaryOp_004(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a, b := 1, 4`,
			`fmt.Print("Value: ", (a*(0x10 - 1))>>b)`),
		Lines(
			`int a = 1`,
			`int b = 4`,
			`fmt.Print("Value: ", ((a * (0x10 - 1)) >> b))`))
}

// Checks the negation operation and negative integer literals.
func TestGoReader_UnaryOp_001(t *testing.T) {
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

// Checks negation and positive unary operations.
func TestGoReader_UnaryOp_002(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a, b := 10, 12`,
			`fmt.Println("Value: ", +a - -b)`),
		Lines(
			`int a = 10`,
			`int b = 12`,
			`fmt.Println("Value: ", (+a - -b))`))
}

// TODO: Add tests for reference and dereference.

// Checks simple if-statement definitions and a comparitive binary operation.
func TestGoReader_IfElse_001(t *testing.T) {
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

// Checks an if-else-statement definition.
func TestGoReader_IfElse_002(t *testing.T) {
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

// Checks assignment within an if-else-statement and that scoping is working for if-else-statements.
func TestGoReader_IfElse_003(t *testing.T) {
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

// Checks else-if-statements in the an if-else-statment.
func TestGoReader_IfElse_004(t *testing.T) {
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

// Checks nested if-else-statments.
func TestGoReader_IfElse_005(t *testing.T) {
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

// TODO: Check assignment of varaibles in the else-if part of an if-else-statement.

// Checks for-statement which increments.
func TestGoReader_For_001(t *testing.T) {
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

// Checks for-statement which decrements.
func TestGoReader_For_002(t *testing.T) {
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

// Checks for-statement with multi-assignment and multiple post operators.
func TestGoReader_For_003(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`for i, j := 9, 0; j < 10; i, j = i-1, i+1 {`,
			`  fmt.Println("up = ", i, ", down = ", j)`,
			`}`),
		Lines(
			`{`,
			`  int i = 9`,
			`  int j = 0`,
			`  void func() gozerTemp2 = void func() {`,
			`    int gozerTemp0 = (i - 1)`,
			`    int gozerTemp1 = (i + 1)`,
			`    i = gozerTemp0`,
			`    j = gozerTemp1`,
			`  }`,
			`  for(; (j < 10); gozerTemp2()) {`,
			`    fmt.Println("up = ", i, ", down = ", j)`,
			`  }`,
			`}`))
}

// Checks for-statement which only the check, a while-statement.
func TestGoReader_For_004(t *testing.T) {
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

// Checks break and continue block statements.
func TestGoReader_For_005(t *testing.T) {
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

// Checks multiple-type of pre and post statements.
func TestGoReader_For_006(t *testing.T) {
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
			`  void func() gozerTemp2 = void func() {`,
			`    int gozerTemp0 = (i + 1)`,
			`    bool gozerTemp1 = false`,
			`    i = gozerTemp0`,
			`    first = gozerTemp1`,
			`  }`,
			`  for(; (i < 10); gozerTemp2()) {`,
			`    if first {`,
			`      fmt.Println("First: ", i)`,
			`      continue`,
			`    }`,
			`    fmt.Println("Rest: ", i)`,
			`  }`,
			`}`))
}

// Checks foreach on list with only index.
func TestGoReader_ForListRange_001(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a := []int{1, 2, 3, 4}`,
			`for i := range a {`,
			`  fmt.Println("Index: ", i)`,
			`}`),
		Lines(
			`[]int a = []int{1, 2, 3, 4}`,
			`for(int i = 0; (i < len(a)); i++) {`,
			`  fmt.Println("Index: ", i)`,
			`}`))
}

// Checks foreach on list with index and value.
func TestGoReader_ForListRange_002(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a := []int{1, 2, 3, 4}`,
			`for i, val := range a {`,
			`  fmt.Println("Index: ", i, ", Value: ", val)`,
			`}`),
		Lines(
			`[]int a = []int{1, 2, 3, 4}`,
			`for(int i = 0; (i < len(a)); i++) {`,
			`  int val = a[i]`,
			`  fmt.Println("Index: ", i, ", Value: ", val)`,
			`}`))
}

// Checks foreach with predefined value and index on list.
func TestGoReader_ForListRange_003(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a := []int{1, 2, 3, 4}`,
			`val, i := -1, -2`,
			`for i, val = range a {`,
			`  fmt.Println("Index: ", i, ", Value: ", val)`,
			`}`),
		Lines(
			`[]int a = []int{1, 2, 3, 4}`,
			`int val = -1`,
			`int i = -2`,
			`for(i = 0; (i < len(a)); i++) {`,
			`  val = a[i]`,
			`  fmt.Println("Index: ", i, ", Value: ", val)`,
			`}`))
}

// Checks foreach on a list with no iterators.
func TestGoReader_ForListRange_004(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a := []int{1, 2, 3, 4}`,
			`for range a {`,
			`  fmt.Println("Bleep")`,
			`}`),
		Lines(
			`[]int a = []int{1, 2, 3, 4}`,
			`for(int gozerTemp0 = 0; (gozerTemp0 < len(a)); gozerTemp0++) {`,
			`  fmt.Println("Bleep")`,
			`}`))
}

// Checks foreach with predefined value and no index on list.
func TestGoReader_ForListRange_005(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`a := []int{1, 2, 3, 4}`,
			`for _, val := range a {`,
			`  fmt.Println("Value: ", val)`,
			`}`),
		Lines(
			`[]int a = []int{1, 2, 3, 4}`,
			`for(int gozerTemp0 = 0; (gozerTemp0 < len(a)); gozerTemp0++) {`,
			`  int val = a[gozerTemp0]`,
			`  fmt.Println("Value: ", val)`,
			`}`))
}

// Checks foreach with predefined value and index on unidentified list.
func TestGoReader_ForListRange_006(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`for i, val := range []int{1, 2, 3, 4} {`,
			`  fmt.Println("Index: ", i, ", Value: ", val)`,
			`}`),
		Lines(
			`{`,
			`  []int gozerTemp0 = []int{1, 2, 3, 4}`,
			`  for(int i = 0; (i < len(gozerTemp0)); i++) {`,
			`    int val = gozerTemp0[i]`,
			`    fmt.Println("Index: ", i, ", Value: ", val)`,
			`  }`,
			`}`))
}

// Checks foreach with predefined value and no index on unidentified list.
func TestGoReader_ForListRange_007(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`for _, val := range []int{1, 2, 3, 4} {`,
			`  fmt.Println("Value: ", val)`,
			`}`),
		Lines(
			`{`,
			`  []int gozerTemp0 = []int{1, 2, 3, 4}`,
			`  for(int gozerTemp1 = 0; (gozerTemp1 < len(gozerTemp0)); gozerTemp1++) {`,
			`    int val = gozerTemp0[gozerTemp1]`,
			`    fmt.Println("Value: ", val)`,
			`  }`,
			`}`))
}

// Checks foreach with predefined no value and no index on unidentified list.
func TestGoReader_ForListRange_008(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`for _, _ := range []int{1, 2, 3, 4} {`,
			`  fmt.Println("Bleep")`,
			`}`),
		Lines(
			`{`,
			`  []int gozerTemp0 = []int{1, 2, 3, 4}`,
			`  for(int gozerTemp1 = 0; (gozerTemp1 < len(gozerTemp0)); gozerTemp1++) {`,
			`    fmt.Println("Bleep")`,
			`  }`,
			`}`))
}

// TODO: For-each map

// Checks the creation of a slicee and the build-in length method.
func TestGoReader_Slices_001(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`arr := []int{4, 1, 3, 2}`,
			`fmt.Println("Count = ", len(arr))`),
		Lines(
			`[]int arr = []int{4, 1, 3, 2}`,
			`fmt.Println("Count = ", len(arr))`))
}

// Checks the indices of a slice expression.
func TestGoReader_Slices_002(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`arr := []int{4, 1, 3, 2}`,
			`for i := 0; i < len(arr); i++ {`,
			`  fmt.Println(i, ": ", arr[i])`,
			`}`),
		Lines(
			`[]int arr = []int{4, 1, 3, 2}`,
			`for(int i = 0; (i < len(arr)); i++) {`,
			`  fmt.Println(i, ": ", arr[i])`,
			`}`))
}

// Checks assignment of a slice element with an index.
func TestGoReader_Slices_003(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`arr := []int{4, 1, 3, 2}`,
			`arr[2] = 8`,
			`fmt.Println("arr[2] = ", arr[2])`),
		Lines(
			`[]int arr = []int{4, 1, 3, 2}`,
			`arr[2] = 8`,
			`fmt.Println("arr[2] = ", arr[2])`))
}

// Checks appending to a slice.
func TestGoReader_Slices_004(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`arr := []int{4, 1, 3, 2}`,
			`arr = append(arr, 8)`,
			`fmt.Println("Count = ", len(arr))`),
		Lines(
			`[]int arr = []int{4, 1, 3, 2}`,
			`arr = append(arr, 8)`,
			`fmt.Println("Count = ", len(arr))`))
}

// Checks the creation of a slice with a default length
// and a check of the build-in capacity method.
func TestGoReader_Slices_005(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`arr := make([]int, 4)`,
			`arr[2] = 8`,
			`fmt.Println("arr[2] = ", arr[2])`,
			`fmt.Println("Count = ", len(arr))`,
			`fmt.Println("Cap = ", cap(arr))`),
		Lines(
			`[]int arr = make([]int, 4)`,
			`arr[2] = 8`,
			`fmt.Println("arr[2] = ", arr[2])`,
			`fmt.Println("Count = ", len(arr))`,
			`fmt.Println("Cap = ", cap(arr))`))
}

// Checks the creation of a slice with a default capacity.
func TestGoReader_Slices_006(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`arr := make([]int, 0, 4)`,
			`arr = append(arr, 8)`,
			`fmt.Println("Count = ", len(arr))`,
			`fmt.Println("Cap = ", cap(arr))`),
		Lines(
			`[]int arr = make([]int, 0, 4)`,
			`arr = append(arr, 8)`,
			`fmt.Println("Count = ", len(arr))`,
			`fmt.Println("Cap = ", cap(arr))`))
}

// Checks creating a subslice and assigning to that subslice.
func TestGoReader_Slices_007(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`arr := []int{4, 1, 3, 2}`,
			`arr2 := arr[1:2]`,
			`arr2[0] = 8`,
			`arr[2] = 7`,
			`fmt.Printf("arr = %v", arr)`,
			`fmt.Printf("arr2 = %v", arr2)`),
		Lines(
			`[]int arr = []int{4, 1, 3, 2}`,
			`[]int arr2 = arr[1:2]`,
			`arr2[0] = 8`,
			`arr[2] = 7`,
			`fmt.Printf("arr = %v", arr)`,
			`fmt.Printf("arr2 = %v", arr2)`))
}

// Checks creating a subslice from the beginning to an index.
func TestGoReader_Slices_008(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`arr := []int{4, 1, 3, 2}`,
			`arr2 := arr[:2]`,
			`arr2[0] = 8`,
			`arr[2] = 7`,
			`fmt.Printf("arr = %v", arr)`,
			`fmt.Printf("arr2 = %v", arr2)`),
		Lines(
			`[]int arr = []int{4, 1, 3, 2}`,
			`[]int arr2 = arr[:2]`,
			`arr2[0] = 8`,
			`arr[2] = 7`,
			`fmt.Printf("arr = %v", arr)`,
			`fmt.Printf("arr2 = %v", arr2)`))
}

// Checks creating a subslice from an index to the end.
func TestGoReader_Slices_009(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`arr := []int{4, 1, 3, 2}`,
			`arr2 := arr[1:]`,
			`arr2[0] = 8`,
			`arr[2] = 7`,
			`fmt.Printf("arr = %v", arr)`,
			`fmt.Printf("arr2 = %v", arr2)`),
		Lines(
			`[]int arr = []int{4, 1, 3, 2}`,
			`[]int arr2 = arr[1:]`,
			`arr2[0] = 8`,
			`arr[2] = 7`,
			`fmt.Printf("arr = %v", arr)`,
			`fmt.Printf("arr2 = %v", arr2)`))
}

// Checks creating a subslice from an index to the end.
func TestGoReader_Slices_010(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`arr := []int{4, 1, 3, 2}`,
			`arr2 := arr[1:2:3]`,
			`arr2[0] = 8`,
			`arr[2] = 7`,
			`fmt.Printf("arr = %v", arr)`,
			`fmt.Printf("arr2 = %v", arr2)`),
		Lines(
			`[]int arr = []int{4, 1, 3, 2}`,
			`[]int arr2 = arr[1:2:3]`,
			`arr2[0] = 8`,
			`arr[2] = 7`,
			`fmt.Printf("arr = %v", arr)`,
			`fmt.Printf("arr2 = %v", arr2)`))
}

// Checks initialization of a map.
func TestGoReader_Maps_001(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`numbers := map[int]string{4: "four", 1: "one", 3: "three", 2: "two"}`,
			`fmt.Println("Count = ", len(numbers))`),
		Lines(
			`map[int]string numbers = map[int]string{4: "four", 1: "one", 3: "three", 2: "two"}`,
			`fmt.Println("Count = ", len(numbers))`))
}

// Checks assignment of a map element with an int key.
func TestGoReader_Maps_002(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`numbers := map[int]string{4: "four", 1: "one", 3: "three", 2: "two"}`,
			`numbers[2] = "duo"`,
			`fmt.Println("numbers[2] = ", numbers[2])`),
		Lines(
			`map[int]string numbers = map[int]string{4: "four", 1: "one", 3: "three", 2: "two"}`,
			`numbers[2] = "duo"`,
			`fmt.Println("numbers[2] = ", numbers[2])`))
}

// Checks assignment of a map element with an string key.
func TestGoReader_Maps_003(t *testing.T) {
	MainMethodBodyTest(t,
		Lines(
			`numbers := map[string]int{"four": 4, "one": 1, "three": 3, "two": 2}`,
			`numbers["two"] = 42`,
			`fmt.Println("numbers[two] = ", numbers["two"])`),
		Lines(
			`map[string]int numbers = map[string]int{"four": 4, "one": 1, "three": 3, "two": 2}`,
			`numbers["two"] = 42`,
			`fmt.Println("numbers[two] = ", numbers["two"])`))
}

// TODO: delete map

// TODO: Switch, Class, Struct defs, etc

/*
// TODO: complex for-each indices
package main

import (
	"fmt"
)

type A struct {
	i int
	v float64
}

func (a A) String() string {
	return fmt.Sprint("<", a.i,", ", a.v, ">")
}

func main() {
	a := A{i: -1, v: -1.0}
	for a.i, a.v = range []float64{0.12, 0.23, 0.34, 0.45} {
		fmt.Print("A = ", a, "\n")
	}
}
*/

/*
	Need to make the following fail because of types not matching
	`numbers := map[int]string{4: "four", 1: "one", 3: "three", 2: "two"}`,
	`numbers[2] = 8`,
*/
