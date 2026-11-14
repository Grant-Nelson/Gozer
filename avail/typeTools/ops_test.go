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
			`{ Key:untyped int, Elem:byte, Slice:string, Range1:int, Range2:rune }`,
	}.Check(t)

	testee{
		prime:    `type Foo string`,
		typeExpr: `Foo`,
		exp: `Foo: Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:byte, Slice:Foo, Range1:int, Range2:rune }`,
	}.Check(t)

	testee{
		typeParams: `T ~string`,
		typeExpr:   `T`,
		exp: `T: Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:byte, Slice:T, Range1:int, Range2:rune }`,
	}.Check(t)

	testee{
		prime:      `type Foo string`,
		typeParams: `T Foo`,
		typeExpr:   `T`,
		exp: `T: Add|Len|Comparable|Orderable|Ref|GetIndex|Slice|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:byte, Slice:T, Range1:int, Range2:rune }`,
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

	testee{
		typeParams: `T interface{string|[]byte}`,
		typeExpr:   `T`,
		exp: `T: Len|Ref|GetIndex|Slice` +
			`{ Key:untyped int, Elem:byte, Slice:T }`,
	}.Check(t)

	testee{
		prime:      `type A []byte`,
		typeParams: `T interface{string|A}`,
		typeExpr:   `T`,
		exp: `T: Len|Ref|GetIndex|Slice` +
			`{ Key:untyped int, Elem:byte, Slice:T }`,
	}.Check(t)

	testee{
		prime:      `type A []byte`,
		typeParams: `T interface{~string|A}`,
		typeExpr:   `T`,
		exp: `T: Len|Ref|GetIndex|Slice` +
			`{ Key:untyped int, Elem:byte, Slice:T }`,
	}.Check(t)

	testee{
		prime:      `type A []byte`,
		typeParams: `T interface{A|~string}`,
		typeExpr:   `T`,
		exp: `T: Len|Ref|GetIndex|Slice` +
			`{ Key:untyped int, Elem:byte, Slice:T }`,
	}.Check(t)

	testee{
		prime:      `type A string`,
		typeParams: `T interface{[]byte|A}`,
		typeExpr:   `T`,
		exp: `T: Len|Ref|GetIndex|Slice` +
			`{ Key:untyped int, Elem:byte, Slice:T }`,
	}.Check(t)

	testee{
		prime:      `type A string`,
		typeParams: `T interface{~[]byte|A}`,
		typeExpr:   `T`,
		exp: `T: Len|Ref|GetIndex|Slice` +
			`{ Key:untyped int, Elem:byte, Slice:T }`,
	}.Check(t)

	testee{
		prime: lines(
			`type A string`,
			`type B []byte`),
		typeParams: `T interface{A|B}`,
		typeExpr:   `T`,
		exp: `T: Len|Ref|GetIndex|Slice` +
			`{ Key:untyped int, Elem:byte, Slice:T }`,
	}.Check(t)

	testee{
		prime: lines(
			`type A string`,
			`type B []byte`),
		typeParams: `T interface{A|B|string|[]uint8}`,
		typeExpr:   `T`,
		exp: `T: Len|Ref|GetIndex|Slice` +
			`{ Key:untyped int, Elem:byte, Slice:T }`,
	}.Check(t)
}

func Test_Ops_UnionsTypes_UnionReduction(t *testing.T) {
	// From https://cs.opensource.google/go/x/exp/+/e88cd736:typeparams/normalize.go;l=66
	// > In this example, the structural type restriction of P is ~string|int: A|B
	// > expands to ~string|~[]byte|int|string, which reduces to ~string|~[]byte|int,
	// > which when intersected with C (~string|~int) yields ~string|int.
	testee{
		prime: lines(
			`type A interface{ ~string|~[]byte }`,
			`type B interface{ int|string }`,
			`type C interface { ~string|~int }`,
			`type P interface{ A|B; C }`), // ~string|int
		typeParams: `T P`,
		typeExpr:   `T`,
		exp:        `T: Add|Comparable|Orderable|Ref`,
	}.Check(t)
}

func Test_Ops_UnionsTypes_SlicesAndArrays(t *testing.T) {
	testee{
		typeParams: `T interface{~[][]byte|~[]string}`,
		typeExpr:   `T`,
		exp:        `T: Len|Cap|IsNil|Ref`,
	}.Check(t)

	testee{
		typeParams: `T interface{[]int|[3]int}`,
		typeExpr:   `T`,
		exp: `T: Len|Cap|Ref|GetIndex|SetIndex|RefIndex` +
			`{ Key:untyped int, Elem:int }`,
	}.Check(t)

	testee{
		typeParams: `T interface{~[]int|~[3]int}`,
		typeExpr:   `T`,
		exp: `T: Len|Cap|Ref|GetIndex|SetIndex|RefIndex` +
			`{ Key:untyped int, Elem:int }`,
	}.Check(t)

	testee{
		typeParams: `T interface{~[]int|~*[3]int}`,
		typeExpr:   `T`,
		exp: `T: Len|Cap|IsNil|Ref|GetIndex|SetIndex|RefIndex` +
			`{ Key:untyped int, Elem:int }`,
	}.Check(t)

	testee{
		typeParams: `T interface{~[]byte|~[]int}`,
		typeExpr:   `T`,
		exp:        `T: Len|Cap|IsNil|Ref`,
	}.Check(t)

	testee{
		typeParams: `T interface{~[]int|~[3]*int}`,
		typeExpr:   `T`,
		exp:        `T: Len|Cap|Ref`,
	}.Check(t)
}

func Test_Ops_UnionsTypes_CompoundGenerics(t *testing.T) {
	testee{
		typeParams: `T any, S ~[]T`,
		typeExpr:   `S`,
		exp: `S: Len|Cap|IsNil|Ref|Make2|Make3|GetIndex|SetIndex|RefIndex|Slice|Slice3|Range0|Range1|Range2` +
			`{ Key:untyped int, Elem:T, Slice:S, Range1:int, Range2:T }`,
	}.Check(t)

	testee{
		typeParams: `T any, S ~[]T | ~map[int]T`,
		typeExpr:   `S`,
		exp:        `S: Len|IsNil|Ref`,
	}.Check(t)

	testee{
		typeParams: `K comparable, V ~[]T, T ~int, M ~map[K]V`,
		typeExpr:   `M`,
		exp: `M: Len|IsNil|Ref|Make1|Make2|GetIndex|SetIndex|Range0|Range1|Range2` +
			`{ Key:K, Elem:V, Range1:K, Range2:V }`,
	}.Check(t)
}

func Test_Ops_FunctionsTypes(t *testing.T) {
	testee{
		typeExpr: `func()`,
		exp:      `func(): IsNil|Ref`,
	}.Check(t)

	testee{
		typeExpr: `func(p func())`,
		exp:      `func(p func()): IsNil|Ref`,
	}.Check(t)

	testee{
		typeExpr: `func(p func() bool)`,
		exp:      `func(p func() bool): IsNil|Ref|Range0`,
	}.Check(t)

	testee{
		typeExpr: `func(p func() error)`,
		exp:      `func(p func() error): IsNil|Ref`,
	}.Check(t)

	testee{
		typeExpr: `func(p func(int) bool)`,
		exp: `func(p func(int) bool): IsNil|Ref|Range0|Range1` +
			`{ Range1:int }`,
	}.Check(t)

	testee{
		typeParams: `F func(p func(int) bool)`,
		typeExpr:   `F`,
		exp: `F: IsNil|Ref|Range0|Range1` +
			`{ Range1:int }`,
	}.Check(t)

	testee{
		typeParams: `T any, F func(p func(T) bool)`,
		typeExpr:   `F`,
		exp: `F: IsNil|Ref|Range0|Range1` +
			`{ Range1:T }`,
	}.Check(t)

	testee{
		typeExpr: `func(p func(int, float32) bool)`,
		exp: `func(p func(int, float32) bool): IsNil|Ref|Range0|Range1|Range2` +
			`{ Range1:int, Range2:float32 }`,
	}.Check(t)

	testee{
		typeExpr: `func(p func(int, bool, float32) bool)`,
		exp:      `func(p func(int, bool, float32) bool): IsNil|Ref`,
	}.Check(t)
}

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
	ts.checkSetTypes()
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

func (ts *testee) checkSetTypes() {
	ts.t.Helper()
	if ts.opTypes.Type == nil {
		ts.t.Errorf("Expected the Type field of an OpType should always be non-nil\n"+
			"\tOptType: %s", ts.opTypes.String())
	}

	checkType := func(name string, need bool, typ types.Type) {
		if need != (typ != nil) {
			ts.t.Errorf("Expected a %s field to be set when needed or nil when not needed\n"+
				"\tField name: %s\n\tOpType:     %s\n\tIs needed:  %t\n\tField type: %v",
				name, name, ts.opTypes.String(), need, typ)
		}
	}

	ops := ts.opTypes.Ops
	checkType(`Key`, ops.NeedsKeyType(), ts.opTypes.Key)
	checkType(`Elem`, ops.NeedsElemType(), ts.opTypes.Elem)
	checkType(`Complex`, ops.NeedsComplexType(), ts.opTypes.Complex)
	checkType(`Float`, ops.NeedsFloatType(), ts.opTypes.Float)
	checkType(`Slice`, ops.NeedsSliceType(), ts.opTypes.Slice)
	checkType(`Range1`, ops.NeedsRange1Type(), ts.opTypes.Range1)
	checkType(`Range2`, ops.NeedsRange2Type(), ts.opTypes.Range2)
}

func (ts *testee) checkTestCode() {
	ts.t.Helper()
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

// createOpsTestCode creates a Go program using the current type information from the gotten OpType.
// The program runs all the major operators on the types. Since the types will not allow all the operators,
// when the program is parsed and type checked a list of errors will be returned. The list of errors
// represents the operators that are not allowed. The program is returned as a string with a slice of
// operators. The slice of operators represent which operation is associated with which lines of the
// program so that the line number from an error can be used to determine the operator for that error.
// The slice of operators is zero based and the error line numbers is one based so the line number subtract
// one will be the index into the operator slice for the operator for that line number.
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
