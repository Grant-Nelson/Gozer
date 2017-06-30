package types

import (
	"fmt"
	"strings"
	"testing"

	"github.com/grant-nelson/Gozer/common"
)

func TestBaseTypes(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*BaseType)(nil).String(), "nil")
	t.CheckStr((*StringType)(nil).String(), "nil")
	CheckType(t, nil, "nil")
	CheckType(t, Bool(), "bool")
	CheckType(t, Byte(), "byte")
	CheckType(t, Complex64(), "complex64")
	CheckType(t, Complex128(), "complex128")
	CheckType(t, Float32(), "float32")
	CheckType(t, Float64(), "float64")
	CheckType(t, Int(), "int")
	CheckType(t, Int16(), "int16")
	CheckType(t, Int32(), "int32")
	CheckType(t, Int64(), "int64")
	CheckType(t, Int8(), "int8")
	CheckType(t, Rune(), "rune")
	CheckType(t, String(), "string")
	CheckType(t, UInt(), "uint")
	CheckType(t, UInt16(), "uint16")
	CheckType(t, UInt32(), "uint32")
	CheckType(t, UInt64(), "uint64")
	CheckType(t, UInt8(), "uint8")
	CheckType(t, Variant(), "variant")
	CheckType(t, Void(), "void")
	CheckFind(t, String(), "temp", "nil")
}

func TestLookupType(tt *testing.T) {
	t := common.NewTester(tt)
	CheckType(t, LookupType("watson"), "nil")
	CheckLookupType(t, Bool())
	CheckLookupType(t, Byte())
	CheckLookupType(t, Complex64())
	CheckLookupType(t, Complex128())
	CheckLookupType(t, Float32())
	CheckLookupType(t, Float64())
	CheckLookupType(t, Int())
	CheckLookupType(t, Int16())
	CheckLookupType(t, Int32())
	CheckLookupType(t, Int64())
	CheckLookupType(t, Int8())
	CheckLookupType(t, Rune())
	CheckLookupType(t, String())
	CheckLookupType(t, UInt())
	CheckLookupType(t, UInt16())
	CheckLookupType(t, UInt32())
	CheckLookupType(t, UInt64())
	CheckLookupType(t, UInt8())
	CheckLookupType(t, Variant())
	CheckLookupType(t, Void())
}

func TestConstantTypes(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr(((*ConstantType)(nil)).String(), "nil")
	CheckType(t, Constant(nil), "const nil")
	CheckType(t, Constant(Int()), "const int")
	CheckType(t, Constant(String()), "const string")
}

func TestPointerTypes(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr(((*PointerType)(nil)).String(), "nil")
	CheckType(t, Pointer(nil), "*nil")
	CheckType(t, Pointer(Int()), "*int")

	CheckType(t, Pointer(Pointer(Int())), "**int")
	CheckType(t, Pointer(UInt64()), "*uint64")
}

func TestListTypes(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr(((*ListType)(nil)).String(), "nil")
	CheckType(t, List(nil), "[]nil")
	CheckType(t, List(Int()), "[]int")
	CheckType(t, List(Pointer(Int())), "[]*int")
	CheckType(t, List(UInt64()), "[]uint64")
	CheckType(t, List(List(String())), "[][]string")
}

func TestMapTypes(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr(((*MapType)(nil)).String(), "nil")
	CheckType(t, Map(nil, nil), "map[nil]nil")
	CheckType(t, Map(Int(), Int()), "map[int]int")
	CheckType(t, Map(String(), List(Pointer(Int()))), "map[string][]*int")
	CheckType(t, Map(String(), Map(Int(), Int())), "map[string]map[int]int")
}

func TestFunctionTypes(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr(((*FunctionType)(nil)).GetName(), "")
	t.CheckStr(((*FunctionType)(nil)).String(), "nil")
	t.CheckStr(((*FunctionType)(nil)).FullString(), "nil")
	t.CheckStr(((*FunctionType)(nil)).FullBodyString(), "nil")
	t.CheckStr(Function().GetName(), "")
	CheckType(t, Function(), "void func()")
	CheckType(t, Function().AddParam("name", String()), "void func(string name)")
	CheckType(t, Function().AddParam("fmt", String()).
		AddParam("args", String()).SetEllipse(true), "void func(string fmt, string... args)")
	CheckType(t, Function().SetReturn(String()), "string func()")

	f1 := Function().SetName("fibonacci").AddParam("name", Float32()).SetReturn(Float32())
	CheckType(t, f1, "fibonacci")
	t.CheckStr(f1.FullString(), "float32 fibonacci(float32 name)")
	t.CheckStr(f1.GetName(), "fibonacci")

	f2 := Function().SetName("main")
	f2.Body = "{\n  print(\"Hello World\")\n}"
	CheckType(t, f2, "main")
	t.CheckStr(f2.FullString(),
		`void main()`)
	t.CheckStr(f2.FullBodyString(),
		`void main() {`,
		`  print("Hello World")`,
		`}`)
}

func TestInterfaceTypes(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr(((*InterfaceType)(nil)).GetName(), "")
	t.CheckStr(((*InterfaceType)(nil)).String(), "nil")
	t.CheckStr(((*InterfaceType)(nil)).FullString(), "nil")
	i1 := Interface()
	CheckType(t, i1, "interface{}")
	f1 := i1.AddFunction("print").AddParam("msg", String())
	CheckType(t, i1,
		`interface{`,
		`  void print(string msg)`,
		`}`)
	f2 := i1.AddFunction("printf").AddParam("fmt", String()).AddParam("a", Interface()).SetEllipse(true)
	i1.AddFunction("println").AddParam("msg", Interface()).SetEllipse(true)
	CheckType(t, i1,
		`interface{`,
		`  void print(string msg)`,
		`  void printf(string fmt, interface{}... a)`,
		`  void println(interface{}... msg)`,
		`}`)
	CheckFind(t, (*InterfaceType)(nil), "print", "nil")
	CheckFind(t, i1, "temp", "nil")
	CheckFind(t, i1, "print", ToString(f1))
	CheckFind(t, i1, "printf", ToString(f2))
	t.CheckStr(i1.GetName(), "")
	f3 := i1.AddFunction("print")
	t.CheckStr(f3.FullString(),
		`void print(string msg)`)
	i1.Name = "logger"
	CheckType(t, i1, `logger`)
	t.CheckStr(i1.GetName(), "logger")
	CheckType(t, ((*InterfaceType)(nil)).AddFunction("sprint"), `nil`)
}

func TestStructureTypes(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr(((*StructureType)(nil)).GetName(), "")
	t.CheckStr(((*StructureType)(nil)).String(), "nil")
	t.CheckStr(((*StructureType)(nil)).FullString(), "nil")
	s1 := Structure()
	CheckType(t, s1, "struct{}")
	s1.AddMember("name", String())
	CheckType(t, s1,
		`struct{`,
		`  string name`,
		`}`)
	s1.AddMember("age", Int())
	CheckType(t, s1,
		`struct{`,
		`  int age`,
		`  string name`,
		`}`)
	CheckFind(t, (*StructureType)(nil), "age", "nil")
	CheckFind(t, s1, "temp", "nil")
	CheckFind(t, s1, "name", "name")
	CheckFind(t, s1, "age", "age")
	s1.AddMember("age", Float32())
	t.CheckStr(s1.GetName(), "")
	CheckType(t, s1,
		`struct{`,
		`  int age`,
		`  string name`,
		`}`)
	s1.Name = "person"
	CheckType(t, s1, `person`)
	t.CheckStr(s1.GetName(), "person")
	t.CheckStr(s1.FullString(),
		`person{`,
		`  int age`,
		`  string name`,
		`}`)
}

func TestReturnSet(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr(((*ReturnSet)(nil)).GetName(), "")
	t.CheckStr(((*ReturnSet)(nil)).String(), "nil")
	t.CheckStr(((*ReturnSet)(nil)).FullString(), "nil")
	s1 := NewReturnSet()
	CheckType(t, s1, "returns{}")
	s1.AddMember("name", String())
	CheckType(t, s1,
		`returns{`,
		`  string name`,
		`}`)
	s1.AddMember("age", Int())
	CheckType(t, s1,
		`returns{`,
		`  string name`,
		`  int age`,
		`}`)
	s1.Members.Sort()
	CheckType(t, s1,
		`returns{`,
		`  int age`,
		`  string name`,
		`}`)
	CheckFind(t, (*ReturnSet)(nil), "age", "nil")
	CheckFind(t, s1, "temp", "nil")
	CheckFind(t, s1, "name", "name")
	CheckFind(t, s1, "age", "age")
	s1.AddMember("age", Float32())
	t.CheckStr(s1.GetName(), "")
	CheckType(t, s1,
		`returns{`,
		`  int age`,
		`  string name`,
		`}`)
	s1.Name = "person"
	CheckType(t, s1, `person`)
	t.CheckStr(s1.GetName(), "person")
	t.CheckStr(s1.FullString(),
		`person{`,
		`  int age`,
		`  string name`,
		`}`)
}

func TestDeclaration(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr(((*DeclarationType)(nil)).GetName(), "")
	t.CheckStr(((*DeclarationType)(nil)).String(), "nil")
	t.CheckStr(((*DeclarationType)(nil)).FullString(), "nil")
	d1 := Declaration()
	CheckType(t, d1, `nil decl`)
	t.CheckStr(d1.GetName(), "")
	d1.Name = "panda"
	t.CheckStr(d1.GetName(), "panda")
	d1.Data = List(Rune())
	CheckType(t, d1, `panda`)
	d1.Name = ""
	CheckType(t, d1, `[]rune decl`)
	CheckFind(t, (*DeclarationType)(nil), "panda", "nil")
	CheckFind(t, d1, "temp", "nil")
	CheckFind(t, d1, "panda", "nil")
	CheckFind(t, d1, "bear", "nil")
	s1 := Structure()
	s1.AddMember("bear", Int())
	d1.Data = s1
	CheckFind(t, d1, "temp", "nil")
	CheckFind(t, d1, "panda", "nil")
	CheckFind(t, d1, "bear", "bear")
}

func TestClassTypes(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr(((*ClassType)(nil)).GetName(), "")
	t.CheckStr(((*ClassType)(nil)).String(), "nil")
	t.CheckStr(((*ClassType)(nil)).FullString(), "nil")
	CheckType(t, ((*ClassType)(nil)).AddFunction("nope"), "nil")
	c1 := Class()
	CheckType(t, c1,
		`class{}`)
	c1.Data = Int()
	CheckType(t, c1,
		`class{`,
		`  int`,
		`}`)
	t.CheckStr(c1.GetName(), "")
	c1.AddFunction("warning").AddParam("text", String()).SetReturn(Int())
	f1 := c1.AddFunction("count").AddParam("num", Int()).SetReturn(Int())
	CheckType(t, c1,
		`class{`,
		`  int`,
		`  int count(int num)`,
		`  int warning(string text)`,
		`}`)
	CheckFind(t, (*ClassType)(nil), "count", "nil")
	CheckFind(t, c1, "temp", "nil")
	CheckFind(t, c1, "count", ToString(f1))
	c1.Data = nil
	CheckType(t, c1,
		`class{`,
		`  int count(int num)`,
		`  int warning(string text)`,
		`}`)
	s1 := Structure()
	c1.Data = s1
	CheckType(t, c1,
		`class{`,
		`  struct{}`,
		`  int count(int num)`,
		`  int warning(string text)`,
		`}`)
	s1.AddMember("first", Float64())
	s1.AddMember("last", Float32())
	CheckType(t, c1,
		`class{`,
		`  struct{`,
		`    float64 first`,
		`    float32 last`,
		`  }`,
		`  int count(int num)`,
		`  int warning(string text)`,
		`}`)
	CheckFind(t, c1, "first", "first")
	c1.Name = "logger"
	CheckType(t, c1,
		`logger`)
	t.CheckStr(c1.FullString(),
		`logger{`,
		`  struct{`,
		`    float64 first`,
		`    float32 last`,
		`  }`,
		`  int count(int num)`,
		`  int warning(string text)`,
		`}`)
	t.CheckStr(c1.GetName(), "logger")
}

func TestPackageAndProgramTypes(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr(((*PackageType)(nil)).GetName(), "")
	t.CheckStr(((*PackageType)(nil)).String(), "nil")
	t.CheckStr(((*PackageType)(nil)).FullString(), "nil")
	p1 := Package()
	CheckType(t, p1, `import{}`)
	f1 := p1.AddFunction("boom")
	i1 := p1.AddInterface("pow")
	c1 := p1.AddClass("splat")
	s1 := p1.AddStructure("bawmp")
	d1 := p1.AddDeclaration("pop", String())
	r1 := p1.AddReturnSet("zap")
	CheckType(t, p1,
		`import{`,
		`  string pop`,
		`  void boom()`,
		`  pow{}`,
		`  splat{}`,
		`  bawmp{}`,
		`  zap{}`,
		`}`)
	CheckFind(t, p1, "temp", "nil")
	CheckFind(t, p1, "boom", ToString(f1))
	CheckFind(t, p1, "pow", ToString(i1))
	CheckFind(t, p1, "splat", ToString(c1))
	CheckFind(t, p1, "bawmp", ToString(s1))
	CheckFind(t, p1, "pop", ToString(d1))
	CheckFind(t, p1, "zap", ToString(r1))
	p2 := Package()
	p2.AddDeclaration("width", Float32())
	p2.AddDeclaration("height", Float32())
	CheckType(t, p2,
		`import{`,
		`  float32 height`,
		`  float32 width`,
		`}`)
	i2 := p1.AddImport("fmt")
	CheckType(t, i2, "import fmt")
	p1.Imports.AddWithShort("other", p2)
	CheckType(t, p1,
		`import{`,
		`  import fmt`,
		`  import other`,
		`  string pop`,
		`  void boom()`,
		`  pow{}`,
		`  splat{}`,
		`  bawmp{}`,
		`  zap{}`,
		`}`)
	CheckFind(t, p1, "other", ToString(p2))

	badType, foundBad := ((*PackageType)(nil)).Find("bad")
	if foundBad {
		t.Fatal("Unexpected result from Find on nil package receiver: found returned true.")
	}
	CheckType(t, badType, "nil")
	CheckType(t, ((*PackageType)(nil)).AddImport("bad"), "nil")
	CheckType(t, ((*PackageType)(nil)).AddDeclaration("bad", Int()), "nil")
	CheckType(t, ((*PackageType)(nil)).AddFunction("bad"), "nil")
	CheckType(t, ((*PackageType)(nil)).AddInterface("bad"), "nil")
	CheckType(t, ((*PackageType)(nil)).AddClass("bad"), "nil")
	CheckType(t, ((*PackageType)(nil)).AddStructure("bad"), "nil")
	CheckType(t, ((*PackageType)(nil)).AddReturnSet("bad"), "nil")

	t.CheckStr(((*ProgramType)(nil)).String(), "nil")
	prog1 := Program()
	CheckType(t, prog1, "{}")
	t.CheckStr(p1.GetName(), "")
	p1.Name = "sounds"
	t.CheckStr(p1.GetName(), "sounds")
	t.CheckStr(p1.FullString(),
		`import sounds{`,
		`  import fmt`,
		`  import other`,
		`  string pop`,
		`  void boom()`,
		`  pow{}`,
		`  splat{}`,
		`  bawmp{}`,
		`  zap{}`,
		`}`)
	t.CheckStr(p1.FullStringWithShort("what"),
		`import what{`,
		`  import fmt`,
		`  import other`,
		`  string pop`,
		`  void boom()`,
		`  pow{}`,
		`  splat{}`,
		`  bawmp{}`,
		`  zap{}`,
		`}`)

	prog1.AddPackage(p1)
	prog1.AddPackageWithShort("orange", p2)
	p2.Name = "orangeJuice"
	CheckType(t, prog1,
		`{`,
		`  import orange = orangeJuice`,
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
	if ((*ProgramType)(nil)).Contains("sounds") {
		t.Fatal("Expected contains to return false on nil program")
	}
	p3 := Package()
	p3.Functions.Functions = append(p3.Functions.Functions, nil)
	p3.Interfaces.Interfaces = append(p3.Interfaces.Interfaces, nil)
	p3.Classes.Classes = append(p3.Classes.Classes, nil)
	p3.Structures.Structures = append(p3.Structures.Structures, nil)
	CheckType(t, p3,
		`import{`,
		`  nil`,
		`  nil`,
		`  nil`,
		`  nil`,
		`}`)
}

func TestIndexable(tt *testing.T) {
	t := common.NewTester(tt)
	CheckGetElement(t, nil, "nil")
	CheckGetElement(t, (*ListType)(nil), "nil")
	CheckGetElement(t, List(Int()), "int")
	CheckGetElement(t, List(String()), "string")
	CheckGetElement(t, (*MapType)(nil), "nil")
	CheckGetElement(t, Map(Int(), Int()), "int")
	CheckGetElement(t, Map(String(), String()), "string")
	CheckGetElement(t, (*StringType)(nil), "uint8")
	CheckGetElement(t, String(), "uint8")
}

func TestElementTypes(tt *testing.T) {
	t := common.NewTester(tt)
	CheckType(t, (*ListType)(nil).ElementType(), "nil")
	CheckType(t, List(Int()).ElementType(), "int")
	CheckType(t, List(String()).ElementType(), "string")
	CheckType(t, (*MapType)(nil).ElementType(), "nil")
	CheckType(t, Map(Int(), Int()).ElementType(), "int")
	CheckType(t, Map(String(), String()).ElementType(), "string")
	CheckType(t, (*StringType)(nil).ElementType(), "uint8")
	CheckType(t, String().ElementType(), "uint8")
}

func TestClassSet(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*ClassSet)(nil).String(), "nil")
	t.CheckStr((*ClassSet)(nil).FullString(), "nil")
	t1, found1 := (*ClassSet)(nil).AddNew("temp")
	CheckType(t, t1, "nil")
	t.CheckBool(found1, false, "AddNew on a nil class set")
	t.CheckInt((*ClassSet)(nil).Len(), 0, "AddNew on a nil class set")

	set := NewClassSet()
	t2, found2 := set.AddNew("blue")
	CheckType(t, t2, "blue")
	t.CheckBool(found2, true, "AddNew on a class set")
	t.CheckInt(set.Len(), 1, "AddNew on a class set")

	t3, found3 := set.AddNew("blue")
	CheckType(t, t3, "blue")
	t.CheckBool(found3, false, "AddNew for a repeat on a class set")
	t.CheckInt(set.Len(), 1, "AddNew for a repeat on a class set")

	t3.Data = Int()
	set.AddNew("red")
	set.AddNew("green")
	set.Sort()

	t.CheckStr(set.String(),
		`blue`,
		`green`,
		`red`)
	t.CheckStr(set.FullString(),
		`blue{`,
		`  int`,
		`}`,
		`green{}`,
		`red{}`)
}

func TestDeclarationSet(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*DeclarationSet)(nil).String(), "nil")
	t.CheckStr((*DeclarationSet)(nil).FullString(), "nil")
	t1, found1 := (*DeclarationSet)(nil).AddNew("temp", Int())
	CheckType(t, t1, "nil")
	t.CheckBool(found1, false, "AddNew on a nil declaration set")
	t.CheckInt((*DeclarationSet)(nil).Len(), 0, "AddNew on a nil declaration set")

	set := NewDeclarationSet()
	set.AddNew("first", String())
	set.AddNew("second", Int())

	t.CheckStr(set.String(),
		`first`,
		`second`)
	t.CheckStr(set.FullString(),
		`string first`,
		`int second`)
}

func TestFunctionSet(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*FunctionSet)(nil).String(), "nil")
	t.CheckStr((*FunctionSet)(nil).FullString(), "nil")
	t.CheckStr((*FunctionSet)(nil).FullBodyString(), "nil")
	t1, found1 := (*FunctionSet)(nil).AddNew("temp")
	CheckType(t, t1, "nil")
	t.CheckBool(found1, false, "AddNew on a nil function set")
	t.CheckInt((*FunctionSet)(nil).Len(), 0, "AddNew on a nil function set")

	f1 := Function().SetName("fu").AddParam("name", String())
	f2 := Function().SetName("bar").SetReturn(String())
	f2.Body = `{ fu("Hello") }`
	set := NewFunctionSet().Add(f1, f2)

	t.CheckStr(set.String(),
		`bar`,
		`fu`)
	t.CheckStr(set.FullString(),
		`string bar()`,
		`void fu(string name)`)
	t.CheckStr(set.FullBodyString(),
		`string bar() { fu("Hello") }`,
		`void fu(string name)`)
}

func TestInterfaceSet(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*InterfaceSet)(nil).String(), "nil")
	t.CheckStr((*InterfaceSet)(nil).FullString(), "nil")
	t1, found1 := (*InterfaceSet)(nil).AddNew("temp")
	CheckType(t, t1, "nil")
	t.CheckBool(found1, false, "AddNew on a nil interface set")
	t.CheckInt((*InterfaceSet)(nil).Len(), 0, "AddNew on a nil interface set")

	set := NewInterfaceSet()
	i1, found1 := set.AddNew("water")
	t.CheckBool(found1, true, "First created interface in set")
	i1.AddFunction("wamp")

	st := NewReturnSet()
	st.AddMember("val1", Int())
	st.AddMember("val2", Int())
	i1.AddFunction("dull").SetReturn(st)

	i2, found2 := set.AddNew("rocks")
	t.CheckBool(found2, true, "Second created interface in set")
	i2.AddFunction("fu").AddParam("name", String())

	i3, found3 := set.AddNew("rocks")
	t.CheckBool(found3, false, "Finding second created interface in set again")
	i3.AddFunction("bar").SetReturn(String())

	t.CheckStr(set.String(),
		`rocks`,
		`water`)
	t.CheckStr(set.FullString(),
		`rocks{`,
		`  string bar()`,
		`  void fu(string name)`,
		`}`,
		`water{`,
		`  returns{`,
		`    int val1`,
		`    int val2`,
		`  } dull()`,
		`  void wamp()`,
		`}`)
}

func TestPackageSet(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*PackageSet)(nil).String(), "nil")
	t.CheckStr((*PackageSet)(nil).FullString(), "nil")
	t1, found1 := (*PackageSet)(nil).AddNew("temp")
	CheckType(t, t1, "nil")
	t.CheckBool(found1, false, "AddNew on a nil package set")
	t.CheckInt((*PackageSet)(nil).Len(), 0, "AddNew on a nil package set")
	t.CheckStr(fmt.Sprint((*PackageSet)(nil).Shorts()), "[]")
	t.CheckStr(fmt.Sprint((*PackageSet)(nil).Packages()), "[]")

	set := NewPackageSet()
	p1, found1 := set.AddNew("pirate")
	t.CheckBool(found1, true, "first pirate creation")
	p1.AddDeclaration("pegLeg", Int())

	p2, found2 := set.AddNew("zombie")
	t.CheckBool(found2, true, "first zombie creation")
	p2.AddDeclaration("brains", String())

	p3, found3 := set.AddNew("robot")
	t.CheckBool(found3, true, "first robot creation")
	p3.AddDeclaration("bolt", UInt16())

	p4, found4 := set.AddNew("pirate")
	t.CheckBool(found4, false, "repeat pirate creation")
	p4.AddDeclaration("parrot", Rune())

	t.CheckStr(set.String(),
		`import pirate`,
		`import robot`,
		`import zombie`)
	t.CheckStr(set.FullString(),
		`import pirate{`,
		`  rune parrot`,
		`  int pegLeg`,
		`}`,
		`import robot{`,
		`  uint16 bolt`,
		`}`,
		`import zombie{`,
		`  string brains`,
		`}`)

	shortSet1 := set.SetShort("ninja", "zombie")
	t.CheckBool(shortSet1, true, "renaming zombie")

	shortSet2 := set.SetShort("wolf", "moon")
	t.CheckBool(shortSet2, false, "moon 404")

	shortSet3 := set.SetShort("arnold", "robot")
	t.CheckBool(shortSet3, true, "renaming robot")

	set.Sort()
	t.CheckStr(set.String(),
		`import arnold = robot`,
		`import ninja = zombie`,
		`import pirate`)

	t.CheckStr(fmt.Sprint(set.Shorts()), "[arnold ninja ]")
	t.CheckStr(fmt.Sprint(set.Packages()), "[import robot import zombie import pirate]")
}

func TestStructureSet(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*StructureSet)(nil).String(), "nil")
	t.CheckStr((*StructureSet)(nil).FullString(), "nil")
	t1, found1 := (*StructureSet)(nil).AddNew("temp")
	CheckType(t, t1, "nil")
	t.CheckBool(found1, false, "AddNew on a nil structure set")
	t.CheckInt((*StructureSet)(nil).Len(), 0, "AddNew on a nil structure set")

	set := NewStructureSet()
	d1, found1 := set.AddNew("address")
	t.CheckBool(found1, true, "add address")
	d1.AddMember("street", String())
	d1.AddMember("zipCode", Int())
	d1.AddMember("state", String())

	d2, found2 := set.AddNew("customer")
	t.CheckBool(found2, true, "add customer")
	d2.AddMember("first", String())
	d2.AddMember("last", String())
	d2.AddMember("age", Int())

	d3, found3 := set.AddNew("address")
	t.CheckBool(found3, false, "find address again")
	d3.AddMember("county", String())

	t.CheckStr(set.String(),
		`address`,
		`customer`)

	d3.Name = "location"
	set.Sort()
	t.CheckStr(set.String(),
		`customer`,
		`location`)
	t.CheckStr(set.FullString(),
		`customer{`,
		`  int age`,
		`  string first`,
		`  string last`,
		`}`,
		`location{`,
		`  string county`,
		`  string state`,
		`  string street`,
		`  int zipCode`,
		`}`)
}

func TestReturnSetSet(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*ReturnSetSet)(nil).String(), "nil")
	t.CheckStr((*ReturnSetSet)(nil).FullString(), "nil")
	t1, found1 := (*ReturnSetSet)(nil).AddNew("temp")
	CheckType(t, t1, "nil")
	t.CheckBool(found1, false, "AddNew on a nil structure set")
	t.CheckInt((*ReturnSetSet)(nil).Len(), 0, "AddNew on a nil structure set")

	set := NewReturnSetSet()
	d1, found1 := set.AddNew("address")
	t.CheckBool(found1, true, "add address")
	d1.AddMember("street", String())
	d1.AddMember("zipCode", Int())
	d1.AddMember("state", String())

	d2, found2 := set.AddNew("customer")
	t.CheckBool(found2, true, "add customer")
	d2.AddMember("first", String())
	d2.AddMember("last", String())
	d2.AddMember("age", Int())

	d3, found3 := set.AddNew("address")
	t.CheckBool(found3, false, "find address again")
	d3.AddMember("county", String())

	t.CheckStr(set.String(),
		`address`,
		`customer`)

	d3.Name = "location"
	set.Sort()
	t.CheckStr(set.String(),
		`customer`,
		`location`)
	t.CheckStr(set.FullString(),
		`customer{`,
		`  string first`,
		`  string last`,
		`  int age`,
		`}`,
		`location{`,
		`  string street`,
		`  int zipCode`,
		`  string state`,
		`  string county`,
		`}`)
}

//============================================================================

// CheckType checks that the type's string matches the given string.
func CheckType(t *common.Tester, ty Type, exp ...string) {
	t.CheckStr(ToString(ty), exp...)
}

// CheckLookupType checks that the LookupType method returns
// the given type when the name of the given type is used in it.
func CheckLookupType(t *common.Tester, ty Type) {
	str := ToString(ty)
	result := LookupType(str)
	resultStr := ToString(result)
	if str != resultStr {
		t.Failed("Unexpected result from LookupType", common.NewMap().
			Add("Expected", str).
			Add("Gotten", resultStr))
	}
}

// CheckGetElement checks that the GetIndexableType method run on the
// given type returns a type matching the given expected string.
func CheckGetElement(t *common.Tester, t1 Type, exp ...string) {
	t2, found := GetIndexableType(t1)
	expStr := strings.Join(exp, "\n")
	if result := ToString(t2); result != expStr {
		t.Failed("Unexpected construct string", common.NewMap().
			Add("Type", ToString(t1)).
			Add("Found", found).
			Add("Expected", expStr).
			Add("Gotten", result))
	}
}

// CheckFind checks that the FindSubtype method run on the given type with
// the given name returns a type matching the given expected string.
func CheckFind(t *common.Tester, t1 Type, name string, exp ...string) {
	t2, found := FindSubtype(t1, name)
	expStr := strings.Join(exp, "\n")
	if result := ToString(t2); result != expStr {
		t.Failed("Unexpected construct string", common.NewMap().
			Add("Type", ToString(t1)).
			Add("Name", name).
			Add("Found", found).
			Add("Expected", expStr).
			Add("Gotten", result))
	}
}
