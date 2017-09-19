package tests

import (
	"testing"
)

// RunTestSuite runs all the common tests for a given configuration.
func RunTestSuite(t *testing.T, hndl SuiteHandler) {

	// Basic001 checks the print methods found in the built-in and fmt packages.
	t.Run("Basic001", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`print("Hello World!\n")`,
				`println("Hello World!")`,
				`fmt.Printf("Hello %s\n", "World!")`,
				`fmt.Println("Hello World!")`),
			BlockSnippet: BlockMainMethodBody(
				`print("Hello World!\n")`,
				`println("Hello World!")`,
				`fmt.Printf("Hello %s\n", "World!")`,
				`fmt.Println("Hello World!")`),
			DartSnippet: DartMainMethodBody(
				`builtin.print('Hello World!\n');`,
				`builtin.println('Hello World!');`,
				`fmt.printf('Hello %s\n', 'World!');`,
				`fmt.println('Hello World!');`),
		})
	})

	// Literals001 checks string literal with back-tick quotes.
	t.Run("Literals001", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet:    GoMainMethodBody("fmt.Print(`Hello World!`)"),
			BlockSnippet: BlockMainMethodBody(`fmt.Print("Hello World!")`),
			DartSnippet:  BlockMainMethodBody(`fmt.print('Hello World!');`),
		})
	})

	// Literals002 checks string literal with back-tick quotes with a newline.
	t.Run("Literals002", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				"fmt.Print(`Hello",
				"World!`)"),
			BlockSnippet: BlockMainMethodBody(`fmt.Print("Hello\n  World!")`),
			DartSnippet:  BlockMainMethodBody(`fmt.print('Hello\n  World!');`),
		})
	})

	// Literals003 checks string literal with unicodes and escape caracters in back-tick quotes.
	t.Run("Literals003", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet:    GoMainMethodBody("fmt.Print(`\\Hello☕!`)"),
			BlockSnippet: BlockMainMethodBody(`fmt.Print("\\Hello\u2615!")`),
			DartSnippet:  BlockMainMethodBody(`fmt.orint('\\Hello\u2615!');`),
		})
	})

	// Literals004 checks string literal with double quotes with single
	// quote and escaped double quotes.
	t.Run("Literals004", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet:    GoMainMethodBody(`fmt.Print("'ello \"World\"!")`),
			BlockSnippet: BlockMainMethodBody(`fmt.Print("'ello \"World\"!")`),
			DartSnippet:  BlockMainMethodBody(`fmt.print('\'ello "World"!');`),
		})
	})

	// Assignment001 checks string assignment and type determination for the variable definition.
	t.Run("Assignment001", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`msg := "Hello World!"`,
				`fmt.Print(msg)`),
			BlockSnippet: BlockMainMethodBody(
				`string msg = "Hello World!"`,
				`fmt.Print(msg)`),
		})
	})

	// Assignment002 checks multiple assignments of string literals.
	t.Run("Assignment002", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`msg, name := "Hello %s!", "World"`,
				`fmt.Printf(msg, name)`),
			BlockSnippet: BlockMainMethodBody(
				`string msg = "Hello %s!"`,
				`string name = "World"`,
				`fmt.Printf(msg, name)`),
		})
	})

	// Assignment003 checks multiple assignments of string literals.
	t.Run("Assignment003", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`msg, val := "Answer: %d!", 42`,
				`fmt.Printf(msg, val)`),
			BlockSnippet: BlockMainMethodBody(
				`string msg = "Answer: %d!"`,
				`int val = 42`,
				`fmt.Printf(msg, val)`),
		})
	})

	// Assignment004 checks multiple addignment without declaration including a one line swap.
	// Also tests generating temporary identifiers.
	t.Run("Assignment004", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a, b := 10, 12`,
				`fmt.Println("A: ", a, ", B: ", b)`,
				`a, b = b, a`, // Swap
				`fmt.Println("A: ", a, ", B: ", b)`),
			BlockSnippet: BlockMainMethodBody(
				`int a = 10`,
				`int b = 12`,
				`fmt.Println("A: ", a, ", B: ", b)`,
				`int gozerTemp0 = b`,
				`int gozerTemp1 = a`,
				`a = gozerTemp0`,
				`b = gozerTemp1`,
				`fmt.Println("A: ", a, ", B: ", b)`),
		})
	})

	// Assignment005 checks assigment from the result of multiple return function.
	t.Run("Assignment005", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`n, err := fmt.Printf("Two = %d\n", 2)`,
				`fmt.Println("n: ", n, ", err: ", err)`),
			BlockSnippet: BlockMainMethodBody(
				`printResult gozerTemp0 = fmt.Printf("Two = %d\n", 2)`,
				`int n = gozerTemp0.n`,
				`error err = gozerTemp0.err`,
				`fmt.Println("n: ", n, ", err: ", err)`),
		})
	})

	// TODO: Assign to underscore

	// BinaryOp001 checks concatination of string literals.
	t.Run("BinaryOp001", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet:    GoMainMethodBody(`fmt.Print("Hello "+"World"+"!")`),
			BlockSnippet: BlockMainMethodBody(`fmt.Print((("Hello " + "World") + "!"))`),
		})
	})

	// BinaryOp002 checks multiple assignment and binary operation of integers.
	t.Run("BinaryOp002", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a, b := 1, 4`,
				`fmt.Print("Value: ", a + b)`),
			BlockSnippet: BlockMainMethodBody(
				`int a = 1`,
				`int b = 4`,
				`fmt.Print("Value: ", (a + b))`),
		})
	})

	// BinaryOp003 chesks multiplication of integers, subtration, and order of opertaions.
	t.Run("BinaryOp003", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a, b := 1, 4`,
				`fmt.Print("Value: ", a*0x10 - 1)`),
			BlockSnippet: BlockMainMethodBody(
				`int a = 1`,
				`int b = 4`,
				`fmt.Print("Value: ", ((a * 0x10) - 1))`),
		})
	})

	// BinaryOp004 checks some more binary operations and parentheses for ordering operations.
	t.Run("BinaryOp004", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a, b := 1, 4`,
				`fmt.Print("Value: ", (a*(0x10 - 1))>>b)`),
			BlockSnippet: BlockMainMethodBody(
				`int a = 1`,
				`int b = 4`,
				`fmt.Print("Value: ", ((a * (0x10 - 1)) >> b))`),
		})
	})

	// UnaryOp001 checks the negation operation and negative integer literals.
	t.Run("UnaryOp001", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a := 14`,
				`fmt.Print("Value 1: ", a)`,
				`a = -8`,
				`fmt.Print("Value 2: ", -a)`),
			BlockSnippet: BlockMainMethodBody(
				`int a = 14`,
				`fmt.Print("Value 1: ", a)`,
				`a = -8`,
				`fmt.Print("Value 2: ", -a)`),
		})
	})

	// UnaryOp002 checks negation and positive unary operations.
	t.Run("UnaryOp002", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a, b := 10, 12`,
				`fmt.Println("Value: ", +a - -b)`),
			BlockSnippet: BlockMainMethodBody(
				`int a = 10`,
				`int b = 12`,
				`fmt.Println("Value: ", (+a - -b))`),
		})
	})

	// TODO: Add tests for reference and dereference.

	// IfElse001 checks simple if-statement definitions and a comparitive binary operation.
	t.Run("IfElse001", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a := 10 `,
				`if a > 4 {`,
				`  fmt.Println("greater: ", a)`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`int a = 10`,
				`if (a > 4) {`,
				`  fmt.Println("greater: ", a)`,
				`}`),
		})
	})

	// IfElse002 checks an if-else-statement definition.
	t.Run("IfElse002", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a := 10 `,
				`if a > 4 {`,
				`  fmt.Println("greater: ", a)`,
				`} else {`,
				`  fmt.Println("less: ", a)`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`int a = 10`,
				`if (a > 4) {`,
				`  fmt.Println("greater: ", a)`,
				`} else {`,
				`  fmt.Println("less: ", a)`,
				`}`),
		})
	})

	// IfElse003 checks assignment within an if-else-statement and that scoping is working for if-else-statements.
	t.Run("IfElse003", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a := 2`,
				`if b, c := 10, 4; b > c + a {`,
				`  fmt.Println("greater: ", a, " > ", c + a)`,
				`} else {`,
				`  fmt.Println("less: ", a, " > ", c + a)`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`int a = 2`,
				`{`,
				`  int b = 10`,
				`  int c = 4`,
				`  if (b > (c + a)) {`,
				`    fmt.Println("greater: ", a, " > ", (c + a))`,
				`  } else {`,
				`    fmt.Println("less: ", a, " > ", (c + a))`,
				`  }`,
				`}`),
		})
	})

	// IfElse004 checks else-if-statements in the an if-else-statment.
	t.Run("IfElse004", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
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
			BlockSnippet: BlockMainMethodBody(
				`int a = 5`,
				`if (a > 10) {`,
				`  fmt.Println(a, " > 10")`,
				`} else if (a > 8) {`,
				`  fmt.Println(a, " > 8")`,
				`} else if (a > 0) {`,
				`  fmt.Println(a, " > 0")`,
				`} else {`,
				`  fmt.Println(a, " <= 0")`,
				`}`),
		})
	})

	// IfElse005 checks nested if-else-statments.
	t.Run("IfElse005", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
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
			BlockSnippet: BlockMainMethodBody(
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
				`}`),
		})
	})

	// TODO: Check assignment of varaibles in the else-if part of an if-else-statement.

	// For001 checks for-statement which increments.
	t.Run("For001", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`for i := 0; i < 10; i++ {`,
				`  fmt.Println("index = ", i)`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`for(int i = 0; (i < 10); i++) {`,
				`  fmt.Println("index = ", i)`,
				`}`),
		})
	})

	// For002 checks for-statement which decrements.
	t.Run("For002", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`for i := 9; i >= 0; i-- {`,
				`  fmt.Println("index = ", i)`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`for(int i = 9; (i >= 0); i--) {`,
				`  fmt.Println("index = ", i)`,
				`}`),
		})
	})

	// For003 checks for-statement with multi-assignment and multiple post operators.
	t.Run("For003", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`for i, j := 9, 0; j < 10; i, j = i-1, i+1 {`,
				`  fmt.Println("up = ", i, ", down = ", j)`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
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
				`}`),
		})
	})

	// For004 checks for-statement which only the check, a while-statement.
	t.Run("For004", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`i := 0`,
				`for i < 10 {`,
				`  fmt.Println("index = ", i)`,
				`  i++`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`int i = 0`,
				`for(; (i < 10); ) {`,
				`  fmt.Println("index = ", i)`,
				`  i++`,
				`}`),
		})
	})

	// For005 checks break and continue block statements.
	t.Run("For005", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
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
			BlockSnippet: BlockMainMethodBody(
				`for(int i = 0; true; i++) {`,
				`  if (i < 5) {`,
				`    fmt.Println("small: ", i)`,
				`    continue`,
				`  }`,
				`  fmt.Println("large: ", i)`,
				`  if (i > 10) {`,
				`    break`,
				`  }`,
				`}`),
		})
	})

	// For006 checks multiple-type of pre and post statements.
	t.Run("For006", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`for i, first := 0, true; i < 10; i, first = i+1, false {`,
				`  if first {`,
				`    fmt.Println("First: ", i)`,
				`    continue`,
				`  }`,
				`  fmt.Println("Rest: ", i)`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
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
				`}`),
		})
	})

	// ForListRange001 checks foreach on list with only index.
	t.Run("ForListRange001", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a := []int{1, 2, 3, 4}`,
				`for i := range a {`,
				`  fmt.Println("Index: ", i)`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`[]int a = []int{1, 2, 3, 4}`,
				`for(int i = 0; (i < len(a)); i++) {`,
				`  fmt.Println("Index: ", i)`,
				`}`),
		})
	})

	// ForListRange002 checks foreach on list with index and value.
	t.Run("ForListRange002", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a := []int{1, 2, 3, 4}`,
				`for i, val := range a {`,
				`  fmt.Println("Index: ", i, ", Value: ", val)`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`[]int a = []int{1, 2, 3, 4}`,
				`for(int i = 0; (i < len(a)); i++) {`,
				`  int val = a[i]`,
				`  fmt.Println("Index: ", i, ", Value: ", val)`,
				`}`),
		})
	})

	// ForListRange003 checks foreach with predefined value and index on list.
	t.Run("ForListRange003", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a := []int{1, 2, 3, 4}`,
				`val, i := -1, -2`,
				`for i, val = range a {`,
				`  fmt.Println("Index: ", i, ", Value: ", val)`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`[]int a = []int{1, 2, 3, 4}`,
				`int val = -1`,
				`int i = -2`,
				`for(i = 0; (i < len(a)); i++) {`,
				`  val = a[i]`,
				`  fmt.Println("Index: ", i, ", Value: ", val)`,
				`}`),
		})
	})

	// ForListRange004 checks foreach on a list with no iterators.
	t.Run("ForListRange004", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a := []int{1, 2, 3, 4}`,
				`for range a {`,
				`  fmt.Println("Bleep")`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`[]int a = []int{1, 2, 3, 4}`,
				`for(int gozerTemp0 = 0; (gozerTemp0 < len(a)); gozerTemp0++) {`,
				`  fmt.Println("Bleep")`,
				`}`),
		})
	})

	// ForListRange005 checks foreach with predefined value and no index on list.
	t.Run("ForListRange005", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`a := []int{1, 2, 3, 4}`,
				`for _, val := range a {`,
				`  fmt.Println("Value: ", val)`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`[]int a = []int{1, 2, 3, 4}`,
				`for(int gozerTemp0 = 0; (gozerTemp0 < len(a)); gozerTemp0++) {`,
				`  int val = a[gozerTemp0]`,
				`  fmt.Println("Value: ", val)`,
				`}`),
		})
	})

	// ForListRange006 checks foreach with predefined value and index on unidentified list.
	t.Run("ForListRange006", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`for i, val := range []int{1, 2, 3, 4} {`,
				`  fmt.Println("Index: ", i, ", Value: ", val)`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`{`,
				`  []int gozerTemp0 = []int{1, 2, 3, 4}`,
				`  for(int i = 0; (i < len(gozerTemp0)); i++) {`,
				`    int val = gozerTemp0[i]`,
				`    fmt.Println("Index: ", i, ", Value: ", val)`,
				`  }`,
				`}`),
		})
	})

	// ForListRange007 checks foreach with predefined value and no index on unidentified list.
	t.Run("ForListRange007", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`for _, val := range []int{1, 2, 3, 4} {`,
				`  fmt.Println("Value: ", val)`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`{`,
				`  []int gozerTemp0 = []int{1, 2, 3, 4}`,
				`  for(int gozerTemp1 = 0; (gozerTemp1 < len(gozerTemp0)); gozerTemp1++) {`,
				`    int val = gozerTemp0[gozerTemp1]`,
				`    fmt.Println("Value: ", val)`,
				`  }`,
				`}`),
		})
	})

	// ForListRange008 checks foreach with predefined no value and no index on unidentified list.
	t.Run("ForListRange008", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`for _, _ := range []int{1, 2, 3, 4} {`,
				`  fmt.Println("Bleep")`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`{`,
				`  []int gozerTemp0 = []int{1, 2, 3, 4}`,
				`  for(int gozerTemp1 = 0; (gozerTemp1 < len(gozerTemp0)); gozerTemp1++) {`,
				`    fmt.Println("Bleep")`,
				`  }`,
				`}`),
		})
	})

	// TODO: For-each map

	// Slices001 checks the creation of a slicee and the build-in length method.
	t.Run("Slices001", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`arr := []int{4, 1, 3, 2}`,
				`fmt.Println("Count = ", len(arr))`),
			BlockSnippet: BlockMainMethodBody(
				`[]int arr = []int{4, 1, 3, 2}`,
				`fmt.Println("Count = ", len(arr))`),
		})
	})

	// Slices002 checks the indices of a slice expression.
	t.Run("Slices002", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`arr := []int{4, 1, 3, 2}`,
				`for i := 0; i < len(arr); i++ {`,
				`  fmt.Println(i, ": ", arr[i])`,
				`}`),
			BlockSnippet: BlockMainMethodBody(
				`[]int arr = []int{4, 1, 3, 2}`,
				`for(int i = 0; (i < len(arr)); i++) {`,
				`  fmt.Println(i, ": ", arr[i])`,
				`}`),
		})
	})

	// Slices003 checks assignment of a slice element with an index.
	t.Run("Slices003", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`arr := []int{4, 1, 3, 2}`,
				`arr[2] = 8`,
				`fmt.Println("arr[2] = ", arr[2])`),
			BlockSnippet: BlockMainMethodBody(
				`[]int arr = []int{4, 1, 3, 2}`,
				`arr[2] = 8`,
				`fmt.Println("arr[2] = ", arr[2])`),
		})
	})

	// Slices004 checks appending to a slice.
	t.Run("Slices004", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`arr := []int{4, 1, 3, 2}`,
				`arr = append(arr, 8)`,
				`fmt.Println("Count = ", len(arr))`),
			BlockSnippet: BlockMainMethodBody(
				`[]int arr = []int{4, 1, 3, 2}`,
				`arr = append(arr, 8)`,
				`fmt.Println("Count = ", len(arr))`),
		})
	})

	// Slices005 checks the creation of a slice with a default length
	// and a check of the build-in capacity method.
	t.Run("Slices005", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`arr := make([]int, 4)`,
				`arr[2] = 8`,
				`fmt.Println("arr[2] = ", arr[2])`,
				`fmt.Println("Count = ", len(arr))`,
				`fmt.Println("Cap = ", cap(arr))`),
			BlockSnippet: BlockMainMethodBody(
				`[]int arr = make([]int, 4)`,
				`arr[2] = 8`,
				`fmt.Println("arr[2] = ", arr[2])`,
				`fmt.Println("Count = ", len(arr))`,
				`fmt.Println("Cap = ", cap(arr))`),
		})
	})

	// Slices006 checks the creation of a slice with a default capacity.
	t.Run("Slices006", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`arr := make([]int, 0, 4)`,
				`arr = append(arr, 8)`,
				`fmt.Println("Count = ", len(arr))`,
				`fmt.Println("Cap = ", cap(arr))`),
			BlockSnippet: BlockMainMethodBody(
				`[]int arr = make([]int, 0, 4)`,
				`arr = append(arr, 8)`,
				`fmt.Println("Count = ", len(arr))`,
				`fmt.Println("Cap = ", cap(arr))`),
		})
	})

	// Slices007 checks creating a subslice and assigning to that subslice.
	t.Run("Slices007", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`arr := []int{4, 1, 3, 2}`,
				`arr2 := arr[1:2]`,
				`arr2[0] = 8`,
				`arr[2] = 7`,
				`fmt.Printf("arr = %v", arr)`,
				`fmt.Printf("arr2 = %v", arr2)`),
			BlockSnippet: BlockMainMethodBody(
				`[]int arr = []int{4, 1, 3, 2}`,
				`[]int arr2 = arr[1:2]`,
				`arr2[0] = 8`,
				`arr[2] = 7`,
				`fmt.Printf("arr = %v", arr)`,
				`fmt.Printf("arr2 = %v", arr2)`),
		})
	})

	// Slices008 checks creating a subslice from the beginning to an index.
	t.Run("Slices008", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`arr := []int{4, 1, 3, 2}`,
				`arr2 := arr[:2]`,
				`arr2[0] = 8`,
				`arr[2] = 7`,
				`fmt.Printf("arr = %v", arr)`,
				`fmt.Printf("arr2 = %v", arr2)`),
			BlockSnippet: BlockMainMethodBody(
				`[]int arr = []int{4, 1, 3, 2}`,
				`[]int arr2 = arr[:2]`,
				`arr2[0] = 8`,
				`arr[2] = 7`,
				`fmt.Printf("arr = %v", arr)`,
				`fmt.Printf("arr2 = %v", arr2)`),
		})
	})

	// Slices009 checks creating a subslice from an index to the end.
	t.Run("Slices009", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`arr := []int{4, 1, 3, 2}`,
				`arr2 := arr[1:]`,
				`arr2[0] = 8`,
				`arr[2] = 7`,
				`fmt.Printf("arr = %v", arr)`,
				`fmt.Printf("arr2 = %v", arr2)`),
			BlockSnippet: BlockMainMethodBody(
				`[]int arr = []int{4, 1, 3, 2}`,
				`[]int arr2 = arr[1:]`,
				`arr2[0] = 8`,
				`arr[2] = 7`,
				`fmt.Printf("arr = %v", arr)`,
				`fmt.Printf("arr2 = %v", arr2)`),
		})
	})

	// Slices010 checks creating a subslice from an index to the end.
	t.Run("Slices010", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`arr := []int{4, 1, 3, 2}`,
				`arr2 := arr[1:2:3]`,
				`arr2[0] = 8`,
				`arr[2] = 7`,
				`fmt.Printf("arr = %v", arr)`,
				`fmt.Printf("arr2 = %v", arr2)`),
			BlockSnippet: BlockMainMethodBody(
				`[]int arr = []int{4, 1, 3, 2}`,
				`[]int arr2 = arr[1:2:3]`,
				`arr2[0] = 8`,
				`arr[2] = 7`,
				`fmt.Printf("arr = %v", arr)`,
				`fmt.Printf("arr2 = %v", arr2)`),
		})
	})

	// Maps001 checks initialization of a map.
	t.Run("Maps001", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`numbers := map[int]string{4: "four", 1: "one", 3: "three", 2: "two"}`,
				`fmt.Println("Count = ", len(numbers))`),
			BlockSnippet: BlockMainMethodBody(
				`map[int]string numbers = map[int]string{4: "four", 1: "one", 3: "three", 2: "two"}`,
				`fmt.Println("Count = ", len(numbers))`),
		})
	})

	// Maps002 checks assignment of a map element with an int key.
	t.Run("Maps002", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`numbers := map[int]string{4: "four", 1: "one", 3: "three", 2: "two"}`,
				`numbers[2] = "duo"`,
				`fmt.Println("numbers[2] = ", numbers[2])`),
			BlockSnippet: BlockMainMethodBody(
				`map[int]string numbers = map[int]string{4: "four", 1: "one", 3: "three", 2: "two"}`,
				`numbers[2] = "duo"`,
				`fmt.Println("numbers[2] = ", numbers[2])`),
		})
	})

	// Maps003 checks assignment of a map element with an string key.
	t.Run("Maps003", func(t *testing.T) {
		hndl(t, map[SnippetType]string{
			GoSnippet: GoMainMethodBody(
				`numbers := map[string]int{"four": 4, "one": 1, "three": 3, "two": 2}`,
				`numbers["two"] = 42`,
				`fmt.Println("numbers[two] = ", numbers["two"])`),
			BlockSnippet: BlockMainMethodBody(
				`map[string]int numbers = map[string]int{"four": 4, "one": 1, "three": 3, "two": 2}`,
				`numbers["two"] = 42`,
				`fmt.Println("numbers[two] = ", numbers["two"])`),
		})
	})
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
