package converter

import (
	"go/ast"
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/unaryOp"
)

func (c *Converter) exprTypes(e ast.Expr) *types.TypeAndValue {
	tv, ok := c.Info.Types[e]
	if !ok {
		c.addFault(faults.New(`failed to get expression type`).
			WithF(`expr`, `%T`, e).
			With(`pos`, c.pos(e.Pos())))
		return nil
	}
	return &tv
}

func (c *Converter) FromExprSlice(es []ast.Expr) []ir.Expr {
	result := make([]ir.Expr, 0, len(es))
	for _, e := range es {
		result = append(result, c.FromExpr(e))
	}
	return result
}

func (c *Converter) FromExpr(e ast.Expr) ir.Expr {
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

func (c *Converter) FromIdent(e *ast.Ident) *ir.Ident {
	if e == nil {
		return nil
	}
	id := &ir.Ident{
		NamePos: e.NamePos,
		Name:    e.Name,
	}
	if tv, ok := c.Info.Types[e]; ok {
		id.TypeAndValue = &tv
	}
	if in, ok := c.Info.Instances[e]; ok {
		id.Instance = &in
	}
	if ds, ok := c.Info.Defs[e]; ok {
		id.Def = ds
	}
	if us, ok := c.Info.Uses[e]; ok {
		id.Use = us
	}
	return id
}

func (c *Converter) FromEllipsis(e *ast.Ellipsis) ir.Expr {
	if e == nil {
		return nil
	}
	return c.FromExpr(e.Elt)
}

func (c *Converter) FromBasicLit(e *ast.BasicLit) *ir.BasicLit {
	if e == nil {
		return nil
	}
	return &ir.BasicLit{
		ValuePos:     e.ValuePos,
		TypeAndValue: c.exprTypes(e),
	}
}

func (c *Converter) FromCompositeLit(e *ast.CompositeLit) *ir.TypeExpr {
	// TODO: Need to store the initial values
	return c.FromTypeExpr(e)
}

func (c *Converter) FromParenExpr(e *ast.ParenExpr) ir.Expr {
	if e == nil {
		return nil
	}
	// Drop the paren since the order of operators is already kept in the tree shape.
	return c.FromExpr(e.X)
}

func (c *Converter) FromSelectorExpr(e *ast.SelectorExpr) *ir.SelectorExpr {
	if e == nil {
		return nil
	}
	return &ir.SelectorExpr{
		X:       c.FromExpr(e.X),
		Sel:     c.FromIdent(e.Sel),
		SelType: c.exprTypes(e).Type,
	}
}

func (c *Converter) FromIndexExpr(e *ast.IndexExpr) *ir.IndexExpr {
	if e == nil {
		return nil
	}
	// TODO: Should determine if this is for a string, slice, or map
	//       but if it is for generics then it should be an index list.
	return &ir.IndexExpr{
		X:          c.FromExpr(e.X),
		LeftPos:    e.Lbrack,
		Index:      c.FromExpr(e.Index),
		ResultType: c.exprTypes(e).Type,
	}
}

func (c *Converter) FromIndexListExpr(e *ast.IndexListExpr) *ir.IndexListExpr {
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

func (c *Converter) FromSliceExpr(e *ast.SliceExpr) *ir.SliceExpr {
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

func (c *Converter) FromTypeAssertExpr(e *ast.TypeAssertExpr) *ir.TypeAssertExpr {
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

func (c *Converter) FromCallExpr(e *ast.CallExpr) *ir.CallExpr {
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

func (c *Converter) FromStarExpr(e *ast.StarExpr) ir.Expr {
	if e == nil {
		return nil
	}
	tv := c.exprTypes(e)
	if tv.IsType() {
		return &ir.TypeExpr{
			TypePos:      e.Star,
			TypeAndValue: tv,
		}
	}
	return &ir.UnaryExpr{
		OpPos:      e.Star,
		Op:         unaryOp.Dereference,
		X:          c.FromExpr(e.X),
		ResultType: tv.Type,
	}
}

func (c *Converter) FromUnaryExpr(e *ast.UnaryExpr) *ir.UnaryExpr {
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

func (c *Converter) FromBinaryExpr(e *ast.BinaryExpr) *ir.BinaryExpr {
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

func (c *Converter) FromKeyValueExpr(e *ast.KeyValueExpr) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromArrayType(e *ast.ArrayType) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromStructType(e *ast.StructType) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromFuncType(e *ast.FuncType) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromInterfaceType(e *ast.InterfaceType) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromMapType(e *ast.MapType) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromChanType(e *ast.ChanType) *ir.TypeExpr {
	return c.FromTypeExpr(e)
}

func (c *Converter) FromTypeExpr(e ast.Expr) *ir.TypeExpr {
	if e == nil {
		return nil
	}
	return &ir.TypeExpr{
		TypePos:      e.Pos(),
		TypeAndValue: c.exprTypes(e),
	}
}
