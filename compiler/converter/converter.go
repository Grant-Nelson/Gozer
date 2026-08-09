package converter

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
)

const (
	directiveGroup      = `gozer`
	directiveAtomicFunc = `atomic`
)

// TODO: Add `//go:linkname` information for declarations

type Converter struct {
	Info    *types.Info
	FileSet *token.FileSet
}

func (c *Converter) pos(p token.Pos) string {
	return c.FileSet.Position(p).String()
}

func (c *Converter) addFault(f *faults.Fault) {
	// TODO: FINISH
	panic(f)
}

func (c *Converter) FromNodeSlice(ns []ast.Node) []ir.Node {
	result := make([]ir.Node, 0, len(ns))
	for _, n := range ns {
		result = append(result, c.FromNode(n))
	}
	return result
}

func (c *Converter) FromNode(n ast.Node) ir.Node {
	switch n := n.(type) {
	case nil:
		return nil
	case ast.Decl:
		return c.FromDecl(n)
	case ast.Stmt:
		return c.FromStmt(n)
	case ast.Expr:
		return c.FromExpr(n)
	default:
		panic(faults.New(`unexpected AST node type`).
			WithF(`type`, `%T`, n).
			With(`pos`, c.pos(n.Pos())))
	}
}
