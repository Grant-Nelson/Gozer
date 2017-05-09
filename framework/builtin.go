package framework

import (
	c "github.com/grant-nelson/Gozer/constructs"
)

// builtinName gets the name for the builtin package.
const builtinName = "builtin"

// BuiltinPrebuild adds the scope information for the builtin package.
// https://golang.org/ref/spec#Built-in_functions
func BuiltinPrebuild(prog *c.ProgramType) {
	if _, exists := prog.Packages[builtinName]; exists {
		return
	}

	// Describe the buildin package.
	pack := c.Package()
	prog.AddPackage(builtinName, pack)

	pack.AddFunction("print").
		AddParam("a", c.Interface()).
		SetEllipse(true)

	pack.AddFunction("println").
		AddParam("a", c.Interface()).
		SetEllipse(true)

	pack.AddInterface("error").
		AddFunction("Error").
		AddReturn("", c.String())

	pack.AddFunction("cap").
		AddParam("a", c.Variant()).
		AddReturn("", c.Int())

	pack.AddFunction("len").
		AddParam("a", c.Variant()).
		AddReturn("", c.Int())

	// TODO: Implement the following:
	//
	// complex, imag, real
	// append, close, copy, delete, make, new, panic, recover
}
