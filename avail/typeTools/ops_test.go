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
	checkGoOps(t, ``, ``, `[3]int`,
		`[3]int: Len|Cap|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)

	checkGoOps(t, `type Foo [3]int`, ``, `Foo`,
		`Foo: Len|Cap|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)
}

func Test_Ops_PointersToArrays(t *testing.T) {
	checkGoOps(t, ``, ``, `*[3]int`,
		`Len|Cap|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)

	checkGoOps(t, ``, ``, `*[3]byte`,
		`Len|Cap|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:byte, Slice:[]byte, Range1:int, Range2:byte }`)

	checkGoOps(t, ``, ``, `*int`,
		`IsNil|Comparable|Ref`)

	checkGoOps(t, ``, ``, `**[3]byte`,
		`IsNil|Comparable|Ref`)

	checkGoOps(t, `type Foo [3]int`, ``, `*Foo`,
		`Len|Cap|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)
}

func Test_Ops_Slices(t *testing.T) {
	checkGoOps(t, ``, ``, `[]int`,
		`Len|Cap|IsNil|Ref|Make|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`)

	checkGoOps(t, ``, ``, `[]*float64`,
		`Len|Cap|IsNil|Ref|Make|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:*float64, Slice:[]*float64, Range1:int, Range2:*float64 }`)

	checkGoOps(t, ``, ``, `[]byte`,
		`Len|Cap|IsNil|Ref|Make|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:byte, Slice:[]byte, Range1:int, Range2:byte }`)

	checkGoOps(t, ``, ``, `*[]byte`,
		`IsNil|Comparable|Ref`)

	checkGoOps(t, lines(
		`type Bar struct{ v int }`,
		`type Foo []*Bar`,
	), ``, `Foo`,
		`Len|Cap|IsNil|Ref|Make|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:*$.Bar, Slice:$.Foo, Range1:int, Range2:*$.Bar }`)
}

func Test_Ops_BooleanTypes(t *testing.T) {
	checkOps(t, `bool`, `Comparable|Ref`)
	checkOps(t, `false`, `Comparable|Ref`)

	checkOpsWithFile(t, lines(
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
		`func Foo[T any, S ~[]T](s S) {`,
		`	println(s)`,
		`}`,
	), token.Pos(45), `s`,
		`Len|Cap|IsNil|Ref|Make|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:T, Slice:S, Range1:int, Range2:T }`)

	checkOpsWithFile(t, lines(
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

func lines(lines ...string) string { return strings.Join(lines, "\n") }

var regexPathMatch = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^.*t.go:\s`)
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

	file := "package t\n" + prime
	ps, err := astTools.ParsePackage(map[string]string{`t.go`: file}, nil)
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
	checkOpsWithFile(t, ``, token.NoPos, expr, exp)
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

func checkGoOps(t *testing.T, prime, typeParams, typeExpr, exp string) {
	t.Helper()

	if len(typeParams) > 0 {
		typeParams = `[` + typeParams + `]`
	}

	file := lines(
		`package t`,
		prime,
		`func testFn`+typeParams+`(x `+typeExpr+`) {`,
		``)
	posOff := len(file)
	file += lines(`print(x)`,
		`}`)
	ps, err := astTools.ParsePackage(map[string]string{`t.go`: file}, nil)
	if err != nil {
		t.Fatalf("Failed to parse package: %v\n%s", err, file)
	}

	pos := ps.Syntax[0].FileStart + token.Pos(posOff)
	tv, err := types.Eval(ps.Fset, ps.Types, pos, `x`)
	if err != nil {
		t.Fatalf("Failed to parse expression: %v\n%s", err, file)
	}

	expOp := Ops(tv.Type)

	gt := func(typ types.Type) string {
		if typ == nil {
			return `any`
		}
		r := types.Default(typ).String()
		r = strings.ReplaceAll(r, `command-line-arguments.`, ``)
		return r
	}

	build := func(code string) error {
		t.Helper()
		file := lines(
			`package t`,
			prime,
			`func testFn`+typeParams+`(x, y, z `+gt(expOp.Type)+`) {`,
			"\t"+strings.ReplaceAll(code, "\n", "\n\t"),
			`}`)
		_, err := astTools.ParsePackage(map[string]string{`t.go`: file}, nil)
		return err
	}

	check := func(op typeOp.Op, code string) {
		t.Helper()
		err := build(code)
		if err == nil {
			if !expOp.Ops.Any(op) {
				t.Errorf(`%v was allowed by Go but not in the expected ops`, op.String())
			}
		} else if expOp.Ops.Any(op) {
			errStr := removePaths(err.Error())
			t.Errorf("%v was in the expected ops but Go errored with\n%s\n%s", op.String(), errStr, file)
		}
	}

	check(typeOp.Add, lines(
		`x = y + z`,
		`x += y`))
	check(typeOp.Arith, lines(
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
		`_ = float64(y)`))
	check(typeOp.Mod, lines(
		`x = y % z`,
		`x %= y`))
	check(typeOp.Bitwise, lines(
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
		`x >>= y`))
	check(typeOp.Len, lines(
		`var _ int = len(y)`))
	check(typeOp.Cap, lines(
		`var _ int = cap(y)`))
	check(typeOp.IsNil, lines(
		`var _ bool = y == nil`))
	check(typeOp.Comparable, lines(
		`var _ bool = y == z`,
		`var _ bool = y != z`))
	check(typeOp.Orderable, lines(
		`var _ bool = y < z`,
		`var _ bool = y <= z`,
		`var _ bool = y > z`,
		`var _ bool = y >= z`))
	check(typeOp.Ref, lines(
		`var p *`+gt(expOp.Type),
		`p = &y`,
		`p = new(y)`,
		`p = new(`+gt(expOp.Type)+`)`,
		`_ = p`))

	check(typeOp.Make, lines(
		`_ = make(`+gt(expOp.Type)+`)`,
		`_ = make(`+gt(expOp.Type)+`, 1)`,
		`_ = make(`+gt(expOp.Type)+`, 1, 2)`))

	check(typeOp.GetIndex, lines(
		`var k `+gt(expOp.Key),
		`var e `+gt(expOp.Elem),
		`e = y[k]`,
		`_ = e`))
	check(typeOp.SetIndex, lines(
		`var k `+gt(expOp.Key),
		`var e `+gt(expOp.Elem),
		`x[k] = e`,
		`_ = e`))
	check(typeOp.RefIndex, lines(
		`var _ *`+gt(expOp.Elem)+` = &y[0]`))
	check(typeOp.Slice, lines(
		`var s `+gt(expOp.Slice),
		`s = y[0:1]`,
		`_ = s`))
	check(typeOp.Slice3, lines(
		`var s `+gt(expOp.Slice),
		`s = y[0:1:2]`,
		`_ = s`))
	check(typeOp.Range0, lines(
		`for range x {`,
		`	print("hi")`,
		`}`))
	check(typeOp.Range1, lines(
		`var k `+gt(expOp.Range1),
		`for k = range x {`,
		`	print(k)`,
		`}`,
		`for _ = range x {`,
		`	print(k)`,
		`}`))
	check(typeOp.Range2, lines(
		`var k `+gt(expOp.Range1),
		`var v `+gt(expOp.Range2),
		`for k, v = range x {`,
		`	print(k, v)`,
		`}`,
		`for _, v = range x {`,
		`	print(k, v)`,
		`}`,
		`for k, _ = range x {`,
		`	print(k, v)`,
		`}`,
		`for _, _ = range x {`,
		`	print(k, v)`,
		`}`))
	check(typeOp.Recv, lines(
		`var e `+gt(expOp.Elem),
		`var ok bool`,
		`e = <- x`,
		`e, ok = <- x`,
		`_, _ = e, ok`))
	check(typeOp.Send, lines(
		`var e `+gt(expOp.Elem),
		`x <- e`,
		`close(x)`))
	check(typeOp.Float, lines(
		`var c `+gt(expOp.Complex),
		`c = complex(x, y)`,
		`_ = c`))
	check(typeOp.Complex, lines(
		`var f `+gt(expOp.Float),
		`f = real(x)`,
		`f = imag(x)`,
		`_ = f`))
}
