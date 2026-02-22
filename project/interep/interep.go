package interep

import (
	"fmt"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/project"
)

type Config struct {

	// Logger to log verbose messages with. Has no affect if verbose was false.
	Logger *logger.Logger

	// ErrGroup is the collector to handle multiple errors.
	ErrGroup *faults.ErrGroup
}

func Remodel(pkg *project.Package, cfg *Config) (err error) {
	defer faults.Recover(&err)

	fmt.Printf("=================================================\n")
	fmt.Printf("Package:   %q\n", pkg.Ast.Name)
	fmt.Printf("Path:      %q\n", pkg.PkgPath())
	fmt.Printf("AST Files: %d\n", len(pkg.Ast.Syntax))
	fmt.Printf("AST Types: %v\n", pkg.Ast.Types)

	return nil
}
