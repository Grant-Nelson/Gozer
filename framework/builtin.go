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

	// Do not define (special implementaions):
	// append, cap, close, copy, delete, len, make, new, panic, recover

	pack.AddFunction("print").
		AddParam("a", c.Interface()).
		SetEllipse(true)

	pack.AddFunction("println").
		AddParam("a", c.Interface()).
		SetEllipse(true)

	pack.AddInterface("error").
		AddFunction("Error").
		AddReturn("", c.String())

	// TODO: Implement the following:
	//
	// complex, imag, real
}
