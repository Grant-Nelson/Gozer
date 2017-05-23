package expressions

import (
	"strings"
	"testing"

	"github.com/grant-nelson/Gozer/constructs/types"
)

func checkString(t *testing.T, result string, exp ...string) {
	expStr := strings.Join(exp, "\n")
	if result != expStr {
		t.Fatal("Unexpected construct string:",
			"\n   Expected: ", expStr,
			"\n   Gotten:   ", result)
	}
}

func checkExp(t *testing.T, e Expression, exp ...string) {
	checkString(t, ToString(e), exp...)
}

func checkReturn(t *testing.T, e Expression, exp ...string) {
	returnType := e.ReturnType()
	result := types.ToString(returnType)
	expStr := strings.Join(exp, "\n")
	if result != expStr {
		t.Fatal("Unexpected return type from expression:",
			"\n   Expression: ", ToString(e),
			"\n   Expected:   ", expStr,
			"\n   Gotten:     ", result)
	}
}

func TestAssignment(t *testing.T) {
	checkString(t, ((*AssignmentExp)(nil)).String(), "nil")
	checkReturn(t, (*AssignmentExp)(nil), "nil")
	id := Identifier("age", types.Float32())
	val := Literal("42.3", types.Float32())
	assm := Assignment(id, val)
	checkExp(t, assm, `age = 42.3`)
	checkReturn(t, assm, `float32`)
}

func TestBinaryOp(t *testing.T) {
	checkString(t, ((*BinaryOpExp)(nil)).String(), "nil")
	checkReturn(t, (*BinaryOpExp)(nil), "nil")
	id := Identifier("value", types.Float32())
	val1 := Literal("2.2", types.Float32())
	val2 := Literal("-1.7", types.Int32())
	bin := BinaryOp(BinaryOp(val1, id, AddOp, types.Float32()), val2, MultiplyOp, types.Float32())
	checkExp(t, bin, `((2.2 + value) * -1.7)`)
	checkReturn(t, bin, `float32`)
}

func TestCall(t *testing.T) {
	checkString(t, ((*CallExp)(nil)).String(), "nil")
	checkReturn(t, (*CallExp)(nil), "nil")
	log := types.Class()
	print := log.Interface.AddFunction("print").AddParam("msg", types.String()).SetReturn(types.Int())
	id := Identifier("log", log)
	checkExp(t, id, `log`)
	checkReturn(t, id,
		`{`,
		`  interface{`,
		`    int print(string msg)`,
		`  }`,
		`}`)

	pfun := Selector(id, "print", print)
	checkExp(t, pfun, `log.print`)
	checkReturn(t, pfun, "int print(string msg)")

	param := Literal("Hello World", types.String())
	checkExp(t, param, `Hello World`)
	checkReturn(t, param, "string")

	call1 := Call(print, pfun, []Expression{param})
	checkExp(t, call1, `log.print(Hello World)`)
	checkReturn(t, call1, "int")

	call2 := Call(nil, nil, []Expression{param})
	checkExp(t, call2, `nil(Hello World)`)
	checkReturn(t, call2, "nil")
}

func TestCompoundLit(t *testing.T) {
	checkString(t, ((*CompoundLiteralExp)(nil)).String(), "nil")
	checkReturn(t, (*CompoundLiteralExp)(nil), "nil")
	val1 := Literal("1.2", types.Float32())
	val2 := Literal("2.3", types.Float32())
	val3 := Literal("3.4", types.Float32())
	comp := CompoundLiteral([]Expression{val1, val2, val3}, types.List(types.Float32()))
	checkExp(t, comp, `[]float32{1.2, 2.3, 3.4}`)
	checkReturn(t, comp, `[]float32`)
}

func TestDefinition(t *testing.T) {
	checkString(t, ((*DefinitionExp)(nil)).String(), "nil")
	checkReturn(t, (*DefinitionExp)(nil), "nil")
	id := Identifier("dogName", types.String())
	def := Definition(id, Literal(`"Gizmo"`, types.String()))
	checkExp(t, def, `string dogName = "Gizmo"`)
	checkReturn(t, def, `string`)
}

func TestIdentifier(t *testing.T) {
	checkString(t, ((*IdentifierExp)(nil)).String(), "nil")
	checkReturn(t, (*IdentifierExp)(nil), "nil")
	id := Identifier("happy", types.Bool())
	checkExp(t, id, `happy`)
	checkReturn(t, id, `bool`)
}

func TestIndexer(t *testing.T) {
	checkString(t, ((*IndexerExp)(nil)).String(), "nil")
	checkReturn(t, (*IndexerExp)(nil), "nil")
	index := Literal("4", types.Int())

	notIndexable := Indexer(Identifier("notIndexable", types.Int64()), index)
	checkExp(t, notIndexable, `notIndexable[4]`)
	checkReturn(t, notIndexable, `nil`)

	listIndex := Indexer(Identifier("list", types.List(types.Int64())), index)
	checkExp(t, listIndex, `list[4]`)
	checkReturn(t, listIndex, `int64`)

	mapIndex := Indexer(Identifier("map", types.Map(types.Int(), types.Float32())), index)
	checkExp(t, mapIndex, `map[4]`)
	checkReturn(t, mapIndex, `float32`)

	strIndex := Indexer(Identifier("str", types.String()), index)
	checkExp(t, strIndex, `str[4]`)
	checkReturn(t, strIndex, `uint8`)
}

func TestLiteral(t *testing.T) {
	checkString(t, ((*LiteralExp)(nil)).String(), "nil")
	checkReturn(t, (*LiteralExp)(nil), "nil")
	lit := Literal("Hello World", types.String())
	checkExp(t, lit, `Hello World`)
	checkReturn(t, lit, `string`)
}

func TestSelector(t *testing.T) {
	checkString(t, ((*SelectorExp)(nil)).String(), "nil")
	checkReturn(t, (*SelectorExp)(nil), "nil")
	dog := types.Interface()
	bark := dog.AddFunction("bark").AddParam("volume", types.Float64())
	id := Identifier("dog", dog)
	pfun := Selector(id, "bark", bark)
	checkExp(t, pfun, `dog.bark`)
	checkReturn(t, pfun, "void bark(float64 volume)")
}

func TestUnaryOp(t *testing.T) {
	checkString(t, ((*UnaryOpExp)(nil)).String(), "nil")
	checkReturn(t, (*UnaryOpExp)(nil), "nil")

	val := Literal("2.2", types.Float32())
	bin := UnaryOp(val, NegateOp, types.Float32())
	checkExp(t, bin, `-2.2`)
	checkReturn(t, bin, `float32`)
}
