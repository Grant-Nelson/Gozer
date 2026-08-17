package converter

import (
	"fmt"
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
		Sizes:   pkg.TypesSizes,
	}

	c := &converter{
		Source:  pkg,
		Errors:  errGroup,
		Package: p,
		Decls:   map[*ast.Ident]ir.Decl{},
		Refs:    map[*ast.Ident]ir.Ref{},
		Imports: map[string]*ir.ImportDecl{},
	}
	c.PrepareDecls()
	c.PopulateDecls()

	// TODO: Add imports to the ir.Package
	return p, nil
}

const (
	directiveGroup      = `gozer`
	directiveAtomicFunc = `atomic`
)

type converter struct {
	Source  *packages.Package
	Errors  *faults.ErrGroup
	Package *ir.Package
	Decls   map[*ast.Ident]ir.Decl
	Refs    map[*ast.Ident]ir.Ref
	Imports map[string]*ir.ImportDecl
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
		c.prepareDef(id, def)
	}
	for id, def := range c.Source.TypesInfo.Uses {
		c.prepareUse(id, def)
	}
}

func (c *converter) prepareDef(id *ast.Ident, def types.Object) {
	switch def := def.(type) {
	case *types.TypeName:
		c.prepareDefType(id, def)
	case *types.Const:
		c.prepareDefConst(id, def)
	case *types.Var:
		c.prepareDefVar(id, def)
	case *types.Func:
		c.prepareDefFunc(id, def)
	case *types.Label:
		c.prepareDefLabel(id, def)
	default:
		c.addFault(faults.New(`unexpected Def Object`).
			With(`id`, def.Name()).
			WithF(`type`, `%T`, def).
			With(`pos`, c.pos(def.Pos())))
	}
}

func (c *converter) prepareDefType(id *ast.Ident, def *types.TypeName) {
	c.Decls[id] = &ir.TypeDecl{TypeObj: def}
}

func (c *converter) prepareDefConst(id *ast.Ident, def *types.Const) {
	c.Decls[id] = &ir.ConstDecl{ConstObj: def}
}

func (c *converter) prepareDefVar(id *ast.Ident, def *types.Var) {
	c.Decls[id] = &ir.VarDecl{VarObj: def}
}

func (c *converter) prepareDefFunc(id *ast.Ident, def *types.Func) {
	c.Decls[id] = &ir.FuncDecl{FuncObj: def}
}

func (c *converter) prepareDefLabel(id *ast.Ident, def *types.Label) {
	c.Decls[id] = &ir.LabeledStmt{LabelObj: def}
}

func (c *converter) prepareUse(id *ast.Ident, def types.Object) {
	switch def := def.(type) {
	case *types.PkgName:
		c.prepareUseImport(id, def)
	case *types.TypeName:
		c.prepareUseType(id, def)
	case *types.Const:
		c.prepareUseConst(id, def)
	case *types.Var:
		c.prepareUseVar(id, def)
	case *types.Func:
		c.prepareUseFunc(id, def)
	case *types.Label:
		c.prepareUseLabel(id, def)
	case *types.Builtin:
		c.prepareUseBuiltin(id, def)
	case *types.Nil:
		c.Decls[id] = &ir.NilType{Obj: def}
	default:
		c.addFault(faults.New(`unexpected Use Object`).
			With(`id`, def.Name()).
			WithF(`type`, `%T`, def).
			With(`pos`, c.pos(def.Pos())))
	}
}

func (c *converter) prepareUseImport(id *ast.Ident, def *types.PkgName) {
	pkgPath := def.Pkg().Path()
	imp, ok := c.Imports[pkgPath]
	if !ok {
		imp = &ir.ImportDecl{PkgObj: def}
		c.Imports[pkgPath] = imp
	}

	c.Refs[id] = &ir.ImportRef{
		RefPos:     id.NamePos,
		ImportDecl: imp,
	}
}

func (c *converter) prepareUseType(id *ast.Ident, def *types.TypeName) {
	decl, ok := c.Decls[id]
	if !ok {
		// TODO: Determine what to do for externally defined types if that's a problem
		panic(fmt.Errorf(`failed to find decl for %s at %s`, id.Name, id.Pos()))
	}
	td, ok := decl.(*ir.TypeDecl)
	if !ok {
		// TODO: Update
		panic(fmt.Errorf(`failed to cast to type decl for %s at %s: %v`, id.Name, id.Pos(), decl))
	}

	tr := &ir.TypeRef{
		RefPos:   id.Pos(),
		TypeDecl: td,
	}
	if inst, ok := c.Source.TypesInfo.Instances[id]; ok {
		count := inst.TypeArgs.Len()
		typeArgs := make([]types.Type, count)
		for i := range count {
			typeArgs[i] = inst.TypeArgs.At(i)
		}
		tr.TypeArgs = typeArgs
		tr.Instance = inst.Type
	}
	c.Refs[id] = tr
}

func (c *converter) prepareUseConst(id *ast.Ident, def *types.Const) {
	decl, ok := c.Decls[id]
	if !ok {
		// TODO: Determine what to do for externally defined const if that's a problem
		panic(fmt.Errorf(`failed to find decl for %s at %s`, id.Name, id.Pos()))
	}
	cd, ok := decl.(*ir.ConstDecl)
	if !ok {
		// TODO: Update
		panic(fmt.Errorf(`failed to cast to const decl for %s at %s: %v`, id.Name, id.Pos(), decl))
	}

	c.Refs[id] = &ir.ConstRef{
		RefPos:    id.Pos(),
		ConstDecl: cd,
	}
}

func (c *converter) prepareUseVar(id *ast.Ident, def *types.Var) {
	decl, ok := c.Decls[id]
	if !ok {
		// TODO: Determine what to do for externally defined var if that's a problem
		panic(fmt.Errorf(`failed to find decl for %s at %s`, id.Name, id.Pos()))
	}
	vd, ok := decl.(*ir.VarDecl)
	if !ok {
		// TODO: Update
		panic(fmt.Errorf(`failed to cast to var decl for %s at %s: %v`, id.Name, id.Pos(), decl))
	}

	c.Refs[id] = &ir.VarRef{
		RefPos:  id.Pos(),
		VarDecl: vd,
	}
}

func (c *converter) prepareUseFunc(id *ast.Ident, def *types.Func) {
	decl, ok := c.Decls[id]
	if !ok {
		// TODO: Determine what to do for externally defined func if that's a problem
		panic(fmt.Errorf(`failed to find decl for %s at %s`, id.Name, id.Pos()))
	}
	fd, ok := decl.(*ir.FuncDecl)
	if !ok {
		// TODO: Update
		panic(fmt.Errorf(`failed to cast to func decl for %s at %s: %v`, id.Name, id.Pos(), decl))
	}

	fr := &ir.FuncRef{
		RefPos:   id.Pos(),
		FuncDecl: fd,
	}
	if inst, ok := c.Source.TypesInfo.Instances[id]; ok {
		count := inst.TypeArgs.Len()
		typeArgs := make([]types.Type, count)
		for i := range count {
			typeArgs[i] = inst.TypeArgs.At(i)
		}
		fr.TypeArgs = typeArgs
		fr.Instance = inst.Type
	}
	c.Refs[id] = fr
}

func (c *converter) prepareUseLabel(id *ast.Ident, def *types.Label) {
	decl, ok := c.Decls[id]
	if !ok {
		// TODO: Update, the label must be local
		panic(fmt.Errorf(`failed to find decl for %s at %s`, id.Name, id.Pos()))
	}
	ld, ok := decl.(*ir.LabeledStmt)
	if !ok {
		// TODO: Update
		panic(fmt.Errorf(`failed to cast to label decl for %s at %s: %v`, id.Name, id.Pos(), decl))
	}

	c.Refs[id] = &ir.BranchStmt{
		TokPos: id.Pos(),
		Label:  ld,
	}
}

func (c *converter) prepareUseBuiltin(id *ast.Ident, def *types.Builtin) {
	// TODO: Implement
	//c.Refs[id] = &ir.FuncDecl{FuncObj: def}
}

func (c *converter) PopulateDecls() {

	// TODO: FIX

	/*
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
	*/

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
	return &ir.TypeDecl{
		Name:         s.Name.Name,
		NamePos:      s.Name.NamePos,
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
