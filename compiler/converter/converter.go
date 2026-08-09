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
		c.addFault(faults.New(`unexpected AST node type`).
			WithF(`type`, `%T`, n).
			With(`pos`, c.pos(n.Pos())))
		return nil
	}
}

func (c *Converter) FromDecl(d ast.Decl) ir.Stmt {
	switch d := d.(type) {
	case nil, *ast.BadDecl:
		return nil
	case *ast.GenDecl:
		return c.FromGenDecl(d)
	case *ast.FuncDecl:
		return c.FromFuncDecl(d)
	default:
		c.addFault(faults.New(`unexpected AST decl node type`).
			WithF(`type`, `%T`, d).
			With(`pos`, c.pos(d.Pos())))
		return nil
	}
}

func (c *Converter) FromGenDecl(d *ast.GenDecl) ir.Stmt {
	constant := d.Tok == token.CONST
	ss := &ir.StmtListStmt{}
	for _, s := range d.Specs {
		if v, ok := s.(*ast.ValueSpec); ok {
			ss.List = append(ss.List, c.FromValueSpec(v, constant))
		}
	}
	return c.SimplifyStmt(ss)
}

func (c *Converter) FromValueSpec(s *ast.ValueSpec, constant bool) *ir.StmtListStmt {
	ss := &ir.StmtListStmt{}
	for i, n := range s.Names {
		v := &ir.ValueDecl{
			Constant: constant,
			Name:     c.FromIdent(n),
		}
		if len(s.Values) > i {
			v.Value = c.FromExpr(s.Values[i])
		}
		ss.List = append(ss.List, v)
	}
	return ss
}
