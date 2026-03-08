package analysis

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/astTools"
)

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
