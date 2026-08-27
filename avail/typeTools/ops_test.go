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
	if diff := cmp.Diff(exp, got); len(diff) > 0 {
		diff = strings.TrimSpace(diff)
		diff = strings.ReplaceAll(diff, "\n", "\n    ")
		t.Errorf("Unexpected results:\ntype: %s\ndiff: %s\ngot: %q\nexp: %q",
			tv.Type.String(), diff, got, exp)
	}
}

func Test_Ops_ArraysAndSlices(t *testing.T) {
	checkOps(t, `[3]int`, `Clear|GetIndex|IsNil|Len|Make|Make3|Ref|RefIndex|SetIndex|Slice|Slice3`)
	checkOps(t, `[3]byte`, `ByteSlice|Clear|GetIndex|IsNil|Len|Make|Make3|Ref|RefIndex|SetIndex|Slice|Slice3`)

	checkOps(t, `&[3]int{1,2,3}`, `Clear|Comparable|Deref|GetIndex|IsNil|Len|Make|Make3|Orderable|Ref|RefIndex|SetIndex|Slice|Slice3`)
	checkOps(t, `*[3]int`, `Clear|Comparable|Deref|GetIndex|IsNil|Len|Make|Make3|Orderable|Ref|RefIndex|SetIndex|Slice|Slice3`)
	checkOps(t, `*[3]byte`, `ByteSlice|Clear|Comparable|Deref|GetIndex|IsNil|Len|Make|Make3|Orderable|Ref|RefIndex|SetIndex|Slice|Slice3`)
	checkOps(t, `(*int)(nil)`, `Comparable|Deref|IsNil|Orderable|Ref`)

	checkOpsWithFile(t, lines(
		`package t`,
		`type Foo [3]int`,
	), token.NoPos, `new(Foo)`, `Clear|Comparable|Deref|GetIndex|IsNil|Len|Make|Make3|Orderable|Ref|RefIndex|SetIndex|Slice|Slice3`)

	checkOps(t, `[]int`, `Cap|Clear|GetIndex|IsNil|Len|Make|Make3|Ref|RefIndex|SetIndex|Slice|Slice3`)
	checkOps(t, `[]byte`, `ByteSlice|Cap|Clear|GetIndex|IsNil|Len|Make|Make3|Ref|RefIndex|SetIndex|Slice|Slice3`)
	checkOps(t, `*[]byte`, `Comparable|Deref|IsNil|Orderable|Ref`)
}

func Test_Ops_BasicTypes(t *testing.T) {
	checkOps(t, `bool`, `Comparable|Ref`)
	checkOps(t, `false`, `Comparable|Ref`)

	checkOps(t, `int`, `Add|Arith|Bitwise|Comparable|Mod|Orderable|Ref`)
	checkOps(t, `uint`, `Add|Arith|Bitwise|Comparable|Mod|Orderable|Ref`)
	checkOps(t, `int64`, `Add|Arith|Bitwise|Comparable|Mod|Orderable|Ref`)
	checkOps(t, `42`, `Add|Arith|Bitwise|Comparable|Mod|Orderable|Ref`)

	checkOps(t, `float64`, `Add|Arith|Comparable|Complex|Orderable|Ref`)
	checkOps(t, `1.42`, `Add|Arith|Comparable|Complex|Orderable|Ref`)

	checkOps(t, `complex128`, `Add|Arith|Comparable|RealImag|Ref`)
	checkOps(t, `1.42+6.2i`, `Add|Arith|Comparable|RealImag|Ref`)

	checkOps(t, `string`, `Add|ByteSlice|Comparable|GetIndex|Len|Orderable|Ref|Slice`)
	checkOps(t, `"Hello"`, `Add|ByteSlice|Comparable|GetIndex|Len|Orderable|Ref|Slice`)

	checkOps(t, `rune`, `Add|Arith|Bitwise|Comparable|Mod|Orderable|Ref`)
	checkOps(t, `'☺'`, `Add|Arith|Bitwise|Comparable|Mod|Orderable|Ref`)

	checkOpsWithFile(t, lines(
		`package t`,
		`import "unsafe"`,
		`var x = 10`,
		`var p = unsafe.Pointer(&x)`,
	), token.NoPos, `p`, `Comparable|IsNil|Orderable|Ref`)

	checkOps(t, `nil`, `IsNil|Ref`)
}

func Test_Ops_ChanTypes(t *testing.T) {
	checkOps(t, `chan int`, `Cap|Comparable|IsNil|Len|Recv|Send`)
	checkOps(t, `make(chan int, 3)`, `Cap|Comparable|IsNil|Len|Recv|Send`)
	checkOps(t, `new(chan int)`, `Comparable|Deref|IsNil|Orderable|Ref`)
	checkOps(t, `chan<- int`, `Cap|Comparable|IsNil|Len|Send`)
	checkOps(t, `<-chan int`, `Cap|Comparable|IsNil|Len|Recv`)
}

func Test_Ops_UnionsTypes(t *testing.T) {
	checkOps(t, `interface{String() string}`, `IsNil`)
	checkOps(t, `interface{int|uint}`, `Add|Arith|Bitwise|Comparable|Mod|Orderable|Ref`)

	checkOpsWithFile(t, lines(
		`package t`,
		`func Foo[T int|uint](t T) {`,
		`	println(t)`,
		`}`,
	), token.Pos(40), `t`, `Add|Arith|Bitwise|Comparable|Mod|Orderable|Ref`)

	checkOps(t, `interface{int64|float64}`, `Add|Arith|Comparable|Orderable|Ref`)
	checkOps(t, `interface{String() string; int|uint}`, `Add|Arith|Bitwise|Comparable|Mod|Orderable|Ref`)
	checkOps(t, `interface{~[]byte|~string}`, `ByteSlice|GetIndex|Len|Ref|Slice`)
	checkOps(t, `interface{~[]int|~[3]int}`, `Clear|GetIndex|IsNil|Len|Make|Make3|Ref|RefIndex|SetIndex|Slice|Slice3`)
	checkOps(t, `interface{~[]int|~*[3]int}`, `Clear|GetIndex|IsNil|Len|Make|Make3|Ref|RefIndex|SetIndex|Slice|Slice3`)

	checkOpsWithFile(t, lines(
		`package t`,
		`func Foo[T any, S ~[]T](s S) {`,
		`	println(s)`,
		`}`,
	), token.Pos(45), `s`, `Cap|Clear|GetIndex|IsNil|Len|Make|Make3|Ref|RefIndex|SetIndex|Slice|Slice3`)

	// TODO: Fix, needs to check value and other types.
	checkOps(t, `interface{~[]byte|~[]int}`, `Cap|Clear|IsNil|Len|Make|Make3|Ref`)
}
