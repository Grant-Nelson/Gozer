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
)

func lines(lines ...string) string { return strings.Join(lines, "\n") }

func checkOps(t *testing.T, expr, exp string) {
	t.Helper()
	checkOpsWithFile(t, "package t\n", token.NoPos, expr, exp)
}

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

func checkOpsAllowed(t *testing.T, typeParams, body, exp string) {
	t.Helper()

	prime := "package main\n" +
		"func Foo[" + typeParams + "](x X, i int) {\n" +
		"	" + body + "\n" +
		"}\n" +
		"func main() {}\n"

	got := `Ok`
	_, err := astTools.ParsePackage(map[string]string{`t.go`: prime}, nil)
	if err != nil {
		got = removePaths(err.Error())
	}

	if diff := cmp.Diff(exp, got); len(diff) > 0 {
		diff = strings.TrimSpace(diff)
		diff = strings.ReplaceAll(diff, "\n", "\n    ")
		t.Errorf("Unexpected results:\ntype params: %s\nbody:        %s\ndiff: %s\ngot: %q\nexp: %q",
			typeParams, body, diff, got, exp)
	}
}

func checkOpsWithFile(t *testing.T, prime string, pos token.Pos, expr, exp string) {
	t.Helper()

	ps, err := astTools.ParsePackage(map[string]string{`t.go`: prime}, nil)
	if err != nil {
		t.Errorf(`Failed to parse package: %v`, err)
		return
	}

	if !pos.IsValid() {
		pos = ps.Syntax[0].FileEnd - 1
	}

	tv, err := types.Eval(ps.Fset, ps.Types, pos, expr)
	if err != nil {
		t.Errorf(`Failed to parse expression: %v`, err)
		return
	}

	got := Ops(tv.Type).String()
	got = strings.ReplaceAll(got, `command-line-arguments`, `$`)

	if diff := cmp.Diff(exp, got); len(diff) > 0 {
		diff = strings.TrimSpace(diff)
		diff = strings.ReplaceAll(diff, "\n", "\n    ")
		t.Errorf("Unexpected results:\ntype: %s\ndiff: %s\ngot: %q\nexp: %q",
			tv.Type.String(), diff, got, exp)
	}
}

func Test_Ops_Arrays(t *testing.T) {
	checkOps(t, `[3]int`,
		`Len|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)

	checkOps(t, `[...]byte{1, 2, 3}`,
		`Len|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:byte, Slice:[]byte, Range1:int, Range2:byte }`)
}

func Test_Ops_PointersToArrays(t *testing.T) {
	checkOps(t, `&[3]int{1,2,3}`,
		`Len|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)

	checkOps(t, `*[3]int`,
		`Len|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)

	checkOps(t, `*[3]byte`,
		`Len|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:byte, Slice:[]byte, Range1:int, Range2:byte }`)

	checkOps(t, `(*int)(nil)`,
		`IsNil|Comparable|Ref`)

	checkOps(t, `**[3]byte`,
		`IsNil|Comparable|Ref`)

	checkOpsWithFile(t, lines(
		`package t`,
		`type Foo [3]int`,
	), token.NoPos, `new(Foo)`,
		`Len|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)
}

func Test_Ops_Slices(t *testing.T) {
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
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1`+ // TODO: Fix, Range1 not valid, Range0 is then also not valid
			`{ Range1:{int | uint} }`)

	checkOpsWithFile(t, lines(
		`package t`,
		`func Foo[T int|uint](t T) {`,
		`	println(t)`,
		`}`,
	), token.Pos(40), `t`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1`+
			`{ Range1:{int | uint} }`)

	checkOps(t, `interface{int64|float64}`,
		`Add|Arith|Comparable|Orderable|Ref`)

	checkOps(t, `interface{String() string; int|uint}`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1`+
			`{ Range1:{int | uint} }`)
}

func Test_Ops_UnionsTypes_WithSharedOps(t *testing.T) {
	checkOps(t, `interface{~[]byte|~string}`,
		`Len|Ref|GetIndex|Slice|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:byte, Slice:{~[]byte | ~string}, Range1:int, Range2:{byte | rune} }`)

	checkOps(t, `interface{~[]byte|string}`,
		`Len|Ref|GetIndex|Slice|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:byte, Slice:{~[]byte | string}, Range1:int, Range2:{byte | rune} }`)

	checkOps(t, `interface{[]byte|~string}`,
		`Len|Ref|GetIndex|Slice|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:byte, Slice:{[]byte | ~string}, Range1:int, Range2:{byte | rune} }`)

	checkOps(t, `interface{[]byte|string}`,
		`Len|Ref|GetIndex|Slice|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:byte, Slice:{[]byte | string}, Range1:int, Range2:{byte | rune} }`)
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
		`Add|Comparable|Orderable|Ref|Range0|Range1`+
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
		`T any, X ~[]T | ~[4]T`,
		`var t T = x[0]; x[1] = t; _ = len(x)`,
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
		`T any, X ~[]int | ~[]int8`,
		`_ = x[0]`,
		`cannot index x (variable of type X constrained by ~[]int | ~[]int8)`)
}

func Test_Ops_Go_Allowance_MapAndMap(t *testing.T) {
	checkOpsAllowed(t,
		`T any, X ~map[int]T | ~map[int]T`,
		`for k, v := range x { _, _ = k, v }`,
		`overlapping terms ~map[int]T and ~map[int]T`)
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

//==============================================================================================
// Meaning of the above confirmations:
//
//	┌──────────────────────────┬──────────────────────────┬────────────────────────────────────┐
//	│ Constraint               │ Allowed                  │ Rejected                           │
//	├──────────────────────────┼──────────────────────────┼────────────────────────────────────┤
//	│ ~[]T | ~map[int]T        │ len, == nil              │ indexing, range, make              │
//	├──────────────────────────┼──────────────────────────┼────────────────────────────────────┤
//	│ ~[]T | ~[4]T             │ indexing, set index, len │ range (different underlying types) │
//	├──────────────────────────┼──────────────────────────┼────────────────────────────────────┤
//	│ ~[]T | ~*[4]T            │ indexing, &x[0]          │                                    │
//	├──────────────────────────┼──────────────────────────┼────────────────────────────────────┤
//	│ ~[]byte | ~string        │ x[0] as byte, slicing    │ range                              │
//	├──────────────────────────┼──────────────────────────┼────────────────────────────────────┤
//	│ ~[]int | ~[]int8         │                          │ indexing (element types differ)    │
//	├──────────────────────────┼──────────────────────────┼────────────────────────────────────┤
//	│ ~map[int]T | ~map[int8]T │                          │ indexing, range (key types differ) │
//	├──────────────────────────┼──────────────────────────┼────────────────────────────────────┤
//	│ ~uint8 | ~int16          │ x + 200                  │ x + 300                            │
//	├──────────────────────────┼──────────────────────────┼────────────────────────────────────┤
//	│ ~int | ~float64          │ x + 1                    │ x + 1.5, %                         │
//	└──────────────────────────┴──────────────────────────┴────────────────────────────────────┘
//
// The rules behind this:
//   - Indexing: every type in the set must have an identical element type.
//     If any type is a map, all of them must be maps with identical key types.
//   - Range and make: every type in the set must have the identical underlying type.
//   - Untyped constants: a constant is accepted only if it can be represented in every type in the set.
//
// Try to encode the per-operation rules in unionTermsOps instead of one generic intersect-then-combine pass:
//
//   - Range and make should require all terms to share one underlying type;
//     checking the underlying types before collecting results is enough.
//
//   - Indexing should require a common Elem.
//
//   - For Key, if any term is a map, require all terms to be maps with a common key. Otherwise Key stays untyped int.
//
//     Why this works. I ran a few more compiler checks, and they all came back the same way. Wherever an operation returns a type,
//     either the terms agree on that type or Go rejects the code:
//
//   - Indexing (x[i], set, &x[i]): every term needs an identical element type. If any term is a map, all terms must be maps with
//     an identical key type. Otherwise the key is an integer, and your castIn handles constants.
//
//   - Range and make: every term needs the identical underlying type. This includes integers: ranging over ~int | ~int8 is
//     rejected, but ranging over ~int is allowed.
//
//   - Slicing: every term needs the identical underlying type. The one exception is ~[]byte | ~string, which is allowed and
//     returns X itself.
//
//   - real, imag and complex: never allowed on a type parameter, even with a single term (see go.dev/issue/50937
//     (https://go.dev/issue/50937)).
//
//   - Operations that don't return a new type, like +, ==, len, == nil and clear, are where your dict.add approach applies. The
//     code compiles, so you only need the set of operations every term supports.
//
//     How I'd restructure unionTermsOps:
//
//     1. All terms share one underlying type. If you strip the tildes, this covers ~[]T and also A | B where both are defined as
//     []T. Return getOps(orig, thatUnderlying) and treat it like a single type. Range, make and slicing are only ever valid in
//     this case.
//     2. The underlying types differ.
//
//   - Intersect the operation bits.
//
//   - Clear Range, Make, Complex and Float unconditionally.
//
//   - Clear Slice unless the terms are exactly []byte and string, in which case Slice is orig.
//
//   - For indexing, keep Elem if it's identical across all terms (using types.Identical), otherwise clear the index
//     operations.
//
//   - Set Key to the shared map key if every term is a map with the same key. If no term is a map, set it to untyped int.
//     Otherwise clear the index operations.
//     3. Drop innerUnionType and NormalTerms from the combining step. Keep NormalTerms only where interfaceTypeOps gets the
//     constraint's terms.
//
//     Bug in the concrete path: integerTypeOps sets Range1: t2, which is the basic type. Ranging over an integer gives you the
//     operand's own type, so type MyInt int or a single-term ~int type parameter should yield MyInt or X, not int. It should be
//     Range1: orig, the same as stringTypeOps already does for Slice.
//
//     Possibly missing: clear(x) compiles for ~[]T | ~map[int]T. If your typeOp set doesn't have a Clear operation yet, you'll
//     probably want one, since slices and maps need different dictionary implementations for it.
//
//     Your untyped int convention fits these rules, and both cases you found are covered. For indexing, Key is only ever untyped int
//     (every term is a slice, array, pointer to array or string) or one identical map key type. Go never allows a mix, so you never
//     need a union for it.
//
//     The new checks also show that most of the mixed-set behavior falls out of two simple steps:
//
//     1. Intersect the operation bits across terms. For ~[]byte | ~string, Go rejects setting an index, taking &x[i], the
//     three-index slice x[a:b:c] and cap. Those are exactly the operations string lacks, so the intersection removes them.
//     Reading x[i] and two-index slices still work.
//     2. Keep a result type only if it's types.Identical across terms.
//
//   - []T | [3]T | *[3]T all have key untyped int and element T, so indexing is kept, which matches the compiler.
//
//   - []byte | string | [4]byte all have element byte, so reading an index is kept.
//
//   - Slicing mostly works out too. sliceTypeOps and stringTypeOps both return orig, so []byte | string keeps Slice: X. An
//     array returns []T, which doesn't match a slice's X, so []T | [3]T drops slicing, which is also correct.
//
//     You need one explicit rule on top of that, because the identity check alone can't catch it. If the terms' underlying types
//     differ, clear Range0, Range1, Range2 and Make:
//
//   - []T | [3]T both report range types int and T, so they'd survive the identity check, but Go rejects ranging over that set.
//
//   - []T | map[int]T both have Make, but Go rejects make for it.
//
//   - The same goes for slicing [3]T | *[3]T: both report []T, but Go rejects it. Since your input compiles, keeping that bit
//     only means the dictionary carries an unused function, so it's harmless either way.
//
//     Bug found along the way: arrayTypeOps is missing typeOp.Cap. The compiler accepts cap(x) for ~[]T | ~[3]T, and cap works on
//     arrays and pointers to arrays. pointerTypeOps reuses arrayTypeOps, so adding it there fixes both. That's in addition to the
//     integerTypeOps fix I mentioned last time: Range1 should be orig, not t2.
