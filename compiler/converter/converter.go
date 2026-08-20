package converter

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/Grant-Nelson/Gozer/avail/assert"
	"github.com/Grant-Nelson/Gozer/avail/astTools"
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
		Sizes:   pkg.TypesSizes,
		Imports: map[string]*ir.ImportDecl{},
	}

	c := &converter{
		FileSet: pkg.Fset,
		Info:    pkg.TypesInfo,
		Errors:  errGroup,
		Package: p,
	}
	c.FromPackage(pkg)
	return p, nil
}

const (
	directiveGroup      = `gozer`
	directiveAtomicFunc = `atomic`
)

type converter struct {
	FileSet *token.FileSet
	Info    *types.Info
	Errors  *faults.ErrGroup
	Package *ir.Package
}

func (c *converter) pos(p token.Pos) string {
	return c.FileSet.Position(p).String()
}

func (c *converter) addFault(f *faults.Fault) {
	if c == nil || c.Errors == nil {
		panic(f)
	}
	if err := c.Errors.Add(f); err != nil {
		panic(err)
	}
}

func (c *converter) FromPackage(pkg *packages.Package) {
	for _, f := range pkg.Syntax {
		// With package set to `c` the nodes will be set correctly.
		c.FromFile(f)
	}
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
			if constant {
				ss.Add(c.FromConstSpec(s))
			} else {
				ss.Add(c.FromVarSpec(s))
			}
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

func (c *converter) FromTypeSpec(s *ast.TypeSpec) *ir.TypeDecl {
	obj, ok := c.Info.Defs[s.Name]
	if !ok {
		c.addFault(faults.New(`expected a def for a TypeSpec to exist`).
			With(`pos`, c.pos(s.Pos())).
			With(`id`, s.Name.Name))
	}

	typ, ok := obj.(*types.TypeName)
	if !ok {
		c.addFault(faults.New(`expected the object for the TypeSpec to be a TypeName`).
			WithF(`type`, `%T`, obj).
			With(`pos`, c.pos(s.Pos())).
			With(`id`, s.Name.Name))
	}

	td := &ir.TypeDecl{TypeObj: typ}
	if c.Package != nil {
		c.Package.Types = append(c.Package.Types, td)
	}
	return td
}

func (c *converter) FromConstSpec(s *ast.ValueSpec) ir.Stmt {
	ss := &ir.StmtListStmt{}
	for _, n := range s.Names {
		obj, ok := c.Info.Defs[n]
		if !ok {
			c.addFault(faults.New(`expected a def for a constant ValueSpec to exist`).
				With(`pos`, c.pos(s.Pos())).
				With(`id`, n.Name))
		}

		tc, ok := obj.(*types.Const)
		if !ok {
			c.addFault(faults.New(`expected the object for the constant ValueSpec to be a Const`).
				WithF(`type`, `%T`, obj).
				With(`pos`, c.pos(s.Pos())).
				With(`id`, n.Name))
		}

		cd := &ir.ConstDecl{ConstObj: tc}
		if c.Package != nil {
			c.Package.Consts = append(c.Package.Consts, cd)
		}
		ss.Add(cd)
	}
	return c.SimplifyStmt(ss)
}

func (c *converter) FromVarSpec(s *ast.ValueSpec) ir.Stmt {
	ss := &ir.StmtListStmt{}
	for i, n := range s.Names {
		obj, ok := c.Info.Defs[n]
		if !ok {
			c.addFault(faults.New(`expected a def for a variable ValueSpec to exist`).
				With(`pos`, c.pos(s.Pos())).
				With(`id`, n.Name))
		}

		tv, ok := obj.(*types.Var)
		if !ok {
			c.addFault(faults.New(`expected the object for the variable TypeSpec to be a Var`).
				WithF(`type`, `%T`, obj).
				With(`pos`, c.pos(s.Pos())).
				With(`id`, n.Name))
		}

		vd := &ir.VarDecl{VarObj: tv}
		if len(s.Values) > i {
			vd.Value = c.FromExpr(s.Values[i])
		}
		if c.Package != nil {
			c.Package.Vars = append(c.Package.Vars, vd)
		}
		ss.Add(vd)
	}
	return c.SimplifyStmt(ss)
}

func (c *converter) FromFuncDecl(astFunc *ast.FuncDecl) *ir.FuncDecl {
	if astFunc == nil {
		return nil
	}

	obj, ok := c.Info.Defs[astFunc.Name]
	if !ok {
		c.addFault(faults.From(`expected a def object for a function declaration`).
			With(`label`, astFunc.Name).
			With(`pos`, c.pos(astFunc.Pos())))
	}

	fnObj, ok := obj.(*types.Func)
	if !ok {
		c.addFault(faults.From(`unexpected object type for a function declaration`).
			With(`object`, obj).
			WithF(`type`, `%T`, obj).
			With(`pos`, c.pos(astFunc.Pos())))
	}

	fn := &ir.FuncDecl{
		FuncObj: fnObj,
		Func: &ir.Func{
			FuncPos:   astFunc.Type.Func,
			Signature: c.exprTypes(astFunc.Name).Type,
		},
	}
	c.setInitialBlock(fn.Func, astFunc.Body, astFunc.Type)

	if astFunc.Doc != nil {
		dv := astTools.Directives(astFunc.Doc.List, directiveGroup)
		if s, ok := dv[directiveAtomicFunc]; ok {
			// The atomic directive should have no fields.
			// Any fields will be ignored (if asserts are off).
			assert.EmptySlice(s)
			fn.Atomic = true
		}
	}

	if c.Package != nil {
		c.Package.Funcs = append(c.Package.Funcs, fn)
	}
	return fn
}
