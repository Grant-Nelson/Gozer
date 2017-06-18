package framework

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

// ioName gets the name for the io package.
const ioName = "io"

// IOPrebuild adds the scope information for the io package.
// https://golang.org/pkg/io/
func IOPrebuild(prog *types.ProgramType) {
	if prog.Contains(ioName) {
		return
	}

	// Add required imports.
	BuiltinPrebuild(prog)

	// Describe the io package.
	pack := types.Package()
	pack.Name = ioName
	prog.AddPackage(pack)

	// TODO: Add io package
}
