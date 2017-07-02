package expressions

import (
	"strings"
	"testing"

	"github.com/grant-nelson/Gozer/common"
	"github.com/grant-nelson/Gozer/constructs/types"
)

func TestAssignment(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*AssignmentExp)(nil).String(), "nil")
	CheckReturn(t, (*AssignmentExp)(nil), "nil")
	id := Identifier("age", types.Float32())
	val := Literal("42.3", types.Float32())
	assm := Assignment(id, val)
	CheckExp(t, assm, `age = 42.3`)
	CheckReturn(t, assm, `float32`)
}

func TestBinaryOp(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*BinaryOpExp)(nil).String(), "nil")
	CheckReturn(t, (*BinaryOpExp)(nil), "nil")
	id := Identifier("value", types.Float32())
	val1 := Literal("2.2", types.Float32())
	val2 := Literal("-1.7", types.Int32())
	bin := BinaryOp(BinaryOp(val1, id, AddOp, types.Float32()), val2, MultiplyOp, types.Float32())
	CheckExp(t, bin, `((2.2 + value) * -1.7)`)
	CheckReturn(t, bin, `float32`)
}

func TestCall(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*CallExp)(nil).String(), "nil")
	CheckReturn(t, (*CallExp)(nil), "nil")
	log := types.Class()
	print := log.AddFunction("print").AddParam("msg", types.String()).SetReturn(types.Int())
	id := Identifier("log", log)
	CheckExp(t, id, `log`)
	CheckReturn(t, id,
		`class{`,
		`  int print(string msg)`,
		`}`)

	pfun := Selector(id, "print", print)
	CheckExp(t, pfun, `log.print`)
	CheckReturn(t, pfun, "print")

	param := Literal("Hello World", types.String())
	CheckExp(t, param, `Hello World`)
	CheckReturn(t, param, "string")

	call1 := Call(print, pfun, []Expression{param})
	CheckExp(t, call1, `log.print(Hello World)`)
	CheckReturn(t, call1, "int")

	call2 := Call(nil, nil, []Expression{param})
	CheckExp(t, call2, `nil(Hello World)`)
	CheckReturn(t, call2, "nil")
}

func TestCompoundLit(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*CompoundLiteralExp)(nil).String(), "nil")
	CheckReturn(t, (*CompoundLiteralExp)(nil), "nil")
	val1 := Literal("1.2", types.Float32())
	val2 := Literal("2.3", types.Float32())
	val3 := Literal("3.4", types.Float32())
	comp := CompoundLiteral([]Expression{val1, val2, val3}, types.List(types.Float32()))
	CheckExp(t, comp, `[]float32{1.2, 2.3, 3.4}`)
	CheckReturn(t, comp, `[]float32`)
}

func TestLambda(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*LambdaExp)(nil).String(), "nil")
	CheckReturn(t, (*LambdaExp)(nil), "nil")
	f1 := types.Function().AddParam("a", types.Int()).
		AddParam("b", types.Float64()).AddParam("c", types.Bool()).
		SetReturn(types.Int())
	f1.Body = "{ doSomething() }"
	lam1 := Lambda(f1)
	CheckExp(t, lam1, `int func(int a, float64 b, bool c) { doSomething() }`)
	CheckReturn(t, lam1, `int func(int a, float64 b, bool c)`)
}

func TestMake(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*MakeExp)(nil).String(), "nil")
	CheckReturn(t, (*MakeExp)(nil), "nil")
	m1 := Make()
	CheckExp(t, m1, `make(void)`)
	CheckReturn(t, m1, `void`)

	m1.Type = types.List(types.Int())
	CheckExp(t, m1, `make([]int)`)
	CheckReturn(t, m1, `[]int`)

	m1.Length = Literal("12", types.Int())
	CheckExp(t, m1, `make([]int, 12)`)
	CheckReturn(t, m1, `[]int`)

	m1.Capacity = Literal("36", types.Int())
	CheckExp(t, m1, `make([]int, 12, 36)`)
	CheckReturn(t, m1, `[]int`)

	m1.Length = nil
	CheckExp(t, m1, `make([]int, 0, 36)`)
	CheckReturn(t, m1, `[]int`)
}

func TestDefinition(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*DefinitionExp)(nil).String(), "nil")
	CheckReturn(t, (*DefinitionExp)(nil), "nil")
	id := Identifier("dogName", types.String())
	def := Definition(id, Literal(`"Gizmo"`, types.String()))
	CheckExp(t, def, `string dogName = "Gizmo"`)
	CheckReturn(t, def, `string`)
}

func TestIdentifier(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*IdentifierExp)(nil).String(), "nil")
	CheckReturn(t, (*IdentifierExp)(nil), "nil")
	id := Identifier("happy", types.Bool())
	CheckExp(t, id, `happy`)
	CheckReturn(t, id, `bool`)
}

func TestIndexer(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*IndexerExp)(nil).String(), "nil")
	CheckReturn(t, (*IndexerExp)(nil), "nil")
	index := Literal("4", types.Int())

	notIndexable := Indexer(Identifier("notIndexable", types.Int64()), index)
	CheckExp(t, notIndexable, `notIndexable[4]`)
	CheckReturn(t, notIndexable, `nil`)

	listIndex := Indexer(Identifier("list", types.List(types.Int64())), index)
	CheckExp(t, listIndex, `list[4]`)
	CheckReturn(t, listIndex, `int64`)

	mapIndex := Indexer(Identifier("map", types.Map(types.Int(), types.Float32())), index)
	CheckExp(t, mapIndex, `map[4]`)
	CheckReturn(t, mapIndex, `float32`)

	strIndex := Indexer(Identifier("str", types.String()), index)
	CheckExp(t, strIndex, `str[4]`)
	CheckReturn(t, strIndex, `uint8`)
}

func TestLiteral(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*LiteralExp)(nil).String(), "nil")
	CheckReturn(t, (*LiteralExp)(nil), "nil")
	lit := Literal("Hello World", types.String())
	CheckExp(t, lit, `Hello World`)
	CheckReturn(t, lit, `string`)
}

func TestSelector(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*SelectorExp)(nil).String(), "nil")
	CheckReturn(t, (*SelectorExp)(nil), "nil")
	dog := types.Interface()
	bark := dog.AddFunction("bark").AddParam("volume", types.Float64())
	id := Identifier("dog", dog)
	pfun := Selector(id, "bark", bark)
	CheckExp(t, pfun, `dog.bark`)
	CheckReturn(t, pfun, "bark")
}

func TestUnaryOp(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*UnaryOpExp)(nil).String(), "nil")
	CheckReturn(t, (*UnaryOpExp)(nil), "nil")

	val := Literal("2.2", types.Float32())
	bin := UnaryOp(val, NegateOp, types.Float32())
	CheckExp(t, bin, `-2.2`)
	CheckReturn(t, bin, `float32`)
}

//============================================================================

// CheckExp checks that the expression's string matches the given string.
func CheckExp(t *common.Tester, e Expression, exp ...string) {
	t.CheckStr(ToString(e), exp...)
}

// CheckReturn checks that the given expression returns
// the given expected type from the ReturnType method.
func CheckReturn(t *common.Tester, e Expression, exp ...string) {
	returnType := e.ReturnType()
	result := types.ToString(returnType)
	expStr := strings.Join(exp, "\n")
	if result != expStr {
		t.Failed("Unexpected return type from expression:", common.NewMap().
			Add("Expression", ToString(e)).
			Add("Expected", expStr).
			Add("Gotten", result))
	}
}
