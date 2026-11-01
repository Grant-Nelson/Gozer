package typeTools

import (
	"go/token"
	"go/types"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/Grant-Nelson/Gozer/avail/typeTools/typeOp"
)

func Test_Ops_Arrays(t *testing.T) {
	checkOpsAllowance(t, lines(
		`func Foo(x, y [3]int) {`,
		`	_, _ = len(x), cap(x)`,               // Len | Cap
		`	_ = x == y`,                          // Comparable
		`	_ = &x`,                              // Ref
		`	x[1] = x[0]`,                         // GetIndex, SetIndex
		`	p := &x[0]; *p = 4`,                  // RefIndex
		`	_, _ = x[0:1], x[0:1:2]`,             // Slice | Slice3
		`	for range x { print("hi") }`,         // Range0
		`	for i := range x { print(i) }`,       // Range1
		`	for i, v := range x { print(i, v) }`, // Range2
		`}`), `Ok`)

	checkOps(t, `[3]int`,
		`Len|Cap|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)

	checkOps(t, `[...]byte{1, 2, 3}`,
		`Len|Cap|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:byte, Slice:[]byte, Range1:int, Range2:byte }`)
}

func Test_Ops_PointersToArrays(t *testing.T) {
	checkOpsAllowance(t, lines(
		`func Foo(x, y *[3]int) {`,
		`	_ = x == nil`,                        // IsNil
		`	_, _ = len(x), cap(x)`,               // Len | Cap
		`	_ = x == y`,                          // Comparable
		`	_, _ = &x, *x`,                       // Ref
		`	x[1] = x[0]`,                         // GetIndex, SetIndex
		`	p := &x[0]; *p = 4`,                  // RefIndex
		`	_, _ = x[0:1], x[0:1:2]`,             // Slice | Slice3
		`	for range x { print("hi") }`,         // Range0
		`	for i := range x { print(i) }`,       // Range1
		`	for i, v := range x { print(i, v) }`, // Range2
		`}`), `Ok`)

	checkOps(t, `&[3]int{1,2,3}`,
		`Len|Cap|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)

	checkOps(t, `*[3]int`,
		`Len|Cap|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)

	checkOps(t, `*[3]byte`,
		`Len|Cap|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:byte, Slice:[]byte, Range1:int, Range2:byte }`)

	checkOps(t, `(*int)(nil)`,
		`IsNil|Comparable|Ref`)

	checkOps(t, `**[3]byte`,
		`IsNil|Comparable|Ref`)

	checkOpsWithFile(t, lines(
		`package t`,
		`type Foo [3]int`,
	), token.NoPos, `new(Foo)`,
		`Len|Cap|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)
}

func Test_Ops_Slices(t *testing.T) {
	checkOpsAllowance(t, lines(
		`func Foo(x, y []int) {`,
		`	_ = x == nil`,                        // IsNil
		`	_, _ = len(x), cap(x)`,               // Len | Cap
		`	_ = &x`,                              // Ref
		`	x[1] = x[0]`,                         // GetIndex, SetIndex
		`	p := &x[0]; *p = 4`,                  // RefIndex
		`	_, _ = x[0:1], x[0:1:2]`,             // Slice | Slice3
		`	for range x { print("hi") }`,         // Range0
		`	for i := range x { print(i) }`,       // Range1
		`	for i, v := range x { print(i, v) }`, // Range2
		`}`), `Ok`)

	checkOps(t, `[]int`,
		`Len|Cap|IsNil|Ref|Make|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)

	checkOps(t, `[]*float64`,
		`Len|Cap|IsNil|Ref|Make|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:*float64, Slice:[]*float64, Range1:int, Range2:*float64 }`)

	checkOps(t, `[]byte`,
		`Len|Cap|IsNil|Ref|Make|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:byte, Slice:[]byte, Range1:int, Range2:byte }`)

	checkOps(t, `*[]byte`,
		`IsNil|Comparable|Ref`)

	checkOpsWithFile(t, lines(
		`package t`,
		`type Bar struct{ v int }`,
		`type Foo []*Bar`,
	), token.NoPos, `Foo`,
		`Len|Cap|IsNil|Ref|Make|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:*$.Bar, Slice:$.Foo, Range1:int, Range2:*$.Bar }`)
}

func Test_Ops_BooleanTypes(t *testing.T) {
	checkOps(t, `bool`, `Comparable|Ref`)
	checkOps(t, `false`, `Comparable|Ref`)

	checkOpsWithFile(t, lines(
		`package t`,
		`type Foo bool`,
	), token.NoPos, `Foo`,
		`Comparable|Ref`)
}

func Test_Ops_IntegerTypes(t *testing.T) {
	checkOps(t, `int`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1`+
			`{ Range1:int }`)

	checkOps(t, `uint`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1`+
			`{ Range1:uint }`)

	checkOps(t, `int64`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1`+
			`{ Range1:int64 }`)

	checkOps(t, `42`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1`+
			`{ Range1:untyped int }`)
}

func Test_Ops_FloatTypes(t *testing.T) {
	checkOps(t, `float32`,
		`Add|Arith|Comparable|Orderable|Ref|Float`+
			`{ Complex:complex64 }`)

	checkOps(t, `float64`,
		`Add|Arith|Comparable|Orderable|Ref|Float`+
			`{ Complex:complex128 }`)

	checkOps(t, `1.42`,
		`Add|Arith|Comparable|Orderable|Ref|Float`+
			`{ Complex:untyped complex }`)
}

func Test_Ops_ComplexTypes(t *testing.T) {
	checkOps(t, `complex128`,
		`Add|Arith|Comparable|Ref|Complex`+
			`{ Float:float64 }`)

	checkOps(t, `complex64`,
		`Add|Arith|Comparable|Ref|Complex`+
			`{ Float:float32 }`)

	checkOps(t, `1.42+6.2i`,
		`Add|Arith|Comparable|Ref|Complex`+
			`{ Float:untyped float }`)

	checkOps(t, `interface{float64|complex128}`,
		`Add|Arith|Comparable|Ref`)
}

func Test_Ops_StringTypes(t *testing.T) {
	checkOps(t, `string`,
		`Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:uint8, Slice:string, Range1:int, Range2:rune }`)

	checkOps(t, `"Hello"`,
		`Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:uint8, Slice:untyped string, Range1:int, Range2:rune }`)

	checkOpsWithFile(t, lines(
		`package t`,
		`type Foo string`,
	), token.NoPos, `Foo`,
		`Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:uint8, Slice:$.Foo, Range1:int, Range2:rune }`)
}

func Test_Ops_RuneTypes(t *testing.T) {
	checkOps(t, `rune`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1`+
			`{ Range1:rune }`)

	checkOps(t, `'☺'`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1`+
			`{ Range1:untyped rune }`)
}

func Test_Ops_PointersTypes(t *testing.T) {
	checkOpsWithFile(t, lines(
		`package t`,
		`import "unsafe"`,
		`var x = 10`,
		`var p = unsafe.Pointer(&x)`,
	), token.NoPos, `p`,
		`IsNil|Comparable|Ref`)

	checkOps(t, `nil`, `IsNil|Ref`)
}

func Test_Ops_ChanTypes(t *testing.T) {
	checkOps(t, `chan int`,
		`Len|Cap|IsNil|Comparable|Ref|Make|Range0|Range1|Recv|Send`+
			`{ Range1:int }`)

	checkOps(t, `make(chan int, 3)`,
		`Len|Cap|IsNil|Comparable|Ref|Make|Range0|Range1|Recv|Send`+
			`{ Range1:int }`)

	checkOps(t, `new(chan int)`,
		`IsNil|Comparable|Ref`)

	checkOps(t, `chan<- int`,
		`Len|Cap|IsNil|Comparable|Ref|Make|Send`)

	checkOps(t, `<-chan int`,
		`Len|Cap|IsNil|Comparable|Ref|Make|Range0|Range1|Recv`+
			`{ Range1:int }`)
}

func Test_Ops_UnionsTypes_Basics(t *testing.T) {
	checkOps(t, `interface{String() string}`, `IsNil|Ref`)

	checkOps(t, `interface{int|uint}`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref`)

	checkOpsWithFile(t, lines(
		`package t`,
		`func Foo[T int|uint](t T) {`,
		`	println(t)`,
		`}`,
	), token.Pos(40), `t`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref`)

	checkOps(t, `interface{int64|float64}`,
		`Add|Arith|Comparable|Orderable|Ref`)

	checkOps(t, `interface{String() string; int|uint}`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref`)
}

func Test_Ops_UnionsTypes_WithSharedOps(t *testing.T) {
	checkOps(t, `interface{~[]byte|~string}`,
		`Len|Ref|GetIndex`+
			`{ Key:untyped int, Elem:byte }`)

	checkOps(t, `interface{~[]byte|string}`,
		`Len|Ref|GetIndex`+
			`{ Key:untyped int, Elem:byte }`)

	checkOps(t, `interface{[]byte|~string}`,
		`Len|Ref|GetIndex`+
			`{ Key:untyped int, Elem:byte }`)

	checkOps(t, `interface{[]byte|string}`,
		`Len|Ref|GetIndex`+
			`{ Key:untyped int, Elem:byte }`)
}

func Test_Ops_UnionsTypes_UnionReduction(t *testing.T) {
	// From https://cs.opensource.google/go/x/exp/+/e88cd736:typeparams/normalize.go;l=66
	// > In this example, the structural type restriction of P is ~string|int: A|B
	// > expands to ~string|~[]byte|int|string, which reduces to ~string|~[]byte|int,
	// > which when intersected with C (~string|~int) yields ~string|int.
	checkOpsWithFile(t, lines(
		`package t`,
		`type A interface{ ~string|~[]byte }`,
		`type B interface{ int|string }`,
		`type C interface { ~string|~int }`,
		`type P interface{ A|B; C }`, // ~string|int
	), token.Pos(45), `P`,
		`Add|Comparable|Orderable|Ref`+
			`{ Range1:int }`)
}

func Test_Ops_UnionsTypes_SlicesAndArrays(t *testing.T) {
	checkOps(t, `interface{~[][]byte|~[]string}`,
		`Len|Cap|IsNil|Ref|Make|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:{[]byte | string}, Slice:{~[][]byte | ~[]string}, Range1:int, Range2:{[]byte | string} }`)

	checkOps(t, `interface{~[]int|~[3]int}`,
		`Len|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:~[]int, Range1:int, Range2:int }`)

	checkOps(t, `interface{~[]int|~*[3]int}`,
		`Len|IsNil|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:~[]int, Range1:int, Range2:int }`)

	checkOps(t, `interface{~[]byte|~[]int}`,
		`Len|Cap|IsNil|Ref|Make|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:{byte | int}, Slice:{~[]byte | ~[]int}, Range1:int, Range2:{byte | int} }`)

	checkOps(t, `interface{~[]int|~[3]*int}`,
		`Len|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:{int | *int}, Slice:{~[]int | []*int}, Range1:int, Range2:{int | *int} }`)
}

func Test_Ops_UnionsTypes_CompoundGenerics(t *testing.T) {
	checkOpsWithFile(t, lines(
		`package t`,
		`func Foo[T any, S ~[]T](s S) {`,
		`	println(s)`,
		`}`,
	), token.Pos(45), `s`,
		`Len|Cap|IsNil|Ref|Make|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:T, Slice:S, Range1:int, Range2:T }`)

	checkOpsWithFile(t, lines(
		`package t`,
		`func Foo[T any, S ~[]T | ~map[int]T](s S) {`,
		`	println(s)`,
		`}`,
	), token.Pos(65), `s`,
		`Len|IsNil|Ref|Make|GetIndex|SetIndex|Range0|Range1|Range2`+
			`{ Key:int, Elem:T, Range1:int, Range2:T }`)
}

func Test_Ops_FunctionsTypes(t *testing.T) {
	// TODO: Create a test to check c.Walk has the correct value
	//
	//	package main
	//
	//	import "fmt"
	//
	//	type Cat struct{ lives int }
	//
	//	func (c *Cat) Walk(y func(int) bool) {
	//		for i := range c.lives {
	//			if !y(i) {
	//				return
	//			}
	//		}
	//	}
	//
	//	func main() {
	//		c := &Cat{lives: 3}
	//		for i := range c.Walk {
	//			fmt.Println(i)
	//		}
	//	}
}

// TODO: Check comparable

func checkOpsAllowance(t *testing.T, code, exp string) {
	t.Helper()

	prime := "package main\n" + code + "\n"
	got := `Ok`
	_, err := astTools.ParsePackage(map[string]string{`t.go`: prime}, nil)
	if err != nil {
		got = removePaths(err.Error())
	}

	if diff := cmp.Diff(exp, got); len(diff) > 0 {
		diff = strings.TrimSpace(diff)
		diff = strings.ReplaceAll(diff, "\n", "\n    ")
		t.Errorf("Unexpected results:\ncode: %s\ndiff: %s\ngot: %q\nexp: %q",
			code, diff, got, exp)
	}
}

func checkOpsAllowed(t *testing.T, typeParams, body, exp string) {
	t.Helper()
	checkOpsAllowance(t,
		"func Foo["+typeParams+"](x X, i int) {\n"+
			"	"+body+"\n"+
			"}", exp)
}

func Test_Ops_Go_Allowance_SliceAndMap(t *testing.T) {
	checkOpsAllowed(t,
		`T any, X ~[]T | ~map[int]T`,
		`_ = x == nil`,
		`Ok`)
	checkOpsAllowed(t,
		`T any, X ~[]T | ~map[int]T`,
		`var t T = x[0]; _ = t`,
		`cannot index x (variable of type X constrained by ~[]T | ~map[int]T)`)
	checkOpsAllowed(t,
		`T any, X ~[]T | ~map[int]T`,
		`clear(x)`,
		`Ok`)
}

func Test_Ops_Go_Allowance_SliceAndArray(t *testing.T) {
	checkOpsAllowed(t,
		`X ~[4]int`,
		`x[1] = x[0]; _, _ = len(x), cap(x)`,
		`Ok`)
	checkOpsAllowed(t,
		`T any, X ~[]T | ~[4]T`,
		`var t T = x[0]; x[1] = t; _, _ = len(x), cap(x)`,
		`Ok`)
	checkOpsAllowed(t,
		`T any, X ~[]T | ~[4]T`,
		`for i, v := range x { _, _ = i, v }`,
		`cannot range over x (variable of type X constrained by ~[]T | ~[4]T): []T and [4]T have different underlying types`)
	checkOpsAllowed(t,
		`T any, X ~[]T | ~*[4]T`,
		`var t T = x[0]; _ = &x[0]`,
		`declared and not used: t`)
	checkOpsAllowed(t,
		`T any, X ~[]T | ~*[4]T`,
		`var y []T = x[1:2]; _ = y`,
		`cannot slice x (variable of type X constrained by ~[]T | ~*[4]T): []T and *[4]T have different underlying types`)
	checkOpsAllowed(t,
		`T any, X ~[]T | ~[4]T`,
		`var y []T = x[1:2]; _ = y`,
		`cannot slice x (variable of type X constrained by ~[]T | ~[4]T): []T and [4]T have different underlying types`)
	checkOpsAllowed(t,
		`T any, X ~[]T | ~[3]T`,
		`_ = cap(x)`,
		`Ok`)
	checkOpsAllowed(t,
		`T any, X ~[]T | ~[3]T | ~*[3]T`,
		`var t T = x[i]; x[i] = t; _ = &x[i]`,
		`Ok`)
}

func Test_Ops_Go_Allowance_SliceAndString(t *testing.T) {
	checkOpsAllowed(t,
		`T any, X ~[]byte | ~string`,
		`var b byte = x[0]; _ = x[1:2]`,
		`declared and not used: b`)
	checkOpsAllowed(t,
		`T any, X ~[]byte | ~string`,
		`for j, v := range x { _, _ = j, v }`,
		`cannot range over x (variable of type X constrained by ~[]byte | ~string): []byte and string have different underlying types`)
	checkOpsAllowed(t,
		`X ~[]byte | ~string`,
		`var y X = x[1:2]; _ = y`,
		`Ok`)
	checkOpsAllowed(t,
		`X ~string | ~[]byte`,
		`var y X = x + x; _ = y`,
		`invalid operation: operator + not defined on x (variable of type X constrained by ~string | ~[]byte)`)
	checkOpsAllowed(t,
		`X ~[]byte | ~string | ~[4]byte`,
		`var b byte = x[i]; _ = b`,
		`Ok`)
	checkOpsAllowed(t,
		`X ~[]byte | ~string`,
		`x[i] = 1`,
		`cannot assign to x[i] (neither addressable nor a map index expression)`)
	checkOpsAllowed(t,
		`X ~[]byte | ~string`,
		`_ = &x[i]`,
		`invalid operation: cannot take address of x[i] (value of type byte)`)
	checkOpsAllowed(t,
		`X ~[]byte | ~string`,
		`var y X = x[1:2:3]; _ = y`,
		`invalid operation: 3-index slice of string`)
	checkOpsAllowed(t,
		`X ~[]byte | ~string`,
		`_ = len(x); _ = cap(x)`,
		`invalid argument: x (variable of type X constrained by ~[]byte | ~string) for built-in cap`)
	checkOpsAllowed(t,
		`X ~string | ~[]byte`,
		`var y X = x[i:]; _ = y`,
		`Ok`)
}

func Test_Ops_Go_Allowance_SliceAndSlice(t *testing.T) {
	checkOpsAllowed(t,
		`X ~[]int | ~[]int8`,
		`_ = x[i]`,
		`cannot index x (variable of type X constrained by ~[]int | ~[]int8)`)
	checkOpsAllowed(t,
		`X ~[]int | ~[]int8`,
		`_, _ = len(x), cap(x)`,
		`Ok`)
	checkOpsAllowed(t,
		`X ~[]int | ~[]int8`,
		`_ = x == nil`,
		`Ok`)
}

func Test_Ops_Go_Allowance_MapAndMap(t *testing.T) {
	checkOpsAllowed(t,
		`T any, X ~map[int]T | ~map[int]T`,
		`for k, v := range x { _, _ = k, v }`,
		`overlapping terms ~map[int]T and ~map[int]T`)
	checkOpsAllowed(t,
		`T any, X ~map[int8]T | ~map[int16]T`,
		`for k, v := range x { _, _ = k, v }`,
		`cannot range over x (variable of type X constrained by ~map[int8]T | ~map[int16]T): `+
			`map[int8]T and map[int16]T have different underlying types`)
	checkOpsAllowed(t,
		`X ~map[int][]byte | ~map[int]string`,
		`_ = x[i]`,
		`cannot index x (variable of type X constrained by ~map[int][]byte | ~map[int]string)`)
	checkOpsAllowed(t,
		`T ~string, X ~map[T]int`,
		`for k, v := range x { _, _ = k, v }`,
		`Ok`)
	checkOpsAllowed(t,
		`T ~[]byte, X ~map[T]int`,
		`for k, v := range x { _, _ = k, v }`,
		`invalid map key type T (missing comparable constraint)`)
	checkOpsAllowed(t,
		`T ~[]byte | ~string, X ~map[T]int`,
		`for k, v := range x { _, _ = k, v }`,
		`invalid map key type T (missing comparable constraint)`)
	checkOpsAllowed(t,
		`T interface{~[]byte | ~string; comparable}, X ~map[T]int`,
		`for k, v := range x { _, _ = k, v }`,
		// This seems like it would require `~string` and `~[]byte` with
		// `comparable` would be impossible. Seems like this shouldn't be allowed.
		`Ok`)
	checkOpsAllowed(t,
		`T interface{~[]byte; comparable}, X ~map[T]int`,
		`for k, v := range x { _, _ = k, v }`,
		`invalid map key type T (missing comparable constraint)`)
}

func Test_Ops_Go_Allowance_BasicAndBasic(t *testing.T) {
	checkOpsAllowed(t,
		`T any, X ~int | ~float64`,
		`_ = x + 1`,
		`Ok`)
	checkOpsAllowed(t,
		`T any, X ~int | ~float64`,
		`_ = x + 1.5`,
		`cannot convert 1.5 (untyped float constant) to type X`)
	checkOpsAllowed(t,
		`T any, X ~uint8 | ~int16`,
		`_ = x + 200`,
		`Ok`)
	checkOpsAllowed(t,
		`T any, X ~uint8 | ~int16`,
		`_ = x + 300`,
		`cannot convert 300 (untyped int constant) to type X`)
	checkOpsAllowed(t,
		`T any, X ~int | ~float64`,
		`_ = x % 2`,
		`invalid operation: operator % not defined on x (variable of type X constrained by ~int | ~float64)`)
	checkOpsAllowed(t,
		`X ~int | ~int8`,
		`for i := range x { var y X = i; _ = y }`,
		`cannot range over x (variable of type X constrained by ~int | ~int8): int and int8 have different underlying types`)
	checkOpsAllowed(t,
		`X ~int`,
		`for i := range x { var y X = i; _ = y }`,
		`Ok`)
	checkOpsAllowed(t,
		`X ~int | ~int8`,
		`for i := range 10 { var y X = i; _ = y }`,
		`cannot use i (variable of type int) as X value in variable declaration`)
}

func Test_Ops_Go_Allowance_Complex(t *testing.T) {
	checkOpsAllowed(t,
		`X ~complex64 | ~complex128`,
		`_ = real(x)`,
		`x (variable of type X constrained by ~complex64 | ~complex128) not supported as argument to built-in real for go1.18 (see go.dev/issue/50937)`)
	checkOpsAllowed(t,
		`X ~complex64`,
		`var f float32 = real(x); _ = f`,
		"x (variable of type X constrained by ~complex64) not supported as argument to built-in real for go1.18 (see go.dev/issue/50937)\n"+
			"cannot use real(x) (value of type X constrained by interface{~float32}) as float32 value in variable declaration")
	checkOpsAllowed(t,
		`X ~float32 | ~float64`,
		`_ = complex(x, x)`,
		`x (variable of type X constrained by ~float32 | ~float64) not supported as argument to built-in complex for go1.18 (see go.dev/issue/50937)`)
	checkOpsAllowed(t,
		`X ~float32`,
		`var c complex64 = complex(x, x); _ = c`,
		"x (variable of type X constrained by ~float32) not supported as argument to built-in complex for go1.18 (see go.dev/issue/50937)\n"+
			"cannot use complex(x, x) (value of type X constrained by interface{~complex64}) as complex64 value in variable declaration")
}

func lines(lines ...string) string { return strings.Join(lines, "\n") }

var regexPathMatch = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^.*t.go:(?:\d+:){0,2}\s`)
})

func removePaths(s string) string {
	buf := &strings.Builder{}
	for ln := range strings.Lines(s) {
		ln = regexPathMatch().ReplaceAllString(ln, ``)
		_, _ = buf.WriteString(ln)
	}
	return buf.String()
}

func getOpsWithFile(t *testing.T, prime string, pos token.Pos, expr string) OpTypes {
	t.Helper()

	ps, err := astTools.ParsePackage(map[string]string{`t.go`: prime}, nil)
	if err != nil {
		t.Fatalf(`Failed to parse package: %v`, err)
	}

	if !pos.IsValid() {
		pos = ps.Syntax[0].FileEnd - 1
	}

	tv, err := types.Eval(ps.Fset, ps.Types, pos, expr)
	if err != nil {
		t.Fatalf(`Failed to parse expression: %v`, err)
	}

	return Ops(tv.Type)
}

func checkOps(t *testing.T, expr, exp string) {
	t.Helper()
	checkOpsWithFile(t, "package t\n", token.NoPos, expr, exp)
}

func checkOpsWithFile(t *testing.T, prime string, pos token.Pos, expr, exp string) {
	t.Helper()
	opt := getOpsWithFile(t, prime, pos, expr)
	got := opt.String()
	got = strings.ReplaceAll(got, `command-line-arguments`, `$`)

	if diff := cmp.Diff(exp, got); len(diff) > 0 {
		diff = strings.TrimSpace(diff)
		diff = strings.ReplaceAll(diff, "\n", "\n    ")
		t.Errorf("Unexpected results:\ntype: %s\ndiff: %s\ngot: %q\nexp: %q",
			opt.Type.String(), diff, got, exp)
	}
}

func checkGoOps(t *testing.T, prime, typeParams, typeExp string, expOp typeOp.Op) {
	t.Helper()

	opGroups := []struct {
		op   typeOp.Op
		code string
	}{
		{
			op: typeOp.Add,
			code: lines(
				`x = y + z`,
				`x += y`),
		}, {
			op: typeOp.Arith,
			code: lines(
				`x = -y`,
				`x = y - z`,
				`x = y * z`,
				`x = y / z`,
				`x -= y`,
				`x *= y`,
				`x /= y`,
				`x++`,
				`x--`,
				`_ = int(y)`,
				`_ = float64(y)`),
		}, {
			op: typeOp.Mod,
			code: lines(
				`x = y % z`,
				`x %= y`),
		}, {
			op: typeOp.Bitwise,
			code: lines(
				`x = ^y`,
				`x = y & z`,
				`x = y | z`,
				`x = y ^ z`,
				`x = y &^ z`,
				`x = y << z`,
				`x = y >> z`,
				`x &= y`,
				`x |= y`,
				`x ^= y`,
				`x &^= y`,
				`x <<= y`,
				`x >>= y`),
		}, {
			op:   typeOp.Len,
			code: `_ = len(y)`,
		}, {
			op:   typeOp.Cap,
			code: `_ = cap(y)`,
		}, {
			op:   typeOp.IsNil,
			code: `_ = y == nil`,
		}, {
			op: typeOp.Comparable,
			code: lines(
				`_ = y == z`,
				`_ = y != z`),
		}, {
			op: typeOp.Orderable,
			code: lines(
				`_ = y < z`,
				`_ = y <= z`,
				`_ = y > z`,
				`_ = y >= z`),
		}, {
			op: typeOp.Ref,
			code: lines(
				`_ = &y`,
				`_ = new(y)`,
				`_ = new(`+typeExp+`)`),
		}, {
			op: typeOp.Make,
			code: lines(
				`_ = make(`+typeExp+`)`,
				`_ = make(`+typeExp+`, 1)`,
				`_ = make(`+typeExp+`, 1, 2)`),
		}, {
			op: typeOp.GetIndex,
			code: lines(
				`_ = y[0]`,
				`_, _ = y[0]`),
		}, {
			op: typeOp.SetIndex,
			code: lines(
				`x[0] = y[0]`,
				`clear(x)`),
		}, {
			op:   typeOp.RefIndex,
			code: `_ = &y[0]`,
		}, {
			op:   typeOp.Slice,
			code: `x = y[0:1]`,
		}, {
			op:   typeOp.Slice3,
			code: `x = y[0:1:2]`,
		}, {
			op: typeOp.Range0,
			code: lines(
				`for range x {`,
				`	print("hi")`,
				`}`),
		}, {
			op: typeOp.Range1,
			code: lines(
				`for i := range x {`,
				`	print(i)`,
				`}`),
		}, {
			op: typeOp.Range2,
			code: lines(
				`for i, j := range x {`,
				`	print(i, j)`,
				`}`),
		}, {
			op: typeOp.Recv,
			code: lines(
				`_ = <- x`,
				`_, _ = <- x`),
		}, {
			op:   typeOp.Send,
			code: `close(x)`,
		}, {
			op:   typeOp.Float,
			code: `_ = complex(x, y)`,
		}, {
			op: typeOp.Complex,
			code: lines(
				`_ = real(x)`,
				`_ = imag(x)`),
		},
	}

	if len(typeParams) > 0 {
		typeParams = `[` + typeParams + `]`
	}

	for _, tt := range opGroups {
		prime := lines(
			`package t`,
			prime,
			`func Foo`+typeParams+`(x, y, z `+typeExp+`) {`,
			tt.code,
			`}`)
		_, err := astTools.ParsePackage(map[string]string{`t.go`: prime}, nil)
		if err == nil {
			if !expOp.Any(tt.op) {
				t.Errorf(`%v was allowed by Go but not in the expected ops`, tt.op.String())
			}
		} else if expOp.Any(tt.op) {
			t.Errorf(`%v was in the expected ops but Go errored with %s`, tt.op.String(), err.Error())
		}
	}
}
