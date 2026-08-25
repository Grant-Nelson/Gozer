package converter

import (
	"go/ast"
	"go/token"
	"go/types"
	"maps"
	"slices"

	"golang.org/x/tools/go/packages"

	"github.com/Grant-Nelson/Gozer/avail/assert"
	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
)

func ConvertPackage(pkg *packages.Package, errGroup *faults.ErrGroup) (p *ir.Package, err error) {
	//defer errGroup.Recover(&err) // TODO: Uncomment
	assert.NotNil(pkg)
	assert.NotNil(pkg.TypesInfo)
	assert.NotNil(pkg.Fset)
	c := &converter{
		FileSet: pkg.Fset,
		Info:    pkg.TypesInfo,
		Errors:  errGroup,
	}
	return c.FromPackage(pkg), nil
}

const (
	directiveGroup      = `gozer`
	directiveAtomicFunc = `atomic`
)

type converter struct {
	FileSet *token.FileSet
	Info    *types.Info
	Errors  *faults.ErrGroup
	Imports map[string]*ir.ImportDecl
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

func (c *converter) FromPackage(pkg *packages.Package) *ir.Package {
	p := &ir.Package{
		Name:    pkg.Name,
		Path:    pkg.PkgPath,
		FileSet: pkg.Fset,
		Sizes:   pkg.TypesSizes,
	}

	c.Imports = map[string]*ir.ImportDecl{}
	defer func() { c.Imports = nil }()

	for _, f := range pkg.Syntax {
		for _, stmt := range c.ExpandStmt(c.FromFile(f)) {
			switch stmt := stmt.(type) {
			case *ir.ConstDecl:
				p.Consts = append(p.Consts, stmt)
			case *ir.VarDecl:
				p.Vars = append(p.Vars, stmt)
			case *ir.FuncDecl:
				p.Funcs = append(p.Funcs, stmt)
			case *ir.TypeDecl:
				p.Types = append(p.Types, stmt)
			default:
				c.addFault(faults.New(`unexpected AST node type`).
					WithF(`type`, `%T`, stmt).
					With(`statement`, stmt).
					With(`pos`, c.pos(stmt.Pos())))
			}
		}
	}

	paths := slices.Collect(maps.Keys(c.Imports))
	slices.Sort(paths)
	imports := make([]*ir.ImportDecl, len(paths))
	for i, path := range paths {
		imports[i] = c.Imports[path]
	}
	p.Imports = imports

	return p
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

	return &ir.TypeDecl{TypeObj: typ}
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
			Signature: fnObj.Signature(),
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

	return fn
}

// setInitialBlock creates an initial block, populate it with current statements, add params.
// Named return values are appended after input params so they are in scope
// for the function body just like declared locals.
func (c *converter) setInitialBlock(fn *ir.Func, fnBody *ast.BlockStmt, fnType *ast.FuncType) *ir.Block {
	body := c.ExpandStmt(c.FromBlockStmt(fnBody))
	params := append(c.ExpandParams(fnType), c.ExpandResults(fnType)...)
	return fn.NewBlock(`initial`, body, params)
}
