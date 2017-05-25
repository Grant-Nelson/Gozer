package framework

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

// ioName gets the name for the io package.
const ioName = "io"

// IoPrebuild adds the scope information for the io package.
// https://golang.org/pkg/io/
func IoPrebuild(prog *types.ProgramType) {
	if prog.Contains(ioName) {
		return
	}

	// Add required imports.
	BuiltinPrebuild(prog)

	// Describe the io package.
	pack := types.Package()
	prog.AddPackage(ioName, pack)

	// TODO: Add io package
}
