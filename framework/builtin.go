package framework

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

// builtinName gets the name for the builtin package.
const builtinName = "builtin"

// BuiltinPrebuild adds the scope information for the builtin package.
// https://golang.org/ref/spec#Built-in_functions
func BuiltinPrebuild(prog *types.ProgramType) {
	if prog.Contains(builtinName) {
		return
	}

	// Describe the buildin package.
	pack := types.Package()
	prog.AddPackage(builtinName, pack)

	// https://golang.org/src/builtin/builtin.go?h=make#L254
	pack.AddInterface("error").
		AddFunction("Error").
		SetReturn(types.String())

	// https://golang.org/src/builtin/builtin.go?h=make#L134
	pack.AddFunction("append").
		AddParam("a", types.Variant()).
		AddParam("b", types.Variant()).
		SetEllipse(true).
		SetReturn(types.Variant())

	// https://golang.org/src/builtin/builtin.go?h=make#L164
	pack.AddFunction("cap").
		AddParam("a", types.Variant()).
		SetReturn(types.Int())

	// https://golang.org/src/builtin/builtin.go?h=make#L155
	pack.AddFunction("len").
		AddParam("a", types.Variant()).
		SetReturn(types.Int())

	// https://golang.org/src/builtin/builtin.go?h=make#L243
	pack.AddFunction("print").
		AddParam("a", types.Interface()).
		SetEllipse(true)

	// https://golang.org/src/builtin/builtin.go?h=make#L250
	pack.AddFunction("println").
		AddParam("a", types.Interface()).
		SetEllipse(true)

	// TODO: Implement the following:
	//
	// complex, imag, real
	// close, copy, delete, make, new, panic, recover
}
