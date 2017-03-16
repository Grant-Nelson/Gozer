package constructs

import (
	"strings"
	"testing"
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

func TestBaseTypes(t *testing.T) {
	checkString(t, ((*BaseType)(nil)).String(), "nil")
	checkConstruct(t, nil, "nil")
	checkConstruct(t, Bool(), "bool")
	checkConstruct(t, Byte(), "byte")
	checkConstruct(t, Float32(), "float32")
	checkConstruct(t, Float64(), "float64")
	checkConstruct(t, Imaginary(), "imaginary")
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
