package converter

import (
	"go/ast"
	"go/types"
	"slices"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/unaryOp"
)

func (c *converter) exprTypes(e ast.Expr) *types.TypeAndValue {
	tv, ok := c.Info.Types[e]
	if !ok {
		c.addFault(faults.New(`failed to get expression type`).
			WithF(`expr`, `%T`, e).
			With(`pos`, c.pos(e.Pos())))
		return nil
	}
	return &tv
}

func (c *converter) FromExprSlice(es []ast.Expr) []ir.Expr {
	result := make([]ir.Expr, 0, len(es))
	for _, e := range es {
		result = append(result, c.FromExpr(e))
	}
	return result
}

func (c *converter) FromExpr(e ast.Expr) ir.Expr {
	switch e := e.(type) {
	case nil, *ast.BadExpr:
		return nil
	case *ast.Ident:
		return c.FromIdent(e)
	case *ast.Ellipsis:
		return c.FromEllipsis(e)
	case *ast.BasicLit:
		return c.FromBasicLit(e)
	case *ast.FuncLit:
		return c.FromFuncLit(e)
	case *ast.CompositeLit:
		return c.FromCompositeLit(e)
	case *ast.ParenExpr:
		return c.FromParenExpr(e)
	case *ast.SelectorExpr:
		return c.FromSelectorExpr(e)
	case *ast.IndexExpr:
		return c.FromIndexExpr(e)
	case *ast.IndexListExpr:
		return c.FromIndexListExpr(e)
	case *ast.SliceExpr:
		return c.FromSliceExpr(e)
	case *ast.TypeAssertExpr:
		return c.FromTypeAssertExpr(e)
	case *ast.CallExpr:
		return c.FromCallExpr(e)
	case *ast.StarExpr:
		return c.FromStarExpr(e)
	case *ast.UnaryExpr:
		return c.FromUnaryExpr(e)
	case *ast.BinaryExpr:
		return c.FromBinaryExpr(e)
	case *ast.KeyValueExpr:
		return c.FromKeyValueExpr(e)
	case *ast.ArrayType:
		return c.FromArrayType(e)
	case *ast.StructType:
		return c.FromStructType(e)
	case *ast.FuncType:
		return c.FromFuncType(e)
	case *ast.InterfaceType:
		return c.FromInterfaceType(e)
	case *ast.MapType:
		return c.FromMapType(e)
	case *ast.ChanType:
		return c.FromChanType(e)
	default:
		c.addFault(faults.New(`unexpected AST expression type`).
			WithF(`type`, `%T`, e).
			With(`pos`, c.pos(e.Pos())))
		return nil
	}
}

func (c *converter) FromIdent(e *ast.Ident) ir.Expr {
	if e == nil {
		return nil
	}
	if obj, ok := c.Info.Defs[e]; ok {
		return c.fromIdentDef(e, obj)
	}
	if obj, ok := c.Info.Uses[e]; ok {
		return c.fromIdentUses(e, obj)
	}
	c.addFault(faults.New(`expected a Def or Uses for a naked identifier`).
		With(`id`, e.Name).
		With(`pos`, c.pos(e.Pos())))
	return nil
}

func (c *converter) fromIdentDef(e *ast.Ident, obj types.Object) ir.Expr {
	switch obj := obj.(type) {
	case nil:
		return nil
	case *types.Var:
		return c.fromIdentDefVar(e, obj)
	default:
		c.addFault(faults.New(`unexpected Def object for a naked identifier in an expression`).
			With(`id`, e.Name).
			With(`pos`, c.pos(e.Pos())).
			WithF(`type`, `%T`, obj).
			With(`object`, obj))
		return nil
	}
}

func (c *converter) fromIdentDefVar(e *ast.Ident, obj *types.Var) *ir.VarRef {
	return &ir.VarRef{
		RefPos: e.NamePos,
		VarObj: obj,
	}
}

func (c *converter) fromIdentUses(e *ast.Ident, obj types.Object) ir.Expr {
	switch obj := obj.(type) {
	case *types.PkgName:
		return c.fromIdentUsesPkgName(e, obj)
	case *types.TypeName:
		return c.fromIdentUsesTypeName(e, obj)
	case *types.Const:
		return c.fromIdentUsesConst(e, obj)
	case *types.Var:
		return c.fromIdentUsesVar(e, obj)
	case *types.Func:
		return c.fromIdentUsesFunc(e, obj)
	case *types.Builtin:
		return c.fromIdentUsesBuiltin(e, obj)
	case *types.Nil:
		return c.fromIdentUsesNil(obj)
	default:
		c.addFault(faults.New(`unexpected Use object for a naked identifier in an expression`).
			With(`id`, e.Name).
			With(`pos`, c.pos(e.Pos())).
			WithF(`type`, `%T`, obj).
			With(`object`, obj))
		return nil
	}
}

func (c *converter) fromIdentUsesPkgName(e *ast.Ident, obj *types.PkgName) *ir.ImportRef {
	pkgPath := obj.Pkg().Path()
	var imp *ir.ImportDecl
	if c.Imports != nil {
		imp = c.Imports[pkgPath]
	}
	if imp == nil {
		imp = &ir.ImportDecl{PkgObj: obj}
		if c.Imports != nil {
			c.Imports[pkgPath] = imp
		}
	}
	return &ir.ImportRef{
		RefPos:     e.NamePos,
		ImportDecl: imp,
	}
}

func (c *converter) fromIdentUsesTypeName(e *ast.Ident, obj *types.TypeName) *ir.TypeRef {
	tr := &ir.TypeRef{
		RefPos:  e.NamePos,
		TypeObj: obj,
	}
	if inst, ok := c.Info.Instances[e]; ok {
		tr.Instance = inst.Type
		if inst.TypeArgs != nil {
			tr.TypeArgs = slices.Collect(inst.TypeArgs.Types())
		}
	}
	return tr
}

func (c *converter) fromIdentUsesConst(e *ast.Ident, obj *types.Const) *ir.ConstRef {
	return &ir.ConstRef{
		RefPos:   e.NamePos,
		ConstObj: obj,
	}
}

func (c *converter) fromIdentUsesVar(e *ast.Ident, obj *types.Var) *ir.VarRef {
	return &ir.VarRef{
		RefPos: e.NamePos,
		VarObj: obj,
	}
}

func (c *converter) fromIdentUsesFunc(e *ast.Ident, obj *types.Func) *ir.FuncRef {
	tr := &ir.FuncRef{
		RefPos:  e.NamePos,
		FuncObj: obj,
	}
	if inst, ok := c.Info.Instances[e]; ok {
		tr.TypeArgs = slices.Collect(inst.TypeArgs.Types())
		tr.Instance = inst.Type
	}
	return tr
}

func (c *converter) fromIdentUsesNil(obj *types.Nil) *ir.NilType {
	return &ir.NilType{Obj: obj}
}

func (c *converter) fromIdentUsesBuiltin(e *ast.Ident, obj *types.Builtin) *ir.BuiltinRef {
	return &ir.BuiltinRef{
		RefPos:  e.NamePos,
		Builtin: obj,
	}
}

func (c *converter) FromEllipsis(e *ast.Ellipsis) ir.Expr {
	if e == nil {
		return nil
	}
	return c.FromExpr(e.Elt)
}

func (c *converter) FromBasicLit(e *ast.BasicLit) *ir.BasicLit {
	if e == nil {
		return nil
	}
	return &ir.BasicLit{
		ValuePos:     e.ValuePos,
		TypeAndValue: c.exprTypes(e),
	}
}

func (c *converter) FromFuncLit(e *ast.FuncLit) *ir.FuncLit {
	if e == nil {
		return nil
	}
	t := c.exprTypes(e).Type
	sig, ok := t.(*types.Signature)
	if !ok {
		c.addFault(faults.New(`expected a signature type for a function lit`).
			WithF(`type`, `%T`, t).
			With(`t`, t).
			With(`pos`, c.pos(e.Pos())))
	}
	fn := &ir.Func{
		FuncPos:   e.Type.Func,
		Signature: sig,
	}
	c.setInitialBlock(fn, e.Body, e.Type)
	return &ir.FuncLit{Func: fn}
}

func (c *converter) FromCompositeLit(e *ast.CompositeLit) *ir.TypeLit {
	if e == nil {
		return nil
	}
	return &ir.TypeLit{
		TypePos: e.Pos(),
		TypeRef: c.exprTypes(e).Type,
		Values:  c.FromExprSlice(e.Elts),
	}
}

func (c *converter) FromParenExpr(e *ast.ParenExpr) ir.Expr {
	if e == nil {
		return nil
	}
	// Drop the paren since the order of operators is already kept in the tree shape.
	return c.FromExpr(e.X)
}

func (c *converter) FromSelectorExpr(e *ast.SelectorExpr) *ir.SelectorExpr {
	if e == nil {
		return nil
	}
	return &ir.SelectorExpr{
		X:       c.FromExpr(e.X),
		Sel:     e.Sel.Name,
		SelPos:  e.Sel.NamePos,
		SelType: c.exprTypes(e).Type,
	}
}

func (c *converter) FromIndexExpr(e *ast.IndexExpr) *ir.IndexExpr {
	if e == nil {
		return nil
	}
	// TODO: Should determine if this is for a string, slice, or map,
	//       however, if it is for generics then it should be an index list.
	return &ir.IndexExpr{
		X:          c.FromExpr(e.X),
		LeftPos:    e.Lbrack,
		Index:      c.FromExpr(e.Index),
		ResultType: c.exprTypes(e).Type,
	}
}

func (c *converter) FromIndexListExpr(e *ast.IndexListExpr) *ir.IndexListExpr {
	if e == nil {
		return nil
	}
	// TODO: Is this needed? Or can all the information be gotten
	//       from type and instantiation information on identifiers?
	return &ir.IndexListExpr{
		X:          c.FromExpr(e.X),
		LeftPos:    e.Lbrack,
		Indices:    c.FromExprSlice(e.Indices),
		ResultType: c.exprTypes(e).Type,
	}
}

func (c *converter) FromSliceExpr(e *ast.SliceExpr) *ir.SliceExpr {
	if e == nil {
		return nil
	}
	return &ir.SliceExpr{
		X:          c.FromExpr(e.X),
		LeftPos:    e.Lbrack,
		Low:        c.FromExpr(e.Low),
		High:       c.FromExpr(e.High),
		Max:        c.FromExpr(e.Max),
		Slice3:     e.Slice3,
		ResultType: c.exprTypes(e).Type,
	}
}

func (c *converter) FromTypeAssertExpr(e *ast.TypeAssertExpr) *ir.TypeAssertExpr {
	if e == nil {
		return nil
	}
	return &ir.TypeAssertExpr{
		X:          c.FromExpr(e.X),
		LparenPos:  e.Lparen,
		AssertType: c.FromExpr(e.Type),
		ResultType: c.exprTypes(e).Type,
	}
}

func (c *converter) FromCallExpr(e *ast.CallExpr) *ir.CallExpr {
	if e == nil {
		return nil
	}
	return &ir.CallExpr{
		Fun:          c.FromExpr(e.Fun),
		LeftParenPos: e.Lparen,
		Args:         c.FromExprSlice(e.Args),
		Expanded:     e.Ellipsis.IsValid(),
		ResultType:   c.exprTypes(e).Type,
	}
}

func (c *converter) FromStarExpr(e *ast.StarExpr) ir.Expr {
	if e == nil {
		return nil
	}
	tv := c.exprTypes(e)
	if tv.IsType() {
		return &ir.TypeExpr{
			TypePos: e.Star,
			TypeRef: tv.Type,
		}
	}
	return &ir.UnaryExpr{
		OpPos:      e.Star,
		Op:         unaryOp.Dereference,
		X:          c.FromExpr(e.X),
		ResultType: tv.Type,
	}
}

func (c *converter) FromUnaryExpr(e *ast.UnaryExpr) *ir.UnaryExpr {
	if e == nil {
		return nil
	}
	return &ir.UnaryExpr{
		OpPos:      e.OpPos,
		Op:         c.FromUnaryOp(e.Op, e.OpPos),
		X:          c.FromExpr(e.X),
		ResultType: c.exprTypes(e).Type,
	}
}

func (c *converter) FromBinaryExpr(e *ast.BinaryExpr) *ir.BinaryExpr {
	if e == nil {
		return nil
	}
	return &ir.BinaryExpr{
		X:          c.FromExpr(e.X),
		OpPos:      e.OpPos,
		Op:         c.FromBinaryOp(e.Op, e.OpPos),
		Y:          c.FromExpr(e.Y),
		ResultType: c.exprTypes(e).Type,
	}
}

func (c *converter) FromKeyValueExpr(e *ast.KeyValueExpr) *ir.TypeExpr {
	if e == nil {
		return nil
	}
	return &ir.TypeExpr{
		TypePos: e.Pos(),
		TypeRef: c.exprTypes(e).Type,
	}
}

func (c *converter) FromArrayType(e *ast.ArrayType) *ir.TypeExpr {
	if e == nil {
		return nil
	}
	return &ir.TypeExpr{
		TypePos: e.Pos(),
		TypeRef: c.exprTypes(e).Type,
	}
}

func (c *converter) FromStructType(e *ast.StructType) *ir.TypeExpr {
	if e == nil {
		return nil
	}
	return &ir.TypeExpr{
		TypePos: e.Pos(),
		TypeRef: c.exprTypes(e).Type,
	}
}

func (c *converter) FromFuncType(e *ast.FuncType) *ir.TypeExpr {
	if e == nil {
		return nil
	}
	return &ir.TypeExpr{
		TypePos: e.Pos(),
		TypeRef: c.exprTypes(e).Type,
	}
}

func (c *converter) FromInterfaceType(e *ast.InterfaceType) *ir.TypeExpr {
	if e == nil {
		return nil
	}
	return &ir.TypeExpr{
		TypePos: e.Pos(),
		TypeRef: c.exprTypes(e).Type,
	}
}

func (c *converter) FromMapType(e *ast.MapType) *ir.TypeExpr {
	if e == nil {
		return nil
	}
	return &ir.TypeExpr{
		TypePos: e.Pos(),
		TypeRef: c.exprTypes(e).Type,
	}
}

func (c *converter) FromChanType(e *ast.ChanType) *ir.TypeExpr {
	if e == nil {
		return nil
	}
	return &ir.TypeExpr{
		TypePos: e.Pos(),
		TypeRef: c.exprTypes(e).Type,
	}
}
