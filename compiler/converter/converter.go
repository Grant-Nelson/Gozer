package converter

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/Grant-Nelson/Gozer/avail/assert"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
)

func ConvertPackage(pkg *packages.Package, errGroup *faults.ErrGroup) (p *ir.Package, err error) {
	defer errGroup.Recover(&err)
	assert.NotNil(pkg)
	assert.NotNil(pkg.TypesInfo)
	assert.NotNil(pkg.Fset)

	p = &ir.Package{
		Name:    pkg.Name,
		Path:    pkg.PkgPath,
		FileSet: pkg.Fset,
	}
	c := &converter{
		Source:  pkg,
		Errors:  errGroup,
		Package: p,
	}
	c.PrepareDecls()
	c.PopulateDecls()
	return c.Package, nil
}

const (
	directiveGroup      = `gozer`
	directiveAtomicFunc = `atomic`
)

// TODO: Add `//go:linkname` information for declarations

type converter struct {
	Source  *packages.Package
	Errors  *faults.ErrGroup
	Package *ir.Package
}

func (c *converter) pos(p token.Pos) string {
	return c.Source.Fset.Position(p).String()
}

func (c *converter) addFault(f *faults.Fault) {
	if c == nil || c.Errors == nil {
		panic(f)
	}
	if err := c.Errors.Add(f); err != nil {
		panic(err)
	}
}

func (c *converter) PrepareDecls() {
	for id, def := range c.Source.TypesInfo.Defs {
		c.prepareDeclDefs(id, def)
	}

	// TODO: Implement

}

func (c *converter) prepareDeclDefs(id *ast.Ident, def types.Object) {

	// TODO: Implement

}

func (c *converter) PopulateDecls() {

	// TODO: FIX

	for _, f := range c.Source.Syntax {
		for _, s := range c.ExpandStmt(c.FromFile(f)) {
			switch d := s.(type) {
			case *ir.TypeStmt:
				c.Package.Types = append(c.Package.Types, d)
			case *ir.ValueDecl:
				if d.Constant {
					c.Package.Consts = append(c.Package.Consts, d)
				} else {
					c.Package.Vars = append(c.Package.Vars, d)
				}
			case *ir.Func:
				p.Funcs = append(p.Funcs, d)
			default:
				c.addFault(faults.New(`unexpected AST package-level node type`).
					With(`package`, c.Source.PkgPath).
					WithF(`type`, `%T`, d).
					With(`pos`, c.pos(d.Pos())))
			}
		}
	}

	// TODO: Finish
	//ir.WalkPackage(pkg, )
}

func (c *converter) FromNodeSlice(ns []ast.Node) []ir.Node {
	result := make([]ir.Node, 0, len(ns))
	for _, n := range ns {
		result = append(result, c.FromNode(n))
	}
	return result
}

func (c *converter) FromNode(n ast.Node) ir.Node {
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

func (c *converter) FromFile(f *ast.File) ir.Stmt {
	ss := &ir.StmtListStmt{}
	for _, d := range f.Decls {
		ss.Add(c.FromDecl(d))
	}
	return c.SimplifyStmt(ss)
}

func (c *converter) FromDecl(d ast.Decl) ir.Stmt {
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

func (c *converter) FromGenDecl(d *ast.GenDecl) ir.Stmt {
	constant := d.Tok == token.CONST
	ss := &ir.StmtListStmt{}
	for _, s := range d.Specs {
		switch s := s.(type) {
		case *ast.ImportSpec:
			// Ignore
		case *ast.ValueSpec:
			ss.Add(c.FromValueSpec(s, constant))
		case *ast.TypeSpec:
			ss.Add(c.FromTypeSpec(s))
		default:
			c.addFault(faults.New(`unexpected AST declaration node type`).
				WithF(`type`, `%T`, s).
				With(`pos`, c.pos(s.Pos())))
		}
	}
	return c.SimplifyStmt(ss)
}

func (c *converter) FromTypeSpec(s *ast.TypeSpec) ir.Stmt {
	return &ir.TypeStmt{
		Name:         c.FromIdent(s.Name),
		TypeAndValue: c.exprTypes(s.Type),
	}
}

func (c *converter) FromValueSpec(s *ast.ValueSpec, constant bool) ir.Stmt {
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
