package constructs

import (
	"strings"
	"testing"

	"github.com/grant-nelson/Gozer/common"
)

func checkString(t *testing.T, result string, exp ...string) {
	expStr := strings.Join(exp, "\n")
	if result != expStr {
		t.Fatal("Unexpected construct string:",
			"\n   Expected: ", expStr,
			"\n   Gotten:   ", result)
	}
}

func checkConstruct(t *testing.T, c Construct, exp ...string) {
	checkString(t, ToString(c), exp...)
}

func checkFind(t *testing.T, t1 Type, name string, exp ...string) {
	t2, found := FindSubtype(t1, name)
	expStr := strings.Join(exp, "\n")
	if result := ToString(t2); result != expStr {
		t.Fatal("Unexpected construct string:",
			"\n   Type:     ", ToString(t1),
			"\n   Name:     ", name,
			"\n   Found:    ", found,
			"\n   Expected: ", expStr,
			"\n   Gotten:   ", result)
	}
}

func checkReturns(t *testing.T, e Expression, exp ...string) {
	returns := e.ReturnTypes()
	parts := make([]string, len(returns))
	for i, ret := range returns {
		parts[i] = ToString(ret)
	}
	result := strings.Join(parts, "\n")
	expStr := strings.Join(exp, "\n")
	if result != expStr {
		t.Fatal("Unexpected return types from expression:",
			"\n   Type:     ", ToString(e),
			"\n   Expected: ", expStr,
			"\n   Gotten:   ", result)
	}
}

func TestBaseTypes(t *testing.T) {
	checkString(t, ((*BaseType)(nil)).String(), "nil")
	checkConstruct(t, nil, "nil")
	checkConstruct(t, Bool(), "bool")
	checkConstruct(t, Byte(), "byte")
	checkConstruct(t, Complex64(), "complex64")
	checkConstruct(t, Complex128(), "complex128")
	checkConstruct(t, Float32(), "float32")
	checkConstruct(t, Float64(), "float64")
	checkConstruct(t, Int(), "int")
	checkConstruct(t, Int16(), "int16")
	checkConstruct(t, Int32(), "int32")
	checkConstruct(t, Int64(), "int64")
	checkConstruct(t, Int8(), "int8")
	checkConstruct(t, Rune(), "rune")
	checkConstruct(t, String(), "string")
	checkConstruct(t, UInt(), "uint")
	checkConstruct(t, UInt16(), "uint16")
	checkConstruct(t, UInt32(), "uint32")
	checkConstruct(t, UInt64(), "uint64")
	checkConstruct(t, UInt8(), "uint8")
	checkConstruct(t, Variant(), "variant")
	checkConstruct(t, Void(), "void")
	checkFind(t, String(), "temp", "nil")
}

func TestConstantTypes(t *testing.T) {
	checkString(t, ((*ConstantType)(nil)).String(), "nil")
	checkConstruct(t, Constant(nil), "const nil")
	checkConstruct(t, Constant(Int()), "const int")
	checkConstruct(t, Constant(String()), "const string")
}

func TestPointerTypes(t *testing.T) {
	checkString(t, ((*PointerType)(nil)).String(), "nil")
	checkConstruct(t, Pointer(nil), "*nil")
	checkConstruct(t, Pointer(Int()), "*int")

	checkConstruct(t, Pointer(Pointer(Int())), "**int")
	checkConstruct(t, Pointer(UInt64()), "*uint64")
}

func TestListTypes(t *testing.T) {
	checkString(t, ((*ListType)(nil)).String(), "nil")
	checkConstruct(t, List(nil), "[]nil")
	checkConstruct(t, List(Int()), "[]int")
	checkConstruct(t, List(Pointer(Int())), "[]*int")
	checkConstruct(t, List(UInt64()), "[]uint64")
	checkConstruct(t, List(List(String())), "[][]string")
}

func TestMapTypes(t *testing.T) {
	checkString(t, ((*MapType)(nil)).String(), "nil")
	checkConstruct(t, Map(nil, nil), "map[nil]nil")
	checkConstruct(t, Map(Int(), Int()), "map[int]int")
	checkConstruct(t, Map(String(), List(Pointer(Int()))), "map[string][]*int")
	checkConstruct(t, Map(String(), Map(Int(), Int())), "map[string]map[int]int")
}

func TestFunctionTypes(t *testing.T) {
	checkString(t, ((*FunctionType)(nil)).String(), "nil")
	checkConstruct(t, Function(), "func()")
	checkConstruct(t, Function().AddParam("name", String()), "func(name string)")
	checkConstruct(t, Function().AddParam("fmt", String()).
		AddParam("args", String()).SetEllipse(true), "func(fmt string, args ...string)")
	checkConstruct(t, Function().AddReturn("", String()), "func()(string)")
	checkConstruct(t, Function().AddReturn("a", Int()).AddReturn("b", Int()), "func()(a int, b int)")

	f1 := Function().SetName("fibonacci").AddParam("name", Float32()).AddReturn("", Float32())
	checkConstruct(t, f1, "func fibonacci(name float32)(float32)")

	f2 := Function().SetName("main")
	f2.Body = Block(Call(Function().AddParam("name", String()),
		Identifier("print", f1), []Expression{Literal(`"Hello World"`, String())}))
	checkConstruct(t, f2,
		`func main() {`,
		`  print("Hello World")`,
		`}`)
}

func TestInterfaceTypes(t *testing.T) {
	checkString(t, ((*InterfaceType)(nil)).String(), "nil")
	i1 := Interface()
	checkConstruct(t, i1, "interface{}")
	f1 := i1.AddFunction("print").AddParam("msg", String())
	checkConstruct(t, i1,
		`interface{`,
		`  print func(msg string)`,
		`}`)
	f2 := i1.AddFunction("printf").AddParam("fmt", String()).AddParam("a", Interface()).SetEllipse(true)
	i1.AddFunction("println").AddParam("msg", Interface()).SetEllipse(true)
	checkConstruct(t, i1,
		`interface{`,
		`  print func(msg string)`,
		`  printf func(fmt string, a ...interface{})`,
		`  println func(msg ...interface{})`,
		`}`)
	checkFind(t, i1, "temp", "nil")
	checkFind(t, i1, "print", ToString(f1))
	checkFind(t, i1, "printf", ToString(f2))
}

func TestStructureTypes(t *testing.T) {
	checkString(t, ((*StructureType)(nil)).String(), "nil")
	s1 := Structure()
	checkConstruct(t, s1, "struct{}")
	s1.AddMember("name", String())
	checkConstruct(t, s1,
		`struct{`,
		`  name string`,
		`}`)
	s1.AddMember("age", Int())
	checkConstruct(t, s1,
		`struct{`,
		`  age int`,
		`  name string`,
		`}`)
	checkFind(t, s1, "temp", "nil")
	checkFind(t, s1, "name", "string")
	checkFind(t, s1, "age", "int")
}

func TestClassTypes(t *testing.T) {
	checkString(t, ((*ClassType)(nil)).String(), "nil")
	c1 := Class()
	checkConstruct(t, c1, `{}`)
	c1.Data = Int()
	checkConstruct(t, c1,
		`{`,
		`  int`,
		`}`)
	c1.Interface.AddFunction("warning").AddParam("text", String()).AddReturn("", Int())
	f1 := c1.Interface.AddFunction("count").AddParam("num", Int()).AddReturn("", Int())
	checkConstruct(t, c1,
		`{`,
		`  int`,
		`  interface{`,
		`    count func(num int)(int)`,
		`    warning func(text string)(int)`,
		`  }`,
		`}`)
	checkFind(t, c1, "temp", "nil")
	checkFind(t, c1, "count", ToString(f1))
	c1.Data = nil
	checkConstruct(t, c1,
		`{`,
		`  interface{`,
		`    count func(num int)(int)`,
		`    warning func(text string)(int)`,
		`  }`,
		`}`)
	s1 := Structure()
	c1.Data = s1
	checkConstruct(t, c1,
		`{`,
		`  struct{}`,
		`  interface{`,
		`    count func(num int)(int)`,
		`    warning func(text string)(int)`,
		`  }`,
		`}`)
	s1.AddMember("first", Float64())
	s1.AddMember("last", Float32())
	checkConstruct(t, c1,
		`{`,
		`  struct{`,
		`    first float64`,
		`    last float32`,
		`  }`,
		`  interface{`,
		`    count func(num int)(int)`,
		`    warning func(text string)(int)`,
		`  }`,
		`}`)
	checkFind(t, c1, "first", ToString(Float64()))
}

func TestPackageAndProgramTypes(t *testing.T) {
	checkString(t, ((*PackageType)(nil)).String(), "nil")
	p1 := Package()
	checkConstruct(t, p1, `{}`)
	f1 := p1.AddFunction("boom")
	i1 := p1.AddInterface("pow")
	c1 := p1.AddClass("splat")
	p1.AddDeclaration("pop", String())
	checkConstruct(t, p1,
		`{`,
		`  pop string`,
		`  boom func()`,
		`  pow interface{}`,
		`  splat {}`,
		`}`)
	checkFind(t, p1, "temp", "nil")
	checkFind(t, p1, "boom", ToString(f1))
	checkFind(t, p1, "pow", ToString(i1))
	checkFind(t, p1, "splat", ToString(c1))
	checkFind(t, p1, "pop", ToString(String()))
	p2 := Package()
	p2.AddDeclaration("width", Float32())
	p2.AddDeclaration("height", Float32())
	checkConstruct(t, p2,
		`{`,
		`  height float32`,
		`  width float32`,
		`}`)
	p1.Imports["other"] = p2
	checkConstruct(t, p1,
		`{`,
		`  import other`,
		`  pop string`,
		`  boom func()`,
		`  pow interface{}`,
		`  splat {}`,
		`}`)
	checkFind(t, p1, "other", ToString(p2))

	checkString(t, ((*ProgramType)(nil)).String(), "nil")
	prog1 := Program()
	checkConstruct(t, prog1, "{}")
	prog1.AddPackage("sounds", p1)
	prog1.AddPackage("other", p2)
	checkConstruct(t, prog1,
		`{`,
		`  import other`,
		`  import sounds`,
		`}`)
	if !prog1.Contains("sounds") {
		t.Fatal("Expected contains to return true for sounds:",
			"\n   Program: ", common.Indent(ToString(prog1), "            "))
	}
	if prog1.Contains("pudding") {
		t.Fatal("Expected contains to return false for pudding:",
			"\n   Program: ", common.Indent(ToString(prog1), "            "))
	}
}

func TestCallExpressions(t *testing.T) {
	checkString(t, ((*CallExp)(nil)).String(), "nil")
	checkString(t, ((*IdentifierExp)(nil)).String(), "nil")
	checkString(t, ((*LiteralExp)(nil)).String(), "")
	checkString(t, ((*SelectorExp)(nil)).String(), "nil")

	log := Class()
	print := log.Interface.AddFunction("print").AddParam("msg", String()).AddReturn("", Int())
	id := Identifier("log", log)
	checkConstruct(t, id, `log`)
	checkReturns(t, id,
		`{`,
		`  interface{`,
		`    print func(msg string)(int)`,
		`  }`,
		`}`)

	pfun := Selector(id, "print", print)
	checkConstruct(t, pfun, `log.print`)
	checkReturns(t, pfun, "func(msg string)(int)")

	param := Literal("Hello World", String())
	checkConstruct(t, param, `Hello World`)
	checkReturns(t, param, "string")

	call1 := Call(print, pfun, []Expression{param})
	checkConstruct(t, call1, `log.print(Hello World)`)
	checkReturns(t, call1, "int")

	call2 := Call(nil, nil, []Expression{param})
	checkConstruct(t, call2, `nil(Hello World)`)
	checkReturns(t, call2, "")
}
