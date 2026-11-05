package typeTools

import (
	"fmt"
	"go/token"
	"go/types"
	"regexp"
	"slices"
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

/*
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
	checkGoOps(t, ``, ``, `bool`,
		`Comparable|Ref`)

	checkGoOps(t, `type Foo bool`, ``, `Foo`,
		`Comparable|Ref`)
}

func Test_Ops_IntegerTypes(t *testing.T) {
	checkGoOps(t, ``, ``, `int`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1`+
			`{ Range1:int }`)

	checkGoOps(t, ``, ``, `uint`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1`+
			`{ Range1:uint }`)

	checkGoOps(t, ``, ``, `int64`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1`+
			`{ Range1:int64 }`)
}

func Test_Ops_FloatTypes(t *testing.T) {
	checkGoOps(t, ``, ``, `float32`,
		`Add|Arith|Comparable|Orderable|Ref|Float`+
			`{ Complex:complex64 }`)

	checkGoOps(t, ``, ``, `float64`,
		`Add|Arith|Comparable|Orderable|Ref|Float`+
			`{ Complex:complex128 }`)

	checkGoOps(t, `type Foo float64`, ``, `Foo`,
		`Add|Arith|Comparable|Orderable|Ref|Float`+
			`{ Complex:untyped complex }`)

	checkGoOps(t, `type Foo float64`, ``, `*Foo`,
		`Add|Arith|Comparable|Orderable|Ref|Float`+
			`{ Complex:untyped complex }`)

	checkGoOps(t, `type Foo *float64`, ``, `Foo`,
		`Add|Arith|Comparable|Orderable|Ref|Float`+
			`{ Complex:untyped complex }`)

	checkGoOps(t, `type Foo *float64`, ``, `*Foo`,
		`Add|Arith|Comparable|Orderable|Ref|Float`+
			`{ Complex:untyped complex }`)
}

func Test_Ops_ComplexTypes(t *testing.T) {
	checkGoOps(t, ``, ``, `complex128`,
		`Add|Arith|Comparable|Ref|Complex`+
			`{ Float:float64 }`)

	checkGoOps(t, ``, ``, `complex64`,
		`Add|Arith|Comparable|Ref|Complex`+
			`{ Float:float32 }`)

	// TODO: Fix when we can mark deviations
	//checkGoOps(t, ``, `T interface{float64|complex128}`, `T`,
	//	`Add|Arith|Comparable|Ref`)
}

func Test_Ops_StringTypes(t *testing.T) {
	checkGoOps(t, ``, ``, `string`,
		`Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:uint8, Slice:string, Range1:int, Range2:rune }`)

	checkGoOps(t, `type Foo string`, ``, `Foo`,
		`Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:uint8, Slice:$.Foo, Range1:int, Range2:rune }`)

	checkGoOps(t, ``, `T ~string`, `T`,
		`Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:uint8, Slice:$.Foo, Range1:int, Range2:rune }`)

	checkGoOps(t, `type Foo string`, `T Foo`, `T`,
		`Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2`+
			`{ Key:untyped int, Elem:uint8, Slice:$.Foo, Range1:int, Range2:rune }`)
}

func Test_Ops_RuneTypes(t *testing.T) {
	checkGoOps(t, ``, ``, `rune`,
		`Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1`+
			`{ Range1:rune }`)
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
*/

//================================================================================= ==  =

func lines(lines ...string) string { return strings.Join(lines, "\n") }

func indent(s string) string { return "\t" + strings.ReplaceAll(s, "\n", "\n\t") }

var regexPathMatch = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^.*t\.go:(?:\d+:){0,2}\s`)
})

func gt(typ types.Type) string {
	if typ == nil {
		return `any`
	}
	r := types.Default(typ).String()
	r = strings.ReplaceAll(r, `command-line-arguments.`, ``)
	return r
}

func removePaths(s string) string {
	buf := &strings.Builder{}
	for ln := range strings.Lines(s) {
		ln = regexPathMatch().ReplaceAllString(ln, ``)
		_, _ = buf.WriteString(ln)
	}
	return buf.String()
}

type testee struct {
	t          *testing.T
	prime      string
	typeParams string
	typeExpr   string
	exp        string
	opTypes    OpTypes
}

func (ts *testee) Check() {
	ts.prep()
	ts.getOpTypes()
	ts.checkExpected()
	if !testing.Short() {
		ts.fullBuildCheck()
	}
}

func (ts *testee) prep() {
	if len(ts.typeParams) > 0 && !strings.HasPrefix(ts.typeParams, `[`) {
		ts.typeParams = `[` + ts.typeParams + `]`
	}
}

func (ts *testee) getOpTypes() {
	file := lines(
		`package t`,
		ts.prime,
		`func testFn`+ts.typeParams+`(x `+ts.typeExpr+`) {`,
		``)
	posOff := len(file)
	file += lines(`print(x)`,
		`}`)
	ps, err := astTools.ParsePackage(map[string]string{`t.go`: file}, nil)
	if err != nil {
		ts.t.Fatalf("Failed to parse package: %v\n%s", err, file)
	}

	pos := ps.Syntax[0].FileStart + token.Pos(posOff)
	tv, err := types.Eval(ps.Fset, ps.Types, pos, `x`)
	if err != nil {
		ts.t.Fatalf("Failed to parse expression: %v\n%s", err, file)
	}

	ts.opTypes = Ops(tv.Type)
}

func (ts *testee) checkExpected() {
	got := ts.opTypes.String()
	got = strings.ReplaceAll(got, `command-line-arguments`, `$`)
	if diff := cmp.Diff(ts.exp, got); len(diff) > 0 {
		diff = strings.TrimSpace(diff)
		diff = strings.ReplaceAll(diff, "\n", "\n    ")
		ts.t.Errorf("Unexpected results:\ntype: %s\ndiff: %s\ngot: %q\nexp: %q",
			gt(ts.opTypes.Type), diff, got, ts.exp)
	}
}

func (ts *testee) buildTrial(code ...string) string {
	file := lines(
		`package t`,
		ts.prime,
		`func testFn`+ts.typeParams+`(x, y, z `+gt(ts.opTypes.Type)+`) {`,
		indent(lines(code...)),
		`}`)
	_, err := astTools.ParsePackage(map[string]string{`t.go`: file}, nil)
	if err != nil {
		return removePaths(err.Error())
	}
	return ``
}

func (ts *testee) anyBuilds(buildErrs ...string) string {
	for _, e := range buildErrs {
		if len(e) <= 0 {
			return ``
		}
	}
	return fmt.Sprintf("expected at least one nil-error out of %d results:\n%s", len(buildErrs),
		indent(lines(buildErrs...)))
}

func (ts *testee) allBuilds(buildErrs ...string) string {
	for _, e := range buildErrs {
		if len(e) > 0 {
			buildErrs = slices.DeleteFunc(buildErrs, func(e string) bool { return len(e) <= 0 })
			return fmt.Sprintf("expected no errors out of %d results:\n\t%s", len(buildErrs),
				indent(lines(buildErrs...)))
		}
	}
	return ``
}

func (ts *testee) buildCheck(op typeOp.Op, buildErr string) {
	if len(buildErr) <= 0 {
		if !ts.opTypes.Ops.Any(op) {
			ts.t.Errorf(`%s: %s was allowed by Go but not in the expected ops`, gt(ts.opTypes.Type), op.String())
		}
	} else if ts.opTypes.Ops.Any(op) {
		ts.t.Errorf("%s: %s was in the expected ops but Go errored with\n%s", gt(ts.opTypes.Type), op.String(), indent(buildErr))
	}
}

func (ts *testee) fullBuildCheck() {
	ts.buildCheck(typeOp.Add, ts.buildTrial(
		`x = y + z`,
		`x += y`))

	ts.buildCheck(typeOp.Arith, ts.allBuilds(
		ts.buildTrial(
			`x = -y`,
			`x = y - z`,
			`x = y * z`,
			`x = y / z`,
			`x -= y`,
			`x *= y`,
			`x /= y`,
			`x++`,
			`x--`),
		ts.anyBuilds(
			ts.buildTrial(
				`_ = int(y)`,
				`_ = float64(y)`,
				`_ = `+gt(ts.opTypes.Type)+`(int(24))`,
				`_ = `+gt(ts.opTypes.Type)+`(float64(1.2))`,
			),
			ts.buildTrial(
				`_ = complex64(y)`,
				`_ = complex128(y)`,
				`_ = `+gt(ts.opTypes.Type)+`(complex64(1.2+3.4i))`,
				`_ = `+gt(ts.opTypes.Type)+`(complex128(1.2+3.4i))`,
			),
		),
	))

	ts.buildCheck(typeOp.Mod, ts.buildTrial(
		`x = y % z`,
		`x %= y`))

	ts.buildCheck(typeOp.Bitwise, ts.buildTrial(
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

	ts.buildCheck(typeOp.Len, ts.buildTrial(
		`var _ int = len(y)`))

	ts.buildCheck(typeOp.Cap, ts.buildTrial(
		`var _ int = cap(y)`))

	ts.buildCheck(typeOp.IsNil, ts.buildTrial(
		`var _ bool = y == nil`))

	ts.buildCheck(typeOp.Comparable, ts.buildTrial(
		`var _ bool = y == z`,
		`var _ bool = y != z`))

	ts.buildCheck(typeOp.Orderable, ts.buildTrial(
		`var _ bool = y < z`,
		`var _ bool = y <= z`,
		`var _ bool = y > z`,
		`var _ bool = y >= z`))

	ts.buildCheck(typeOp.Ref, ts.buildTrial(
		`var p *`+gt(ts.opTypes.Type),
		`p = &y`,
		`p = new(y)`,
		`p = new(`+gt(ts.opTypes.Type)+`)`,
		`_ = p`))

	ts.buildCheck(typeOp.Make,
		ts.anyBuilds(
			ts.buildTrial(`_ = make(`+gt(ts.opTypes.Type)+`)`),
			ts.buildTrial(`_ = make(`+gt(ts.opTypes.Type)+`, 1)`),
			ts.buildTrial(`_ = make(`+gt(ts.opTypes.Type)+`, 1, 2)`),
		))

	ts.buildCheck(typeOp.GetIndex, ts.buildTrial(
		`var k `+gt(ts.opTypes.Key),
		`var e `+gt(ts.opTypes.Elem),
		`e = y[k]`,
		`_ = e`))

	ts.buildCheck(typeOp.SetIndex, ts.buildTrial(
		`var k `+gt(ts.opTypes.Key),
		`var e `+gt(ts.opTypes.Elem),
		`x[k] = e`,
		`_ = e`))

	ts.buildCheck(typeOp.RefIndex, ts.buildTrial(
		`var _ *`+gt(ts.opTypes.Elem)+` = &y[0]`))

	ts.buildCheck(typeOp.Slice, ts.buildTrial(
		`var s `+gt(ts.opTypes.Slice),
		`s = y[0:1]`,
		`_ = s`))

	ts.buildCheck(typeOp.Slice3, ts.buildTrial(
		`var s `+gt(ts.opTypes.Slice),
		`s = y[0:1:2]`,
		`_ = s`))

	ts.buildCheck(typeOp.Range0, ts.buildTrial(
		`for range x {`,
		`	print("hi")`,
		`}`))

	ts.buildCheck(typeOp.Range1, ts.buildTrial(
		`var k `+gt(ts.opTypes.Range1),
		`for k = range x {`,
		`	print(k)`,
		`}`,
		`for _ = range x {`,
		`	print(k)`,
		`}`))

	ts.buildCheck(typeOp.Range2, ts.buildTrial(
		`var k `+gt(ts.opTypes.Range1),
		`var v `+gt(ts.opTypes.Range2),
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

	ts.buildCheck(typeOp.Recv, ts.buildTrial(
		`var e `+gt(ts.opTypes.Elem),
		`var ok bool`,
		`e = <- x`,
		`e, ok = <- x`,
		`_, _ = e, ok`))

	ts.buildCheck(typeOp.Send, ts.buildTrial(
		`var e `+gt(ts.opTypes.Elem),
		`x <- e`,
		`close(x)`))

	ts.buildCheck(typeOp.Float, ts.buildTrial(
		`var c `+gt(ts.opTypes.Complex),
		`c = complex(x, y)`,
		`_ = c`))

	ts.buildCheck(typeOp.Complex, ts.buildTrial(
		`var f `+gt(ts.opTypes.Float),
		`f = real(x)`,
		`f = imag(x)`,
		`_ = f`))
}
