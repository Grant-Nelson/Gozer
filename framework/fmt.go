package framework

import (
	c "github.com/grant-nelson/Gozer/constructs"
)

// fmtName gets the name for the fmt package.
const fmtName = "fmt"

// FmtPrebuild adds the package to the transpiler.
// https://golang.org/pkg/fmt/
func FmtPrebuild(prog *c.ProgramType) {
	if prog.Contains(fmtName) {
		return
	}

	// Add required imports.
	BuiltinPrebuild(prog)
	IoPrebuild(prog)

	// Get the types from other packages needed by this package.
	builtin := prog.Packages[builtinName]
	errType := builtin.Interfaces["error"]

	// Describe the fmt package.
	pack := c.Package()
	prog.AddPackage(fmtName, pack)

	pack.AddFunction("Errorf").
		AddParam("format", c.String()).
		AddParam("a", c.Interface()).
		SetEllipse(true).
		AddResult("", errType)

	pack.AddFunction("Print").
		AddParam("a", c.Interface()).
		SetEllipse(true).
		AddResult("n", c.Int()).
		AddResult("err", errType)

	pack.AddFunction("Print").
		AddParam("format", c.String()).
		AddParam("a", c.Interface()).
		SetEllipse(true).
		AddResult("n", c.Int()).
		AddResult("err", errType)

	pack.AddFunction("Println").
		AddParam("a", c.Interface()).
		SetEllipse(true).
		AddResult("n", c.Int()).
		AddResult("err", errType)

	pack.AddFunction("Sprint").
		AddParam("a", c.Interface()).
		SetEllipse(true).
		AddResult("", c.String())

	pack.AddFunction("Sprintf").
		AddParam("format", c.String()).
		AddParam("a", c.Interface()).
		SetEllipse(true).
		AddResult("", c.String())

	pack.AddFunction("Sprintln").
		AddParam("a", c.Interface()).
		SetEllipse(true).
		AddResult("", c.String())

	// TODO: Add the following:
	//
	// func Fprint(w io.Writer, a ...interface{}) (n int, err error)
	// func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
	// func Fprintln(w io.Writer, a ...interface{}) (n int, err error)
	//
	// func Fscan(r io.Reader, a ...interface{}) (n int, err error)
	// func Fscanf(r io.Reader, format string, a ...interface{}) (n int, err error)
	// func Fscanln(r io.Reader, a ...interface{}) (n int, err error)
	//
	// func Scan(a ...interface{}) (n int, err error)
	// func Scanf(format string, a ...interface{}) (n int, err error)
	// func Scanln(a ...interface{}) (n int, err error)
	//
	// func Sscan(str string, a ...interface{}) (n int, err error)
	// func Sscanf(str string, format string, a ...interface{}) (n int, err error)
	// func Sscanln(str string, a ...interface{}) (n int, err error)
	//
	// And more...
}
