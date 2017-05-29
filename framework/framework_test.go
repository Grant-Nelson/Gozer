package framework

import (
	"testing"

	"github.com/grant-nelson/Gozer/constructs/types"
)

func TestBuiltin(t *testing.T) {
	prog := types.Program()
	BuiltinPrebuild(prog)
}

func TestFmt(t *testing.T) {
	prog := types.Program()
	FmtPrebuild(prog)
}

func TestIO(t *testing.T) {
	prog := types.Program()
	IOPrebuild(prog)
}
