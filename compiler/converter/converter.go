package converter

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/branchKind"
)

type Converter struct {
	Info    *types.Info
	FileSet *token.FileSet
}

func (c *Converter) FromStmtSlice(ss []ast.Stmt) []ir.Stmt {
	result := make([]ir.Stmt, 0, len(ss))
	for _, s := range ss {
		result = append(result, c.ExpandStmt(c.FromStmt(s))...)
	}
	return result
}

func (c *Converter) ExpandStmtSlice(ss []ir.Stmt) []ir.Stmt {
	st := make([]ir.Stmt, 0, len(ss))
	for _, s := range ss {
		st = append(st, c.ExpandStmt(s)...)
	}
	return st
}

func (c *Converter) ExpandStmt(s ir.Stmt) []ir.Stmt {
	if s == nil {
		return []ir.Stmt{}
	}
	if b, ok := s.(*ir.StmtListStmt); ok {
		return c.ExpandStmtSlice(b.List)
	}
	return []ir.Stmt{s}
}

func (c *Converter) FromStmt(s ast.Stmt) ir.Stmt {
	switch s := s.(type) {
	case nil, *ast.BadStmt, *ast.EmptyStmt:
		return nil
	case *ast.DeclStmt:
		return c.FromDeclStmt(s)
	case *ast.LabeledStmt:
		return c.FromLabeledStmt(s)
	case *ast.ExprStmt:
		return c.FromExprStmt(s)
	case *ast.SendStmt:
		return c.FromSendStmt(s)
	case *ast.IncDecStmt:
		return c.FromIncDecStmt(s)
	case *ast.AssignStmt:
		return c.FromAssignStmt(s)
	case *ast.GoStmt:
		return c.FromGoStmt(s)
	case *ast.DeferStmt:
		return c.FromDeferStmt(s)
	case *ast.ReturnStmt:
		return c.FromReturnStmt(s)
	case *ast.BranchStmt:
		return c.FromBranchStmt(s)
	case *ast.BlockStmt:
		return c.FromBlockStmt(s)
	case *ast.IfStmt:
		return c.FromIfStmt(s)
	case *ast.CaseClause:
		return c.FromCaseClause(s)
	case *ast.SwitchStmt:
		return c.FromSwitchStmt(s)
	case *ast.TypeSwitchStmt:
		return c.FromTypeSwitchStmt(s)
	case *ast.CommClause:
		return c.FromCommClause(s)
	case *ast.SelectStmt:
		return c.FromSelectStmt(s)
	case *ast.ForStmt:
		return c.FromForStmt(s)
	case *ast.RangeStmt:
		return c.FromRangeStmt(s)
	default:
		panic(faults.New(`unexpected AST statement type`).
			WithF(`type`, `%T`, s))
	}
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
	// TODO: FIX by making func lit have an expression reference to the func
	// case *ast.FuncLit:
	//	return c.FromFuncLit(e)
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
		panic(faults.New(`unexpected AST expression type`).
			WithF(`type`, `%T`, e))
	}
}

func (c *Converter) FromAssignStmt(s *ast.AssignStmt) *ir.AssignStmt {
	if s == nil {
		return nil
	}
	return &ir.AssignStmt{
		TokPos: s.TokPos,
		Lhs:    c.FromExprSlice(s.Lhs),
		Define: s.Tok == token.DEFINE,
		Rhs:    c.FromExprSlice(s.Rhs),
	}
}

func (c *Converter) FromBranchToken(t token.Token) branchKind.BranchKind {
	switch t {
	case token.BREAK:
		return branchKind.Break
	case token.CONTINUE:
		return branchKind.Continue
	case token.GOTO:
		return branchKind.Goto
	case token.FALLTHROUGH:
		return branchKind.Fallthrough
	default:
		panic(faults.New(`unexpected token for a branch kind`).
			With(`token`, t.String()))
	}
}

func (c *Converter) FromBranchStmt(s *ast.BranchStmt) *ir.BranchStmt {
	if s == nil {
		return nil
	}
	return &ir.BranchStmt{
		TokPos: s.TokPos,
		Kind:   c.FromBranchToken(s.Tok),
		Label:  s.Label,
	}
}

func (c *Converter) FromCaseClause(s *ast.CaseClause) *ir.CaseClause {
	if s == nil {
		return nil
	}
	return &ir.CaseClause{
		Case: s.Case,
		List: c.FromExprSlice(s.List),
		Body: c.FromStmtSlice(s.Body),
	}
}

func (c *Converter) FromCaseClauseSlice(s *ast.BlockStmt) []*ir.CaseClause {
	if s == nil {
		return []*ir.CaseClause{}
	}
	ccs := make([]*ir.CaseClause, 0, len(s.List))
	for i, ct := range s.List {
		if ct == nil {
			continue
		}
		cc, ok := ct.(*ast.CaseClause)
		if !ok {
			panic(faults.New(`expected case clause`).
				WithF(`type`, `%T`, ct).
				With(`index`, i))
		}
		ccs = append(ccs, c.FromCaseClause(cc))
	}
	return ccs
}

func (c *Converter) FromCommClause(s *ast.CommClause) *ir.CommClause {
	if s == nil {
		return nil
	}
	return &ir.CommClause{
		Case: s.Case,
		Comm: c.FromStmt(s.Comm),
		Body: c.FromStmtSlice(s.Body),
	}
}

func (c *Converter) FromCommClauseSlice(s *ast.BlockStmt) []*ir.CommClause {
	if s == nil {
		return []*ir.CommClause{}
	}
	ccs := make([]*ir.CommClause, 0, len(s.List))
	for i, ct := range s.List {
		if ct == nil {
			continue
		}
		cc, ok := ct.(*ast.CommClause)
		if !ok {
			panic(faults.New(`expected comm clause`).
				WithF(`type`, `%T`, ct).
				With(`index`, i))
		}
		ccs = append(ccs, c.FromCommClause(cc))
	}
	return ccs
}

func (c *Converter) FromDeclStmt(s *ast.DeclStmt) *ir.DeclStmt {
	if s == nil {
		return nil
	}
	return &ir.DeclStmt{
		Ast:  s,
		Decl: s.Decl,
	}
}

func (c *Converter) FromDeferStmt(s *ast.DeferStmt) *ir.DeferStmt {
	if s == nil {
		return nil
	}
	return &ir.DeferStmt{
		Defer: s.Defer,
		Call:  s.Call,
	}
}

func (c *Converter) FromExprStmt(s *ast.ExprStmt) *ir.ExprStmt {
	if s == nil {
		return nil
	}
	return &ir.ExprStmt{
		X: c.FromExpr(s.X),
	}
}

func (c *Converter) FromIncDecStmt(s *ast.IncDecStmt) *ir.ExprStmt {
	if s == nil {
		return nil
	}
	exp := &ast.UnaryExpr{OpPos: s.Pos(), Op: s.Tok, X: s.X}
	c.Info.Types[exp] = c.Info.Types[s.X]
	return &ir.ExprStmt{
		X: c.FromExpr(exp),
	}
}

func (c *Converter) FromForStmt(s *ast.ForStmt) *ir.ForStmt {
	if s == nil {
		return nil
	}
	return &ir.ForStmt{
		Ast:  s,
		Init: c.FromStmt(s.Init),
		Cond: c.FromExpr(s.Cond),
		Post: c.FromStmt(s.Post),
		Body: c.ExpandStmt(c.FromBlockStmt(s.Body)),
	}
}

func (c *Converter) FromGoStmt(s *ast.GoStmt) *ir.GoStmt {
	if s == nil {
		return nil
	}
	return &ir.GoStmt{
		Ast:  s,
		Call: c.FromCallExpr(s.Call),
	}
}

func (c *Converter) FromIfStmt(s *ast.IfStmt) *ir.IfStmt {
	if s == nil {
		return nil
	}
	return &ir.IfStmt{
		Ast:  s,
		Init: c.FromStmt(s.Init),
		Cond: c.FromExpr(s.Cond),
		Body: c.ExpandStmt(c.FromBlockStmt(s.Body)),
		Else: c.ExpandStmt(c.FromStmt(s.Else)),
	}
}

func (c *Converter) FromLabeledStmt(s *ast.LabeledStmt) *ir.LabeledStmt {
	if s == nil {
		return nil
	}
	return &ir.LabeledStmt{
		Ast:   s,
		Label: s.Label,
		Stmt:  c.FromStmt(s.Stmt),
	}
}

func (c *Converter) ExpandParams(ft *ast.FuncType) []*ir.Param {
	if ft == nil {
		return nil
	}
	return c.ExpandFieldList(ft.Params)
}

func (c *Converter) ExpandResults(ft *ast.FuncType) []*ir.Param {
	if ft == nil {
		return nil
	}
	return c.ExpandFieldList(ft.Results)
}

func (c *Converter) ExpandFieldList(fl *ast.FieldList) []*ir.Param {
	if fl == nil || len(fl.List) <= 0 {
		return nil
	}
	params := make([]*ir.Param, 0, len(fl.List))
	for _, field := range fl.List {
		// If field.Names is nil then this is an unnamed field
		// that should be skipped over when defining block params.
		for _, name := range field.Names {
			if name == nil {
				continue
			}
			param := &ir.Param{
				Name: name,
				Expr: field.Type,
			}
			params = append(params, param)
		}
	}
	return params
}

func (c *Converter) FromRangeStmt(s *ast.RangeStmt) *ir.RangeStmt {
	if s == nil {
		return nil
	}
	return &ir.RangeStmt{
		Ast:   s,
		Key:   c.FromExpr(s.Key),
		Value: c.FromExpr(s.Value),
		X:     c.FromExpr(s.X),
		Body:  c.ExpandStmt(c.FromBlockStmt(s.Body)),
	}
}

func (c *Converter) FromReturnStmt(s *ast.ReturnStmt) *ir.ReturnStmt {
	if s == nil {
		return nil
	}
	return &ir.ReturnStmt{
		Ast:     s,
		Results: c.FromExprSlice(s.Results),
	}
}

func (c *Converter) FromSelectStmt(s *ast.SelectStmt) *ir.SelectStmt {
	if s == nil {
		return nil
	}
	return &ir.SelectStmt{
		Ast:  s,
		Body: c.FromCommClauseSlice(s.Body),
	}
}

func (c *Converter) FromSendStmt(s *ast.SendStmt) *ir.SendStmt {
	if s == nil {
		return nil
	}
	return &ir.SendStmt{
		Ast:   s,
		Chan:  c.FromExpr(s.Chan),
		Value: c.FromExpr(s.Value),
	}
}

func (c *Converter) FromBlockStmt(s *ast.BlockStmt) *ir.StmtListStmt {
	if s == nil {
		return nil
	}
	return &ir.StmtListStmt{
		Ast:  s,
		List: c.FromStmtSlice(s.List),
	}
}

func (c *Converter) FromSwitchStmt(s *ast.SwitchStmt) *ir.SwitchStmt {
	if s == nil {
		return nil
	}
	return &ir.SwitchStmt{
		Ast:  s,
		Init: c.FromStmt(s.Init),
		Tag:  c.FromExpr(s.Tag),
		Body: c.FromCaseClauseSlice(s.Body),
	}
}

func (c *Converter) FromTypeSwitchStmt(s *ast.TypeSwitchStmt) *ir.TypeSwitchStmt {
	if s == nil {
		return nil
	}
	return &ir.TypeSwitchStmt{
		Ast:    s,
		Init:   c.FromStmt(s.Init),
		Assign: c.FromStmt(s.Assign),
		Body:   c.FromCaseClauseSlice(s.Body),
	}
}

func (c *Converter) FromIdent(e *ast.Ident) *ir.Ident {
	if e == nil {
		return nil
	}
	return &ir.Ident{
		NamePos: e.NamePos,
		Name:    e.Name,
	}
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
		ValuePos: e.ValuePos,
		Value:    constant.MakeFromLiteral(e.Value, e.Kind, 0),
	}
}

func (c *Converter) FromFuncDecl(astFunc *ast.FuncDecl) *ir.Func {
	if astFunc == nil {
		return nil
	}
	fn := &ir.Func{
		Ast:     astFunc,
		FuncPos: astFunc.Type.Func,
		Name:    astFunc.Name.Name,
	}

	// Create initial block, populate it with current statements, add params.
	// Named return values are appended after input params so they are in scope
	// for the function body just like declared locals.
	body := c.ExpandStmt(c.FromBlockStmt(astFunc.Body))
	params := append(c.ExpandParams(astFunc.Type), c.ExpandResults(astFunc.Type)...)
	fn.NewBlock(`initial`, body, params)
	return fn
}

func (c *Converter) FromFuncLit(astFunc *ast.FuncLit) *ir.Func {
	if astFunc == nil {
		return nil
	}
	pos := astFunc.Type.Func
	fn := &ir.Func{
		FuncPos: pos,
		Name:    fmt.Sprintf(`unnamed-lit-%d`, pos),
	}

	// Create initial block, populate it with current statements, add params.
	// Named return values are appended after input params so they are in scope
	// for the function body just like declared locals.
	body := c.ExpandStmt(c.FromBlockStmt(astFunc.Body))
	params := append(c.ExpandParams(astFunc.Type), c.ExpandResults(astFunc.Type)...)
	fn.NewBlock(`initial`, body, params)
	return fn
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
		X:   c.FromExpr(e.X),
		Sel: c.FromIdent(e.Sel),
	}
}

func (c *Converter) FromIndexExpr(e *ast.IndexExpr) *ir.IndexExpr {
	if e == nil {
		return nil
	}
	// TODO: Should determine if this is for a string, slice, or map
	//       but if it is for generics then it should be an index list.
	return &ir.IndexExpr{
		X:       c.FromExpr(e.X),
		LeftPos: e.Lbrack,
		Index:   c.FromExpr(e.Index),
	}
}

func (c *Converter) FromIndexListExpr(e *ast.IndexListExpr) *ir.IndexListExpr {
	if e == nil {
		return nil
	}
	return &ir.IndexListExpr{
		X:       c.FromExpr(e.X),
		LeftPos: e.Lbrack,
		Indices: c.FromExprSlice(e.Indices),
	}
}

func (c *Converter) FromSliceExpr(e *ast.SliceExpr) *ir.SliceExpr {
	if e == nil {
		return nil
	}
	return &ir.SliceExpr{
		X:       c.FromExpr(e.X),
		LeftPos: e.Lbrack,
		Low:     c.FromExpr(e.Low),
		High:    c.FromExpr(e.High),
		Max:     c.FromExpr(e.Max),
		Slice3:  e.Slice3,
	}
}

func (c *Converter) FromTypeAssertExpr(e *ast.TypeAssertExpr) *ir.TypeAssertExpr {
	if e == nil {
		return nil
	}
	return &ir.TypeAssertExpr{
		X:         c.FromExpr(e.X),
		LparenPos: e.Lparen,
		Type:      c.FromExpr(e.Type),
	}
}

func (c *Converter) FromCallExpr(e *ast.CallExpr) *ir.CallExpr {
	if e == nil {
		return nil
	}
	return &ir.CallExpr{
		Fun:       c.FromExpr(e.Fun),
		LparenPos: e.Lparen,
		Args:      c.FromExprSlice(e.Args),
		Variadic:  e.Ellipsis.IsValid(),
	}
}

func (c *Converter) FromStarExpr(e *ast.StarExpr) ir.Expr {
	if e == nil {
		return nil
	}
	if tv, ok := c.Info.Types[e]; ok && tv.IsType() {
		return &ir.TypeExpr{
			TypePos:      e.Star,
			TypeAndValue: tv,
		}
	}
	return &ir.UnaryExpr{
		OpPos: e.Star,
		Op:    token.MUL,
		X:     c.FromExpr(e.X),
	}
}

func (c *Converter) FromUnaryExpr(e *ast.UnaryExpr) *ir.UnaryExpr {
	if e == nil {
		return nil
	}
	return &ir.UnaryExpr{
		OpPos: e.OpPos,
		Op:    e.Op,
		X:     c.FromExpr(e.X),
	}
}

func (c *Converter) FromBinaryExpr(e *ast.BinaryExpr) *ir.BinaryExpr {
	if e == nil {
		return nil
	}
	return &ir.BinaryExpr{
		X:     c.FromExpr(e.X),
		OpPos: e.OpPos,
		Op:    e.Op,
		Y:     c.FromExpr(e.Y),
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
	tv, ok := c.Info.Types[e]
	if !ok {
		panic(fmt.Errorf(`failed to get type for %T at %s`, e, c.FileSet.Position(e.Pos()).String()))
	}
	return &ir.TypeExpr{
		TypePos:      e.Pos(),
		TypeAndValue: tv,
	}
}
