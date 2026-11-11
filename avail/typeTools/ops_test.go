package typeTools

import (
	"go/token"
	"go/types"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/Grant-Nelson/Gozer/avail/typeTools/typeOp"
)

func Test_Ops_Arrays(t *testing.T) {
	testee{
		typeExpr: `[3]int`,
		exp: `[3]int: Len|Cap|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`,
	}.Check(t)

	testee{
		prime:    `type Foo [3]int`,
		typeExpr: `Foo`,
		exp: `Foo: Len|Cap|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`,
	}.Check(t)
}

func Test_Ops_PointersToArrays(t *testing.T) {
	testee{
		typeExpr: `*[3]int`,
		exp: `*[3]int: Len|Cap|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`,
	}.Check(t)

	testee{
		typeExpr: `*[3]byte`,
		exp: `*[3]byte: Len|Cap|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:byte, Slice:[]byte, Range1:int, Range2:byte }`,
	}.Check(t)

	testee{
		typeExpr: `*int`,
		exp:      `*int: IsNil|Comparable|Ref`,
	}.Check(t)

	testee{
		typeExpr: `**[3]byte`,
		exp:      `**[3]byte: IsNil|Comparable|Ref`,
	}.Check(t)

	testee{
		prime:    `type Foo [3]int`,
		typeExpr: `*Foo`,
		exp: `*Foo: Len|Cap|IsNil|Comparable|Ref|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`,
	}.Check(t)
}

func Test_Ops_Slices(t *testing.T) {
	testee{
		typeExpr: `[]int`,
		exp: `[]int: Len|Cap|IsNil|Ref|Make2|Make3|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:int, Slice:[]int, Range1:int, Range2:int }`,
	}.Check(t)

	testee{
		typeExpr: `[]*float64`,
		exp: `[]*float64: Len|Cap|IsNil|Ref|Make2|Make3|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:*float64, Slice:[]*float64, Range1:int, Range2:*float64 }`,
	}.Check(t)

	testee{
		typeExpr: `[]byte`,
		exp: `[]byte: Len|Cap|IsNil|Ref|Make2|Make3|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:byte, Slice:[]byte, Range1:int, Range2:byte }`,
	}.Check(t)

	testee{
		typeExpr: `*[]byte`,
		exp:      `*[]byte: IsNil|Comparable|Ref`,
	}.Check(t)

	testee{
		prime: lines(
			`type Bar struct{ v int }`,
			`type Foo []*Bar`,
		),
		typeExpr: `Foo`,
		exp: `Foo: Len|Cap|IsNil|Ref|Make2|Make3|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:*Bar, Slice:Foo, Range1:int, Range2:*Bar }`,
	}.Check(t)
}

func Test_Ops_BooleanTypes(t *testing.T) {
	testee{
		typeExpr: `bool`,
		exp:      `bool: Comparable|Ref`,
	}.Check(t)

	testee{
		prime:    `type Foo bool`,
		typeExpr: `Foo`,
		exp:      `Foo: Comparable|Ref`,
	}.Check(t)

	testee{
		prime:    `type Foo *bool`,
		typeExpr: `Foo`,
		exp:      `Foo: IsNil|Comparable|Ref`,
	}.Check(t)

	testee{
		prime:    `type Foo bool`,
		typeExpr: `*Foo`,
		exp:      `*Foo: IsNil|Comparable|Ref`,
	}.Check(t)
}

func Test_Ops_IntegerTypes(t *testing.T) {
	testee{
		typeExpr: `int`,
		exp: `int: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1` +
			`{ Range1:int }`,
	}.Check(t)

	testee{
		typeExpr: `uint`,
		exp: `uint: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1` +
			`{ Range1:uint }`,
	}.Check(t)

	testee{
		typeExpr: `int8`,
		exp: `int8: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1` +
			`{ Range1:int8 }`,
	}.Check(t)

	testee{
		typeExpr: `int16`,
		exp: `int16: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1` +
			`{ Range1:int16 }`,
	}.Check(t)

	testee{
		typeExpr: `int32`,
		exp: `int32: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1` +
			`{ Range1:int32 }`,
	}.Check(t)

	testee{
		typeExpr: `int64`,
		exp: `int64: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1` +
			`{ Range1:int64 }`,
	}.Check(t)

	testee{
		typeExpr: `uint8`,
		exp: `uint8: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1` +
			`{ Range1:uint8 }`,
	}.Check(t)

	testee{
		typeExpr: `uint16`,
		exp: `uint16: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1` +
			`{ Range1:uint16 }`,
	}.Check(t)

	testee{
		typeExpr: `uint32`,
		exp: `uint32: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1` +
			`{ Range1:uint32 }`,
	}.Check(t)

	testee{
		typeExpr: `uint64`,
		exp: `uint64: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1` +
			`{ Range1:uint64 }`,
	}.Check(t)
}

func Test_Ops_FloatTypes(t *testing.T) {
	testee{
		typeExpr: `float32`,
		exp: `float32: Add|Arith|Comparable|Orderable|Ref|Float` +
			`{ Complex:complex64 }`,
	}.Check(t)

	testee{
		typeExpr: `float64`,
		exp: `float64: Add|Arith|Comparable|Orderable|Ref|Float` +
			`{ Complex:complex128 }`,
	}.Check(t)

	testee{
		prime:    `type Foo float64`,
		typeExpr: `Foo`,
		exp: `Foo: Add|Arith|Comparable|Orderable|Ref|Float` +
			`{ Complex:complex128 }`,
	}.Check(t)

	testee{
		prime:    `type Foo float64`,
		typeExpr: `*Foo`,
		exp:      `*Foo: IsNil|Comparable|Ref`,
	}.Check(t)

	testee{
		prime:    `type Foo *float64`,
		typeExpr: `Foo`,
		exp:      `Foo: IsNil|Comparable|Ref`,
	}.Check(t)

	testee{
		prime:    `type Foo *float64`,
		typeExpr: `*Foo`,
		exp:      `*Foo: IsNil|Comparable|Ref`,
	}.Check(t)

	testee{
		// Go doesn't allow complex on mixed types.
		typeParams: `T interface{float32|float64}`,
		typeExpr:   `T`,
		exp:        `T: Add|Arith|Comparable|Orderable|Ref`,
	}.Check(t)
}

func Test_Ops_ComplexTypes(t *testing.T) {
	testee{
		typeExpr: `complex128`,
		exp: `complex128: Add|Arith|Comparable|Ref|Complex` +
			`{ Float:float64 }`,
	}.Check(t)

	testee{
		typeExpr: `complex64`,
		exp: `complex64: Add|Arith|Comparable|Ref|Complex` +
			`{ Float:float32 }`,
	}.Check(t)

	testee{
		typeParams: `T interface{float64|complex128}`,
		typeExpr:   `T`,
		exp:        `T: Add|Arith|Comparable|Ref`,
	}.Check(t)

	testee{
		// Go doesn't allow real or imag on mixed types.
		typeParams: `T interface{complex64|complex128}`,
		typeExpr:   `T`,
		exp:        `T: Add|Arith|Comparable|Ref`,
	}.Check(t)
}

func Test_Ops_StringTypes(t *testing.T) {
	testee{
		typeExpr: `string`,
		exp: `string: Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:uint8, Slice:string, Range1:int, Range2:rune }`,
	}.Check(t)

	testee{
		prime:    `type Foo string`,
		typeExpr: `Foo`,
		exp: `Foo: Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:uint8, Slice:Foo, Range1:int, Range2:rune }`,
	}.Check(t)

	testee{
		typeParams: `T ~string`,
		typeExpr:   `T`,
		exp: `T: Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:uint8, Slice:T, Range1:int, Range2:rune }`,
	}.Check(t)

	testee{
		prime:      `type Foo string`,
		typeParams: `T Foo`,
		typeExpr:   `T`,
		exp: `T: Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:uint8, Slice:T, Range1:int, Range2:rune }`,
	}.Check(t)
}

func Test_Ops_RuneTypes(t *testing.T) {
	testee{
		typeExpr: `rune`,
		exp: `rune: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref|Range0|Range1` +
			`{ Range1:rune }`,
	}.Check(t)
}

func Test_Ops_ChanTypes(t *testing.T) {
	testee{
		typeExpr: `chan int`,
		exp: `chan int: Len|Cap|IsNil|Comparable|Ref|Make1|Make2|Range0|Range1|Recv|Send` +
			`{ Elem:int, Range1:int }`,
	}.Check(t)

	testee{
		typeExpr: `chan<- int`,
		exp: `chan<- int: Len|Cap|IsNil|Comparable|Ref|Make1|Make2|Send` +
			`{ Elem:int }`,
	}.Check(t)

	testee{
		typeExpr: `<-chan int`,
		exp: `<-chan int: Len|Cap|IsNil|Comparable|Ref|Make1|Make2|Range0|Range1|Recv` +
			`{ Elem:int, Range1:int }`,
	}.Check(t)
}

func Test_Ops_UnionsTypes_Basics(t *testing.T) {
	testee{
		typeExpr: `interface{String() string}`,
		exp:      `interface{String() string}: IsNil|Comparable|Ref`,
	}.Check(t)

	// The above will wrap the argument in the pseudo interface object that is always comparable and
	// can compare against nil, however, below's generic will instantiate to the type of the argument,
	// meaning it may or may not be comparable or compare against nil, unless additional constraints
	// are added to the Go type.

	testee{
		typeParams: `T interface{String() string}`,
		typeExpr:   `T`,
		exp:        `T: Ref`,
	}.Check(t)

	testee{
		typeParams: `T interface{comparable; String() string}`,
		typeExpr:   `T`,
		exp:        `T: Comparable|Ref`,
	}.Check(t)

	testee{
		typeParams: `T interface{int|uint}`,
		typeExpr:   `T`,
		exp:        `T: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref`,
	}.Check(t)

	testee{
		typeParams: `T int|uint`,
		typeExpr:   `T`,
		exp:        `T: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref`,
	}.Check(t)

	testee{
		typeParams: `T interface{int64|float64}`,
		typeExpr:   `T`,
		exp:        `T: Add|Arith|Comparable|Orderable|Ref`,
	}.Check(t)

	testee{
		typeParams: `T interface{String() string; int|uint}`,
		typeExpr:   `T`,
		exp:        `T: Add|Arith|Mod|Bitwise|Comparable|Orderable|Ref`,
	}.Check(t)
}

func Test_Ops_UnionsTypes_WithSharedOps(t *testing.T) {
	testee{
		typeParams: `T interface{~[]byte|~string}`,
		typeExpr:   `T`,
		exp: `T: Len|Ref|GetIndex|Slice` +
			`{ Key:untyped int, Elem:byte, Slice:T }`,
	}.Check(t)

	testee{
		typeParams: `T interface{~[]byte|string}`,
		typeExpr:   `T`,
		exp: `T: Len|Ref|GetIndex|Slice` +
			`{ Key:untyped int, Elem:byte, Slice:T }`,
	}.Check(t)

	testee{
		typeParams: `T interface{~string|~[]byte}`,
		typeExpr:   `T`,
		exp: `T: Len|Ref|GetIndex|Slice` +
			`{ Key:untyped int, Elem:byte, Slice:T }`,
	}.Check(t)

	testee{
		typeParams: `T interface{[]byte|~string}`,
		typeExpr:   `T`,
		exp: `T: Len|Ref|GetIndex|Slice` +
			`{ Key:untyped int, Elem:byte, Slice:T }`,
	}.Check(t)

	testee{
		typeParams: `T interface{[]byte|string}`,
		typeExpr:   `T`,
		exp: `T: Len|Ref|GetIndex|Slice` +
			`{ Key:untyped int, Elem:byte, Slice:T }`,
	}.Check(t)
}

/*
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

func gt(typ types.Type) string {
	if typ == nil {
		return `noType` // defined in the code patterns below
	}
	r := types.Default(typ).String()
	r = strings.ReplaceAll(r, `command-line-arguments.`, ``)
	return r
}

type testee struct {
	t          *testing.T
	prime      string
	typeParams string
	typeExpr   string
	exp        string
	opTypes    OpTypes
}

func (ts testee) Check(t *testing.T) {
	ts.t = t
	ts.t.Helper()
	ts.prep()
	ts.getOpTypes()
	ts.checkExpected()
	if !testing.Short() {
		ts.checkTestCode()
	}
}

func (ts *testee) prep() {
	ts.t.Helper()
	if len(ts.typeParams) > 0 && !strings.HasPrefix(ts.typeParams, `[`) {
		ts.typeParams = `[` + ts.typeParams + `]`
	}
}

func (ts *testee) getOpTypes() {
	ts.t.Helper()
	file := lines(
		`package t`,
		`type noType any`,
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
	ts.t.Helper()
	got := ts.opTypes.String()
	got = strings.ReplaceAll(got, `command-line-arguments.`, ``)
	if diff := cmp.Diff(ts.exp, got); len(diff) > 0 {
		diff = strings.TrimSpace(diff)
		diff = strings.ReplaceAll(diff, "\n", "\n    ")
		ts.t.Errorf("Unexpected results:\ntype: %s\ndiff: %s\ngot: %q\nexp: %q",
			gt(ts.opTypes.Type), diff, got, ts.exp)
	}
}

func (ts *testee) checkTestCode() {
	errs := ts.runTestCode()
	for op := range typeOp.Iterate() {
		ts.testCodeErrCheck(op, errs[op])
	}
}

var regexTestCodeErrorMatch = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^.*t\.go:(\d+):(?:\d+:)?\s(.*)$`)
})

func (ts *testee) runTestCode() map[typeOp.Op][]string {
	ts.t.Helper()
	errs := map[typeOp.Op][]string{}
	code, opLines := ts.createOpsTestCode()
	_, err := astTools.ParsePackage(map[string]string{`t.go`: code}, nil)
	if err != nil {
		errMsg := err.Error()
		for ln := range strings.Lines(errMsg) {
			ln = strings.TrimSpace(ln)
			matches := regexTestCodeErrorMatch().FindStringSubmatch(ln)
			if len(matches) == 3 {
				lineNoStr, errMsg := matches[1], matches[2]
				if lineNo, err := strconv.ParseUint(lineNoStr, 10, 32); err == nil {
					op := opLines[lineNo-1]
					if op == typeOp.None {
						ts.t.Error(`error came from typeOp.None location of the test code: `, errMsg)
						continue
					}
					errs[op] = append(errs[op], errMsg)
				}
			}
		}
	}
	return errs
}

func (ts *testee) testCodeErrCheck(op typeOp.Op, errs []string) {
	ts.t.Helper()
	if len(errs) <= 0 {
		if !ts.opTypes.Ops.Any(op) {
			ts.t.Errorf(`%s: %s was allowed by Go but not in the expected ops`, gt(ts.opTypes.Type), op.String())
		}
	} else if ts.opTypes.Ops.Any(op) {
		ts.t.Errorf("%s: %s was in the expected ops but Go errored with:\n%s", gt(ts.opTypes.Type), op.String(), indent(lines(errs...)))
	}
}

func (ts *testee) createOpsTestCode() (string, []typeOp.Op) {
	ts.t.Helper()
	buf := &strings.Builder{}
	currentOp := typeOp.None
	opLines := []typeOp.Op{}

	writeLn := func(parts ...string) {
		for _, s := range parts {
			for range strings.Count(s, "\n") {
				opLines = append(opLines, currentOp)
			}
			_, _ = buf.WriteString(s)
		}
		_, _ = buf.WriteRune('\n')
		opLines = append(opLines, currentOp)
	}

	mark := func(op typeOp.Op) { currentOp = op }

	typeStr := gt(ts.opTypes.Type)
	writeLn(`package t`)
	writeLn(`type noType any`)
	writeLn(ts.prime)
	writeLn(`func testFn`, ts.typeParams, `(x, y, z `, typeStr, `) {`)
	writeLn(`var b bool`)
	writeLn(`var i int`)
	writeLn(`var p *`, typeStr)
	writeLn(`var k `, gt(ts.opTypes.Key))
	writeLn(`var e `, gt(ts.opTypes.Elem))
	writeLn(`var pe *`, gt(ts.opTypes.Elem))
	writeLn(`var s `, gt(ts.opTypes.Slice))
	writeLn(`var rk `, gt(ts.opTypes.Range1))
	writeLn(`var rv `, gt(ts.opTypes.Range2))
	writeLn(`var c `, gt(ts.opTypes.Complex))
	writeLn(`var f `, gt(ts.opTypes.Float))

	mark(typeOp.Add)
	writeLn(`x = y + z`)
	writeLn(`x += y`)

	mark(typeOp.Arith)
	writeLn(`x = -y`)
	writeLn(`x = y - z`)
	writeLn(`x = y * z`)
	writeLn(`x = y / z`)
	writeLn(`x -= y`)
	writeLn(`x *= y`)
	writeLn(`x /= y`)
	writeLn(`x++`)
	writeLn(`x--`)
	//	any(
	//		// TODO: Work on a better way to check and handle casting.
	//		// TODO: Should check `fn[TIn ~int, TOut ~string](x TIn) { return TOut(x) }`
	//		build(
	//			`yi := int(y)`,
	//			`yf := float64(y)`,
	//			`_ = `+gt(ts.opTypes.Type)+`(yi)`,
	//			`_ = `+gt(ts.opTypes.Type)+`(yf)`,
	//		),
	//		build(
	//			`y1 := complex64(y)`,
	//			`y2 := complex128(y)`,
	//			`_ = `+gt(ts.opTypes.Type)+`(y1)`,
	//			`_ = `+gt(ts.opTypes.Type)+`(y2)`,
	//		),
	//	)

	mark(typeOp.Mod)
	writeLn(`x = y % z`)
	writeLn(`x %= y`)

	mark(typeOp.Bitwise)
	writeLn(`x = ^y`)
	writeLn(`x = y & z`)
	writeLn(`x = y | z`)
	writeLn(`x = y ^ z`)
	writeLn(`x = y &^ z`)
	writeLn(`x = y << z`)
	writeLn(`x = y >> z`)
	writeLn(`x &= y`)
	writeLn(`x |= y`)
	writeLn(`x ^= y`)
	writeLn(`x &^= y`)
	writeLn(`x <<= y`)
	writeLn(`x >>= y`)

	mark(typeOp.Len)
	writeLn(`i = len(y)`)

	mark(typeOp.Cap)
	writeLn(`i = cap(y)`)

	mark(typeOp.IsNil)
	writeLn(`b = y == nil`)

	mark(typeOp.Comparable)
	writeLn(`b = y == z`)
	writeLn(`b = y != z`)

	mark(typeOp.Orderable)
	writeLn(`b = y < z`)
	writeLn(`b = y <= z`)
	writeLn(`b = y > z`)
	writeLn(`b = y >= z`)

	mark(typeOp.Ref)
	writeLn(`p = &y`)
	writeLn(`p = new(y)`)
	writeLn(`p = new(`, typeStr, `)`)

	mark(typeOp.Make1)
	writeLn(`_ = make(` + typeStr + `)`)

	mark(typeOp.Make2)
	writeLn(`_ = make(` + typeStr + `, 1)`)

	mark(typeOp.Make3)
	writeLn(`_ = make(` + typeStr + `, 1, 2)`)

	mark(typeOp.GetIndex)
	writeLn(`e = y[k]`)

	mark(typeOp.SetIndex)
	writeLn(`x[k] = e`)

	mark(typeOp.RefIndex)
	writeLn(`pe = &y[0]`)

	mark(typeOp.Slice)
	writeLn(`s = y[0:1]`)

	mark(typeOp.Slice3)
	writeLn(`s = y[0:1:2]`)

	mark(typeOp.Range0)
	writeLn(`for range x {`)
	writeLn(`	print("hi")`)
	writeLn(`}`)

	mark(typeOp.Range1)
	writeLn(`for rk = range x {`)
	writeLn(`	print(rk)`)
	writeLn(`}`)
	writeLn(`for _ = range x {`)
	writeLn(`	print(rk)`)
	writeLn(`}`)

	mark(typeOp.Range2)
	writeLn(`for rk, rv = range x {`)
	writeLn(`	print(rk, rv)`)
	writeLn(`}`)
	writeLn(`for _, rv = range x {`)
	writeLn(`	print(rk, rv)`)
	writeLn(`}`)
	writeLn(`for rk, _ = range x {`)
	writeLn(`	print(rk, rv)`)
	writeLn(`}`)
	writeLn(`for _, _ = range x {`)
	writeLn(`	print(rk, rv)`)
	writeLn(`}`)

	mark(typeOp.Recv)
	writeLn(`e = <- x`)
	writeLn(`e, b = <- x`)

	mark(typeOp.Send)
	writeLn(`x <- e`)
	writeLn(`close(x)`)

	mark(typeOp.Float)
	writeLn(`c = complex(x, y)`)

	mark(typeOp.Complex)
	writeLn(`f = real(x)`)
	writeLn(`f = imag(x)`)

	mark(typeOp.None)
	writeLn(`_, _ = b, i`)
	writeLn(`_, _, _, _ = p, k, e, pe`)
	writeLn(`_, _, _ = s, rk, rv`)
	writeLn(`_, _ = c, f`)
	writeLn(`}`)
	return buf.String(), opLines
}
