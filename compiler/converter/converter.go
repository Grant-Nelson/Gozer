package converter

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"

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
	Errors  *faults.ErrGroup
}

func (c *Converter) pos(p token.Pos) string {
	return c.FileSet.Position(p).String()
}

func (c *Converter) addFault(f *faults.Fault) {
	if c == nil || c.Errors == nil {
		panic(f)
	}
	if err := c.Errors.Add(f); err != nil {
		panic(err)
	}
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
	case *ast.File:
		return c.FromFile(n)
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

// TODO: Finish commenting
//
// This doesn't convert imports.
func (c *Converter) FromPackage(pkg *packages.Package) *ir.Package {
	p := &ir.Package{
		Name:    pkg.Name,
		Path:    pkg.PkgPath,
		FileSet: c.FileSet,
	}
	for _, f := range pkg.Syntax {
		for _, s := range c.ExpandStmt(c.FromFile(f)) {
			switch d := s.(type) {
			case *ir.TypeStmt:
				p.Types = append(p.Types, d)
			case *ir.ValueDecl:
				if d.Constant {
					p.Consts = append(p.Consts, d)
				} else {
					p.Vars = append(p.Vars, d)
				}
			case *ir.Func:
				p.Funcs = append(p.Funcs, d)
			default:
				c.addFault(faults.New(`unexpected AST package-level node type`).
					With(`package`, pkg.PkgPath).
					WithF(`type`, `%T`, d).
					With(`pos`, c.pos(d.Pos())))
			}
		}
	}
	return p
}

func (c *Converter) FromFile(f *ast.File) ir.Stmt {
	ss := &ir.StmtListStmt{}
	for _, d := range f.Decls {
		ss.Add(c.FromDecl(d))
	}
	return c.SimplifyStmt(ss)
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
		switch s := s.(type) {
		case *ast.ImportSpec:
			// Ignore
		case *ast.ValueSpec:
			ss.Add(c.FromValueSpec(s, constant))
		case *ast.TypeSpec:

		default:
			c.addFault(faults.New(`unexpected AST declaration node type`).
				WithF(`type`, `%T`, s).
				With(`pos`, c.pos(s.Pos())))
		}
	}
	return c.SimplifyStmt(ss)
}

func (c *Converter) FromTypeSpec(s *ast.TypeSpec) ir.Stmt {

	// TODO: Finish

	return &ir.TypeStmt{
		TypePos:      s.Pos(),
		TypeAndValue: c.exprTypes(s.Type),
	}
}

func (c *Converter) FromValueSpec(s *ast.ValueSpec, constant bool) ir.Stmt {
	ss := &ir.StmtListStmt{}
	for i, n := range s.Names {
		v := &ir.ValueDecl{
			Constant: constant,
			Name:     c.FromIdent(n),
		}
		if len(s.Values) > i {
			v.Value = c.FromExpr(s.Values[i])
		}
		ss.Add(v)
	}
	return c.SimplifyStmt(ss)
}
