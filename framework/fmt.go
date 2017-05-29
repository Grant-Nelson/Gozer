package framework

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

// fmtName gets the name for the fmt package.
const fmtName = "fmt"

// FmtPrebuild adds the package to the transpiler.
// https://golang.org/pkg/fmt/
func FmtPrebuild(prog *types.ProgramType) {
	if prog.Contains(fmtName) {
		return
	}

	// Add required imports.
	BuiltinPrebuild(prog)
	IOPrebuild(prog)

	// Get the types from other packages needed by this package.
	builtin := prog.Packages[builtinName]
	errType := builtin.Interfaces["error"]

	// Describe the fmt package.
	pack := types.Package()
	prog.AddPackage(fmtName, pack)

	// return structure for methods similar to Print and Println
	lenRet := pack.AddStructure("lengthReturn").
		AddMember("n", types.Int()).
		AddMember("err", errType)

	// https://golang.org/pkg/fmt/#Errorf
	pack.AddFunction("Errorf").
		AddParam("format", types.String()).
		AddParam("a", types.Interface()).
		SetEllipse(true).
		SetReturn(errType)

	// https://golang.org/pkg/fmt/#Print
	pack.AddFunction("Print").
		AddParam("a", types.Interface()).
		SetEllipse(true).
		SetReturn(lenRet)

	// https://golang.org/pkg/fmt/#Printf
	pack.AddFunction("Printf").
		AddParam("format", types.String()).
		AddParam("a", types.Interface()).
		SetEllipse(true).
		SetReturn(lenRet)

	// https://golang.org/pkg/fmt/#Println
	pack.AddFunction("Println").
		AddParam("a", types.Interface()).
		SetEllipse(true).
		SetReturn(lenRet)

	// https://golang.org/pkg/fmt/#Sprint
	pack.AddFunction("Sprint").
		AddParam("a", types.Interface()).
		SetEllipse(true).
		SetReturn(types.String())

	// https://golang.org/pkg/fmt/#Sprintf
	pack.AddFunction("Sprintf").
		AddParam("format", types.String()).
		AddParam("a", types.Interface()).
		SetEllipse(true).
		SetReturn(types.String())

	// https://golang.org/pkg/fmt/#Sprintln
	pack.AddFunction("Sprintln").
		AddParam("a", types.Interface()).
		SetEllipse(true).
		SetReturn(types.String())

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
