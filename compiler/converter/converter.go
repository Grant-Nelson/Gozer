package converter

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/branchKind"
)

type Converter struct {
	Info *types.Info
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

	// TODO: FINISH IMPLEMENTING

	return nil
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
		Ast: s,
		X:   c.FromExpr(s.X),
	}
}

func (c *Converter) FromIncDecStmt(s *ast.IncDecStmt) *ir.ExprStmt {
	if s == nil {
		return nil
	}
	exp := &ast.UnaryExpr{OpPos: s.Pos(), Op: s.Tok, X: s.X}
	c.Info.Types[exp] = c.Info.Types[s.X]
	return &ir.ExprStmt{
		Ast: s,
		X:   c.FromExpr(exp),
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
		Call: s.Call,
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
