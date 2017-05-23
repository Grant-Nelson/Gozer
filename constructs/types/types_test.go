package types

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

func checkType(t *testing.T, ty Type, exp ...string) {
	checkString(t, ToString(ty), exp...)
}

func checkGetElement(t *testing.T, t1 Type, exp ...string) {
	t2, found := GetIndexableType(t1)
	expStr := strings.Join(exp, "\n")
	if result := ToString(t2); result != expStr {
		t.Fatal("Unexpected construct string:",
			"\n   Type:     ", ToString(t1),
			"\n   Found:    ", found,
			"\n   Expected: ", expStr,
			"\n   Gotten:   ", result)
	}
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

func TestBaseTypes(t *testing.T) {
	checkString(t, (*BaseType)(nil).String(), "nil")
	checkString(t, (*StringType)(nil).String(), "nil")
	checkType(t, nil, "nil")
	checkType(t, Bool(), "bool")
	checkType(t, Byte(), "byte")
	checkType(t, Complex64(), "complex64")
	checkType(t, Complex128(), "complex128")
	checkType(t, Float32(), "float32")
	checkType(t, Float64(), "float64")
	checkType(t, Int(), "int")
	checkType(t, Int16(), "int16")
	checkType(t, Int32(), "int32")
	checkType(t, Int64(), "int64")
	checkType(t, Int8(), "int8")
	checkType(t, Rune(), "rune")
	checkType(t, String(), "string")
	checkType(t, UInt(), "uint")
	checkType(t, UInt16(), "uint16")
	checkType(t, UInt32(), "uint32")
	checkType(t, UInt64(), "uint64")
	checkType(t, UInt8(), "uint8")
	checkType(t, Variant(), "variant")
	checkType(t, Void(), "void")
	checkFind(t, String(), "temp", "nil")
}

func TestConstantTypes(t *testing.T) {
	checkString(t, ((*ConstantType)(nil)).String(), "nil")
	checkType(t, Constant(nil), "const nil")
	checkType(t, Constant(Int()), "const int")
	checkType(t, Constant(String()), "const string")
}

func TestPointerTypes(t *testing.T) {
	checkString(t, ((*PointerType)(nil)).String(), "nil")
	checkType(t, Pointer(nil), "*nil")
	checkType(t, Pointer(Int()), "*int")

	checkType(t, Pointer(Pointer(Int())), "**int")
	checkType(t, Pointer(UInt64()), "*uint64")
}

func TestListTypes(t *testing.T) {
	checkString(t, ((*ListType)(nil)).String(), "nil")
	checkType(t, List(nil), "[]nil")
	checkType(t, List(Int()), "[]int")
	checkType(t, List(Pointer(Int())), "[]*int")
	checkType(t, List(UInt64()), "[]uint64")
	checkType(t, List(List(String())), "[][]string")
}

func TestMapTypes(t *testing.T) {
	checkString(t, ((*MapType)(nil)).String(), "nil")
	checkType(t, Map(nil, nil), "map[nil]nil")
	checkType(t, Map(Int(), Int()), "map[int]int")
	checkType(t, Map(String(), List(Pointer(Int()))), "map[string][]*int")
	checkType(t, Map(String(), Map(Int(), Int())), "map[string]map[int]int")
}

type testBody struct{}

func (tb *testBody) String() string {
	return "{\n  print(\"Hello World\")\n}"
}

func TestFunctionTypes(t *testing.T) {
	checkString(t, ((*FunctionType)(nil)).String(), "nil")
	checkType(t, Function(), "void func()")
	checkType(t, Function().AddParam("name", String()), "void func(string name)")
	checkType(t, Function().AddParam("fmt", String()).
		AddParam("args", String()).SetEllipse(true), "void func(string fmt, string... args)")
	checkType(t, Function().SetReturn(String()), "string func()")

	f1 := Function().SetName("fibonacci").AddParam("name", Float32()).SetReturn(Float32())
	checkType(t, f1, "float32 fibonacci(float32 name)")

	f2 := Function().SetName("main")
	f2.Body = &testBody{}
	checkType(t, f2,
		`void main() {`,
		`  print("Hello World")`,
		`}`)
}

func TestInterfaceTypes(t *testing.T) {
	checkString(t, ((*InterfaceType)(nil)).String(), "nil")
	i1 := Interface()
	checkType(t, i1, "interface{}")
	f1 := i1.AddFunction("print").AddParam("msg", String())
	checkType(t, i1,
		`interface{`,
		`  void print(string msg)`,
		`}`)
	f2 := i1.AddFunction("printf").AddParam("fmt", String()).AddParam("a", Interface()).SetEllipse(true)
	i1.AddFunction("println").AddParam("msg", Interface()).SetEllipse(true)
	checkType(t, i1,
		`interface{`,
		`  void print(string msg)`,
		`  void printf(string fmt, interface{}... a)`,
		`  void println(interface{}... msg)`,
		`}`)
	checkFind(t, i1, "temp", "nil")
	checkFind(t, i1, "print", ToString(f1))
	checkFind(t, i1, "printf", ToString(f2))
	f3 := i1.AddFunction("print")
	checkType(t, f3, `void print(string msg)`)
	i1.Name = "logger"
	checkType(t, i1, `logger`)
}

func TestStructureTypes(t *testing.T) {
	checkString(t, ((*StructureType)(nil)).String(), "nil")
	s1 := Structure()
	checkType(t, s1, "struct{}")
	s1.AddMember("name", String())
	checkType(t, s1,
		`struct{`,
		`  string name`,
		`}`)
	s1.AddMember("age", Int())
	checkType(t, s1,
		`struct{`,
		`  string name`,
		`  int age`,
		`}`)
	checkFind(t, s1, "temp", "nil")
	checkFind(t, s1, "name", "string")
	checkFind(t, s1, "age", "int")
	s1.AddMember("age", Float32())
	checkType(t, s1,
		`struct{`,
		`  string name`,
		`  float32 age`,
		`}`)
	s1.Name = "person"
	checkType(t, s1, `person`)
}

func TestClassTypes(t *testing.T) {
	checkString(t, ((*ClassType)(nil)).String(), "nil")
	c1 := Class()
	checkType(t, c1, `{}`)
	c1.Data = Int()
	checkType(t, c1,
		`{`,
		`  int`,
		`}`)
	c1.Interface.AddFunction("warning").AddParam("text", String()).SetReturn(Int())
	f1 := c1.Interface.AddFunction("count").AddParam("num", Int()).SetReturn(Int())
	checkType(t, c1,
		`{`,
		`  int`,
		`  interface{`,
		`    int count(int num)`,
		`    int warning(string text)`,
		`  }`,
		`}`)
	checkFind(t, c1, "temp", "nil")
	checkFind(t, c1, "count", ToString(f1))
	c1.Data = nil
	checkType(t, c1,
		`{`,
		`  interface{`,
		`    int count(int num)`,
		`    int warning(string text)`,
		`  }`,
		`}`)
	s1 := Structure()
	c1.Data = s1
	checkType(t, c1,
		`{`,
		`  struct{}`,
		`  interface{`,
		`    int count(int num)`,
		`    int warning(string text)`,
		`  }`,
		`}`)
	s1.AddMember("first", Float64())
	s1.AddMember("last", Float32())
	checkType(t, c1,
		`{`,
		`  struct{`,
		`    float64 first`,
		`    float32 last`,
		`  }`,
		`  interface{`,
		`    int count(int num)`,
		`    int warning(string text)`,
		`  }`,
		`}`)
	checkFind(t, c1, "first", ToString(Float64()))
	c1.Name = "logger"
	checkType(t, c1,
		`logger`)
}

func TestPackageAndProgramTypes(t *testing.T) {
	checkString(t, ((*PackageType)(nil)).String(), "nil")
	p1 := Package()
	checkType(t, p1, `{}`)
	f1 := p1.AddFunction("boom")
	i1 := p1.AddInterface("pow")
	c1 := p1.AddClass("splat")
	p1.AddDeclaration("pop", String())
	checkType(t, p1,
		`{`,
		`  pop string`,
		`  void boom()`,
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
	checkType(t, p2,
		`{`,
		`  height float32`,
		`  width float32`,
		`}`)
	p1.Imports["other"] = p2
	checkType(t, p1,
		`{`,
		`  import other`,
		`  pop string`,
		`  void boom()`,
		`  pow interface{}`,
		`  splat {}`,
		`}`)
	checkFind(t, p1, "other", ToString(p2))

	checkString(t, ((*ProgramType)(nil)).String(), "nil")
	prog1 := Program()
	checkType(t, prog1, "{}")
	prog1.AddPackage("sounds", p1)
	prog1.AddPackage("other", p2)
	checkType(t, prog1,
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

func TestIndexable(t *testing.T) {
	checkGetElement(t, nil, "nil")
	checkGetElement(t, (*ListType)(nil), "nil")
	checkGetElement(t, List(Int()), "int")
	checkGetElement(t, List(String()), "string")
	checkGetElement(t, (*MapType)(nil), "nil")
	checkGetElement(t, Map(Int(), Int()), "int")
	checkGetElement(t, Map(String(), String()), "string")
	checkGetElement(t, (*StringType)(nil), "uint8")
	checkGetElement(t, String(), "uint8")
}

func TestElementTypes(t *testing.T) {
	checkType(t, (*ListType)(nil).ElementType(), "nil")
	checkType(t, List(Int()).ElementType(), "int")
	checkType(t, List(String()).ElementType(), "string")
	checkType(t, (*MapType)(nil).ElementType(), "nil")
	checkType(t, Map(Int(), Int()).ElementType(), "int")
	checkType(t, Map(String(), String()).ElementType(), "string")
	checkType(t, (*StringType)(nil).ElementType(), "uint8")
	checkType(t, String().ElementType(), "uint8")
}
