package typeTools

import (
	"go/token"
	"go/types"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Grant-Nelson/Gozer/avail/astTools"
)

func lines(lines ...string) string { return strings.Join(lines, "\n") }

func checkOps(t *testing.T, expr, exp string) {
	t.Helper()
	checkOpsWithFile(t, "package t\n", token.NoPos, expr, exp)
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
		`Clear|GetIndex|IsNil|Len|Make|Make3|Range|Range2|Ref|RefIndex|SetIndex|Slice|Slice3`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)

	checkOps(t, `[3]byte`,
		`ByteSlice|Clear|GetIndex|IsNil|Len|Make|Make3|Range|Range2|Ref|RefIndex|SetIndex|Slice|Slice3`+
			`{ Key:untyped int, Elem:byte, Slice:[]byte, Range1:int, Range2:byte }`)
}

func Test_Ops_PointersToArrays(t *testing.T) {
	checkOps(t, `&[3]int{1,2,3}`,
		`Clear|Comparable|Deref|GetIndex|IsNil|Len|Make|Make3|Orderable|Range|Range2|Ref|RefIndex|SetIndex|Slice|Slice3`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Deref:[3]int, Range1:int, Range2:int }`)

	checkOps(t, `*[3]int`,
		`Clear|Comparable|Deref|GetIndex|IsNil|Len|Make|Make3|Orderable|Range|Range2|Ref|RefIndex|SetIndex|Slice|Slice3`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Deref:[3]int, Range1:int, Range2:int }`)

	checkOps(t, `*[3]byte`,
		`ByteSlice|Clear|Comparable|Deref|GetIndex|IsNil|Len|Make|Make3|Orderable|Range|Range2|Ref|RefIndex|SetIndex|Slice|Slice3`+
			`{ Key:untyped int, Elem:byte, Slice:[]byte, Deref:[3]byte, Range1:int, Range2:byte }`)

	checkOps(t, `(*int)(nil)`,
		`Comparable|Deref|IsNil|Orderable|Ref`+
			`{ Deref:int }`)

	checkOpsWithFile(t, lines(
		`package t`,
		`type Foo [3]int`,
	), token.NoPos, `new(Foo)`,
		`Clear|Comparable|Deref|GetIndex|IsNil|Len|Make|Make3|Orderable|Range|Range2|Ref|RefIndex|SetIndex|Slice|Slice3`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Deref:$.Foo, Range1:int, Range2:int }`)
}

func Test_Ops_Slices(t *testing.T) {
	checkOps(t, `[]int`,
		`Cap|Clear|GetIndex|IsNil|Len|Make|Make3|Range|Range2|Ref|RefIndex|SetIndex|Slice|Slice3`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)

	checkOps(t, `[]*float64`,
		`Cap|Clear|GetIndex|IsNil|Len|Make|Make3|Range|Range2|Ref|RefIndex|SetIndex|Slice|Slice3`+
			`{ Key:untyped int, Elem:*float64, Slice:[]*float64, Range1:int, Range2:*float64 }`)

	checkOps(t, `[]byte`,
		`ByteSlice|Cap|Clear|GetIndex|IsNil|Len|Make|Make3|Range|Range2|Ref|RefIndex|SetIndex|Slice|Slice3`+
			`{ Key:untyped int, Elem:byte, Slice:[]byte, Range1:int, Range2:byte }`)

	checkOps(t, `*[]byte`,
		`Comparable|Deref|IsNil|Orderable|Ref`+
			`{ Deref:[]byte }`)

	checkOpsWithFile(t, lines(
		`package t`,
		`type Bar struct{ v int }`,
		`type Foo []*Bar`,
	), token.NoPos, `Foo`,
		`Cap|Clear|GetIndex|IsNil|Len|Make|Make3|Range|Range2|Ref|RefIndex|SetIndex|Slice|Slice3`+
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
		`Add|Arith|Bitwise|Comparable|Mod|Orderable|Range|Ref`+
			`{ Range1:int }`)

	checkOps(t, `uint`,
		`Add|Arith|Bitwise|Comparable|Mod|Orderable|Range|Ref`+
			`{ Range1:uint }`)

	checkOps(t, `int64`,
		`Add|Arith|Bitwise|Comparable|Mod|Orderable|Range|Ref`+
			`{ Range1:int64 }`)

	checkOps(t, `42`,
		`Add|Arith|Bitwise|Comparable|Mod|Orderable|Range|Ref`+
			`{ Range1:untyped int }`)
}

func Test_Ops_FloatTypes(t *testing.T) {
	checkOps(t, `float32`,
		`Add|Arith|Comparable|Complex|Orderable|Ref`+
			`{ Complex:complex64 }`)

	checkOps(t, `float64`,
		`Add|Arith|Comparable|Complex|Orderable|Ref`+
			`{ Complex:complex128 }`)

	checkOps(t, `1.42`,
		`Add|Arith|Comparable|Complex|Orderable|Ref`+
			`{ Complex:untyped complex }`)
}

func Test_Ops_ComplexTypes(t *testing.T) {
	checkOps(t, `complex128`,
		`Add|Arith|Comparable|RealImag|Ref`+
			`{ RealImag:float64 }`)

	checkOps(t, `complex64`,
		`Add|Arith|Comparable|RealImag|Ref`+
			`{ RealImag:float32 }`)

	checkOps(t, `1.42+6.2i`,
		`Add|Arith|Comparable|RealImag|Ref`+
			`{ RealImag:untyped float }`)
}

func Test_Ops_StringTypes(t *testing.T) {
	checkOps(t, `string`,
		`Add|ByteSlice|Comparable|GetIndex|Len|Orderable|Range|Range2|Ref|Slice`+
			`{ Key:untyped int, Elem:uint8, Slice:string, Range1:int, Range2:rune }`)

	checkOps(t, `"Hello"`,
		`Add|ByteSlice|Comparable|GetIndex|Len|Orderable|Range|Range2|Ref|Slice`+
			`{ Key:untyped int, Elem:uint8, Slice:untyped string, Range1:int, Range2:rune }`)

	checkOpsWithFile(t, lines(
		`package t`,
		`type Foo string`,
	), token.NoPos, `Foo`,
		`Add|ByteSlice|Comparable|GetIndex|Len|Orderable|Range|Range2|Ref|Slice`+
			`{ Key:untyped int, Elem:uint8, Slice:$.Foo, Range1:int, Range2:rune }`)
}

func Test_Ops_RuneTypes(t *testing.T) {
	checkOps(t, `rune`,
		`Add|Arith|Bitwise|Comparable|Mod|Orderable|Range|Ref`+
			`{ Range1:rune }`)

	checkOps(t, `'☺'`,
		`Add|Arith|Bitwise|Comparable|Mod|Orderable|Range|Ref`+
			`{ Range1:untyped rune }`)
}

func Test_Ops_PointersTypes(t *testing.T) {
	checkOpsWithFile(t, lines(
		`package t`,
		`import "unsafe"`,
		`var x = 10`,
		`var p = unsafe.Pointer(&x)`,
	), token.NoPos, `p`, `Comparable|IsNil|Orderable|Ref`)

	checkOps(t, `nil`, `IsNil|Ref`)
}

func Test_Ops_ChanTypes(t *testing.T) {
	checkOps(t, `chan int`,
		`Cap|Comparable|IsNil|Len|Range|Recv|Send`+
			`{ Range1:int }`)

	checkOps(t, `make(chan int, 3)`,
		`Cap|Comparable|IsNil|Len|Range|Recv|Send`+
			`{ Range1:int }`)

	checkOps(t, `new(chan int)`,
		`Comparable|Deref|IsNil|Orderable|Ref`+
			`{ Deref:chan int }`)

	checkOps(t, `chan<- int`,
		`Cap|Comparable|IsNil|Len|Send`)

	checkOps(t, `<-chan int`,
		`Cap|Comparable|IsNil|Len|Range|Recv`+
			`{ Range1:int }`)
}

func Test_Ops_UnionsTypes(t *testing.T) {
	checkOps(t, `interface{String() string}`, `IsNil`)

	checkOps(t, `interface{int|uint}`,
		`Add|Arith|Bitwise|Comparable|Mod|Orderable|Ref`)

	checkOpsWithFile(t, lines(
		`package t`,
		`func Foo[T int|uint](t T) {`,
		`	println(t)`,
		`}`,
	), token.Pos(40), `t`,
		`Add|Arith|Bitwise|Comparable|Mod|Orderable|Ref`)

	checkOps(t, `interface{int64|float64}`,
		`Add|Arith|Comparable|Orderable|Ref`)

	checkOps(t, `interface{String() string; int|uint}`,
		`Add|Arith|Bitwise|Comparable|Mod|Orderable|Ref`)

	checkOps(t, `interface{~[]byte|~string}`,
		`ByteSlice|GetIndex|Len|Range|Ref|Slice`+
			`{ Key:untyped int, Range1:int }`)

	checkOps(t, `interface{~[]int|~[3]int}`,
		`Clear|GetIndex|IsNil|Len|Make|Make3|Ref|RefIndex|SetIndex|Slice|Slice3`+
			`{ Key:untyped int, Elem:int, Range1:int, Range2:int }`)

	checkOps(t, `interface{~[]int|~*[3]int}`,
		`Clear|GetIndex|IsNil|Len|Make|Make3|Ref|RefIndex|SetIndex|Slice|Slice3`+
			`{ Key:untyped int, Elem:int, Range1:int, Range2:int }`)

	checkOpsWithFile(t, lines(
		`package t`,
		`func Foo[T any, S ~[]T](s S) {`,
		`	println(s)`,
		`}`,
	), token.Pos(45), `s`,
		`Cap|Clear|GetIndex|IsNil|Len|Make|Make3|Range|Range2|Ref|RefIndex|SetIndex|Slice|Slice3`+
			`{ Key:untyped int, Elem:T, Slice:S, Range1:int, Range2:T }`)

	checkOps(t, `interface{~[]byte|~[]int}`,
		`Cap|Clear|IsNil|Len|Make|Make3|Range|Ref`+
			`{ Key:untyped int, Range1:int }`)
}

// TODO: Create a test to check c.Walk has the correct value
//		package main
//
//		import "fmt"
//
//		type Cat struct{ lives int }
//
//		func (c *Cat) Walk(y func(int) bool) {
//			for i := range c.lives {
//				if !y(i) {
//					return
//				}
//			}
//		}
//
//		func main() {
//			c := &Cat{lives: 3}
//			for i := range c.Walk {
//				fmt.Println(i)
//			}
//		}
