package analyzer

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/project"
)

type Config struct {
	// Logger to log verbose messages with. Has no affect if verbose was false.
	Logger *logger.Logger

	// ErrGroup is the collector to handle multiple errors.
	ErrGroup *faults.ErrGroup

	// Project is the project to analyze.
	Project *project.Project
}

func Analyze(cfg *Config) error {

	// TODO: Implement

	return nil
}

func ContainsBlockingCall(n ast.Node) bool {
	for it := range astTools.Nodes(n) {
		switch t := it.Node.(type) {
		case *ast.SendStmt, *ast.CallExpr:
			return true
		case *ast.UnaryExpr:
			return t.Op == token.ARROW
		}
	}
	return false
}
