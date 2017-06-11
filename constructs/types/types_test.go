package types

import (
	"strings"
	"testing"

	"github.com/grant-nelson/Gozer/common"
)

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

func TestLookupType(t *testing.T) {
	checkType(t, LookupType("watson"), "nil")
	checkLookupType(t, Bool())
	checkLookupType(t, Byte())
	checkLookupType(t, Complex64())
	checkLookupType(t, Complex128())
	checkLookupType(t, Float32())
	checkLookupType(t, Float64())
	checkLookupType(t, Int())
	checkLookupType(t, Int16())
	checkLookupType(t, Int32())
	checkLookupType(t, Int64())
	checkLookupType(t, Int8())
	checkLookupType(t, Rune())
	checkLookupType(t, String())
	checkLookupType(t, UInt())
	checkLookupType(t, UInt16())
	checkLookupType(t, UInt32())
	checkLookupType(t, UInt64())
	checkLookupType(t, UInt8())
	checkLookupType(t, Variant())
	checkLookupType(t, Void())
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

func TestFunctionTypes(t *testing.T) {
	checkString(t, ((*FunctionType)(nil)).String(), "nil")
	checkString(t, ((*FunctionType)(nil)).FullString(), "nil")
	checkType(t, Function(), "void func()")
	checkType(t, Function().AddParam("name", String()), "void func(string name)")
	checkType(t, Function().AddParam("fmt", String()).
		AddParam("args", String()).SetEllipse(true), "void func(string fmt, string... args)")
	checkType(t, Function().SetReturn(String()), "string func()")

	f1 := Function().SetName("fibonacci").AddParam("name", Float32()).SetReturn(Float32())
	checkType(t, f1, "fibonacci")
	checkString(t, f1.FullString(), "float32 fibonacci(float32 name)")

	f2 := Function().SetName("main")
	f2.Body = "{\n  print(\"Hello World\")\n}"
	checkType(t, f2, "main")
	checkString(t, f2.FullString(),
		`void main() {`,
		`  print("Hello World")`,
		`}`)
}

func TestInterfaceTypes(t *testing.T) {
	checkString(t, ((*InterfaceType)(nil)).String(), "nil")
	checkString(t, ((*InterfaceType)(nil)).FullString(), "nil")
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
	checkString(t, f3.FullString(),
		`void print(string msg)`)
	i1.Name = "logger"
	checkType(t, i1, `logger`)
}

func TestStructureTypes(t *testing.T) {
	checkString(t, ((*StructureType)(nil)).String(), "nil")
	checkString(t, ((*StructureType)(nil)).FullString(), "nil")
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
		`  int age`,
		`  string name`,
		`}`)
	checkFind(t, s1, "temp", "nil")
	checkFind(t, s1, "name", "name")
	checkFind(t, s1, "age", "age")
	s1.AddMember("age", Float32())
	checkType(t, s1,
		`struct{`,
		`  int age`,
		`  string name`,
		`}`)
	s1.Name = "person"
	checkType(t, s1, `person`)
}

func TestClassTypes(t *testing.T) {
	checkString(t, ((*ClassType)(nil)).String(), "nil")
	checkString(t, ((*ClassType)(nil)).FullString(), "nil")
	c1 := Class()
	checkType(t, c1,
		`class{`,
		`  nil`,
		`  interface{}`,
		`}`)
	c1.Data = Int()
	checkType(t, c1,
		`class{`,
		`  int`,
		`  interface{}`,
		`}`)
	c1.Interface.AddFunction("warning").AddParam("text", String()).SetReturn(Int())
	f1 := c1.Interface.AddFunction("count").AddParam("num", Int()).SetReturn(Int())
	checkType(t, c1,
		`class{`,
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
		`class{`,
		`  nil`,
		`  interface{`,
		`    int count(int num)`,
		`    int warning(string text)`,
		`  }`,
		`}`)
	s1 := Structure()
	c1.Data = s1
	checkType(t, c1,
		`class{`,
		`  struct{}`,
		`  interface{`,
		`    int count(int num)`,
		`    int warning(string text)`,
		`  }`,
		`}`)
	s1.AddMember("first", Float64())
	s1.AddMember("last", Float32())
	checkType(t, c1,
		`class{`,
		`  struct{`,
		`    float64 first`,
		`    float32 last`,
		`  }`,
		`  interface{`,
		`    int count(int num)`,
		`    int warning(string text)`,
		`  }`,
		`}`)
	checkFind(t, c1, "first", "first")
	c1.Name = "logger"
	checkType(t, c1,
		`logger`)
	checkString(t, c1.FullString(),
		`logger{`,
		`  struct{`,
		`    float64 first`,
		`    float32 last`,
		`  }`,
		`  interface{`,
		`    int count(int num)`,
		`    int warning(string text)`,
		`  }`,
		`}`)
}

func TestPackageAndProgramTypes(t *testing.T) {
	checkString(t, ((*PackageType)(nil)).String(), "nil")
	p1 := Package()
	checkType(t, p1, `import{}`)
	f1 := p1.AddFunction("boom")
	i1 := p1.AddInterface("pow")
	c1 := p1.AddClass("splat")
	s1 := p1.AddStructure("bawmp")
	d1 := p1.AddDeclaration("pop", String())
	checkType(t, p1,
		`import{`,
		`  string pop`,
		`  void boom()`,
		`  pow{}`,
		`  splat{`,
		`    nil`,
		`    interface{}`,
		`  }`,
		`  bawmp{}`,
		`}`)
	checkFind(t, p1, "temp", "nil")
	checkFind(t, p1, "boom", ToString(f1))
	checkFind(t, p1, "pow", ToString(i1))
	checkFind(t, p1, "splat", ToString(c1))
	checkFind(t, p1, "bawmp", ToString(s1))
	checkFind(t, p1, "pop", ToString(d1))
	p2 := Package()
	p2.AddDeclaration("width", Float32())
	p2.AddDeclaration("height", Float32())
	checkType(t, p2,
		`import{`,
		`  float32 height`,
		`  float32 width`,
		`}`)
	p1.Imports.AddWithShort("other", p2)
	checkType(t, p1,
		`import{`,
		`  import other`,
		`  string pop`,
		`  void boom()`,
		`  pow{}`,
		`  splat{`,
		`    nil`,
		`    interface{}`,
		`  }`,
		`  bawmp{}`,
		`}`)
	checkFind(t, p1, "other", ToString(p2))

	badType, foundBad := ((*PackageType)(nil)).Find("bad")
	if foundBad {
		t.Fatal("Unexpected result from Find on nil package receiver: found returned true.")
	}
	checkType(t, badType, "nil")
	checkType(t, ((*PackageType)(nil)).AddDeclaration("bad", Int()), "nil")
	checkType(t, ((*PackageType)(nil)).AddFunction("bad"), "nil")
	checkType(t, ((*PackageType)(nil)).AddInterface("bad"), "nil")
	checkType(t, ((*PackageType)(nil)).AddClass("bad"), "nil")
	checkType(t, ((*PackageType)(nil)).AddStructure("bad"), "nil")

	checkString(t, ((*ProgramType)(nil)).String(), "nil")
	prog1 := Program()
	checkType(t, prog1, "{}")
	p1.Name = "sounds"
	prog1.AddPackage(p1)
	prog1.AddPackageWithShort("orange", p2)
	checkType(t, prog1,
		`{`,
		`  import orange`,
		`  import sounds`,
		`}`)
	if !prog1.Contains("sounds") {
		t.Fatal("Expected contains to return true for sounds:",
			"\n   Program: ", common.Indent(ToString(prog1), "            "))
	}
	if !prog1.Contains("orange") {
		t.Fatal("Expected contains to return true for short name orange:",
			"\n   Program: ", common.Indent(ToString(prog1), "            "))
	}
	if prog1.Contains("pudding") {
		t.Fatal("Expected contains to return false for pudding:",
			"\n   Program: ", common.Indent(ToString(prog1), "            "))
	}
	p3 := Package()
	p3.Functions.Functions = append(p3.Functions.Functions, nil)
	p3.Interfaces.Interfaces = append(p3.Interfaces.Interfaces, nil)
	p3.Classes.Classes = append(p3.Classes.Classes, nil)
	p3.Structures.Structures = append(p3.Structures.Structures, nil)
	checkType(t, p3,
		`import{`,
		`  nil`,
		`  nil`,
		`  nil`,
		`  nil`,
		`}`)
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

//============================================================================

// Failed indicates a test has failed.
func fail(t *testing.T, text string, m common.Map) {
	result := ""
	if !m.Empty() {
		result = ":\n   " + m.FormatMap("   ")
	}
	t.Fatal(text + result)
}

// checkString checks the the given string matches the given expected lines.
// The lines will be joined with newlines.
func checkString(t *testing.T, result string, exp ...string) {
	expStr := strings.Join(exp, "\n")
	if result != expStr {
		fail(t, "Unexpected construct string", common.NewMap().
			Add("Expected", expStr).
			Add("Gotten", result))
	}
}

// checkType checks that the type's string matches the given string.
func checkType(t *testing.T, ty Type, exp ...string) {
	checkString(t, ToString(ty), exp...)
}

// checkLookupType checks that the LookupType method returns
// the given type when the name of the given type is used in it.
func checkLookupType(t *testing.T, ty Type) {
	str := ToString(ty)
	result := LookupType(str)
	resultStr := ToString(result)
	if str != resultStr {
		fail(t, "Unexpected result from LookupType", common.NewMap().
			Add("Expected", str).
			Add("Gotten", resultStr))
	}
}

// checkGetElement checks that the GetIndexableType method run on the
// given type returns a type matching the given expected string.
func checkGetElement(t *testing.T, t1 Type, exp ...string) {
	t2, found := GetIndexableType(t1)
	expStr := strings.Join(exp, "\n")
	if result := ToString(t2); result != expStr {
		fail(t, "Unexpected construct string", common.NewMap().
			Add("Type", ToString(t1)).
			Add("Found", found).
			Add("Expected", expStr).
			Add("Gotten", result))
	}
}

// checkFind checks that the FindSubtype method run on the given type with
// the given name returns a type matching the given expected string.
func checkFind(t *testing.T, t1 Type, name string, exp ...string) {
	t2, found := FindSubtype(t1, name)
	expStr := strings.Join(exp, "\n")
	if result := ToString(t2); result != expStr {
		fail(t, "Unexpected construct string", common.NewMap().
			Add("Type", ToString(t1)).
			Add("Name", name).
			Add("Found", found).
			Add("Expected", expStr).
			Add("Gotten", result))
	}
}
