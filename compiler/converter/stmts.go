package converter

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/assert"
	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/iterator"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
)

func (c *Converter) FromStmtSlice(ss []ast.Stmt) []ir.Stmt {
	result := make([]ir.Stmt, 0, len(ss))
	for _, s := range ss {
		result = append(result, c.ExpandStmt(c.FromStmt(s))...)
	}
	return result
}

func (c *Converter) SimplifyStmt(s ir.Stmt) ir.Stmt {
	if b, ok := s.(*ir.StmtListStmt); ok {
		b.List = iterator.NotZero(iterator.Iterate(c.ExpandStmtSlice(b.List)...)).ToSlice()
		switch len(b.List) {
		case 0:
			return nil
		case 1:
			return b.List[0]
		default:
			return b
		}
	}
	return s
}

func (c *Converter) ExpandStmtSlice(ss []ir.Stmt) []ir.Stmt {
	st := make([]ir.Stmt, 0, len(ss))
	for _, s := range ss {
		for x := range iterator.NotZero(iterator.Iterate(c.ExpandStmt(s)...)) {
			st = append(st, x)
		}
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
		c.addFault(faults.New(`unexpected AST statement type`).
			WithF(`type`, `%T`, s).
			With(`pos`, c.pos(s.Pos())))
		return nil
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
			c.addFault(faults.New(`expected case clause`).
				WithF(`type`, `%T`, ct).
				With(`pos`, c.pos(ct.Pos())).
				With(`index`, i))
			continue
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
			c.addFault(faults.New(`expected comm clause`).
				WithF(`type`, `%T`, ct).
				With(`pos`, c.pos(ct.Pos())).
				With(`index`, i))
			continue
		}
		ccs = append(ccs, c.FromCommClause(cc))
	}
	return ccs
}

func (c *Converter) FromDeclStmt(s *ast.DeclStmt) ir.Stmt {
	if s == nil {
		return nil
	}
	return c.FromDecl(s.Decl)
}

func (c *Converter) FromDeferStmt(s *ast.DeferStmt) *ir.DeferStmt {
	if s == nil {
		return nil
	}
	return &ir.DeferStmt{
		Defer: s.Defer,
		Call:  c.FromCallExpr(s.Call),
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
	return &ir.ExprStmt{
		X: &ir.UnaryExpr{
			OpPos: s.Pos(),
			Op:    c.FromUnaryOp(s.Tok),
			X:     c.FromExpr(s.X),
		},
	}
}

func (c *Converter) FromForStmt(s *ast.ForStmt) *ir.ForStmt {
	if s == nil {
		return nil
	}
	return &ir.ForStmt{
		ForPos: s.For,
		Init:   c.FromStmt(s.Init),
		Cond:   c.FromExpr(s.Cond),
		Post:   c.FromStmt(s.Post),
		Body:   c.ExpandStmt(c.FromBlockStmt(s.Body)),
	}
}

func (c *Converter) FromGoStmt(s *ast.GoStmt) *ir.GoStmt {
	if s == nil {
		return nil
	}
	return &ir.GoStmt{
		GoPos: s.Go,
		Call:  c.FromCallExpr(s.Call),
	}
}

func (c *Converter) FromIfStmt(s *ast.IfStmt) *ir.IfStmt {
	if s == nil {
		return nil
	}
	return &ir.IfStmt{
		IfPos: s.If,
		Init:  c.FromStmt(s.Init),
		Cond:  c.FromExpr(s.Cond),
		Body:  c.ExpandStmt(c.FromBlockStmt(s.Body)),
		Else:  c.ExpandStmt(c.FromStmt(s.Else)),
	}
}

func (c *Converter) FromLabeledStmt(s *ast.LabeledStmt) *ir.LabeledStmt {
	if s == nil {
		return nil
	}
	return &ir.LabeledStmt{
		Label: c.FromIdent(s.Label),
		Stmt:  c.FromStmt(s.Stmt),
	}
}

func (c *Converter) FromRangeStmt(s *ast.RangeStmt) *ir.RangeStmt {
	if s == nil {
		return nil
	}
	return &ir.RangeStmt{
		ForPos: s.For,
		Key:    c.FromExpr(s.Key),
		Value:  c.FromExpr(s.Value),
		Define: s.Tok == token.DEFINE,
		X:      c.FromExpr(s.X),
		Body:   c.ExpandStmt(c.FromBlockStmt(s.Body)),
	}
}

func (c *Converter) FromReturnStmt(s *ast.ReturnStmt) *ir.ReturnStmt {
	if s == nil {
		return nil
	}
	return &ir.ReturnStmt{
		ReturnPos: s.Return,
		Results:   c.FromExprSlice(s.Results),
	}
}

func (c *Converter) FromSelectStmt(s *ast.SelectStmt) *ir.SelectStmt {
	if s == nil {
		return nil
	}
	return &ir.SelectStmt{
		SelectPos: s.Select,
		Body:      c.FromCommClauseSlice(s.Body),
	}
}

func (c *Converter) FromSendStmt(s *ast.SendStmt) *ir.SendStmt {
	if s == nil {
		return nil
	}
	return &ir.SendStmt{
		ArrowPos: s.Arrow,
		Chan:     c.FromExpr(s.Chan),
		Value:    c.FromExpr(s.Value),
	}
}

func (c *Converter) FromBlockStmt(s *ast.BlockStmt) *ir.StmtListStmt {
	if s == nil {
		return nil
	}
	return &ir.StmtListStmt{
		List: c.FromStmtSlice(s.List),
	}
}

func (c *Converter) FromSwitchStmt(s *ast.SwitchStmt) *ir.SwitchStmt {
	if s == nil {
		return nil
	}
	return &ir.SwitchStmt{
		SwitchPos: s.Switch,
		Init:      c.FromStmt(s.Init),
		Tag:       c.FromExpr(s.Tag),
		Body:      c.FromCaseClauseSlice(s.Body),
	}
}

func (c *Converter) FromTypeSwitchStmt(s *ast.TypeSwitchStmt) *ir.TypeSwitchStmt {
	if s == nil {
		return nil
	}
	return &ir.TypeSwitchStmt{
		SwitchPos: s.Switch,
		Init:      c.FromStmt(s.Init),
		Assign:    c.FromStmt(s.Assign),
		Body:      c.FromCaseClauseSlice(s.Body),
	}
}

func (c *Converter) FromFuncDecl(astFunc *ast.FuncDecl) *ir.Func {
	if astFunc == nil {
		return nil
	}
	fn := &ir.Func{
		FuncPos: astFunc.Type.Func,
		Name:    c.FromIdent(astFunc.Name),
	}
	c.setInitialBlock(fn, astFunc.Body, astFunc.Type)

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

func (c *Converter) FromFuncLit(astFunc *ast.FuncLit) *ir.Func {
	if astFunc == nil {
		return nil
	}
	pos := astFunc.Type.Func
	fn := &ir.Func{
		FuncPos: pos,
	}
	c.setInitialBlock(fn, astFunc.Body, astFunc.Type)
	return fn
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
				Name: c.FromIdent(name),
				Expr: c.FromExpr(field.Type),
			}
			params = append(params, param)
		}
	}
	return params
}

// setInitialBlock creates an initial block, populate it with current statements, add params.
// Named return values are appended after input params so they are in scope
// for the function body just like declared locals.
func (c *Converter) setInitialBlock(fn *ir.Func, fnBody *ast.BlockStmt, fnType *ast.FuncType) *ir.Block {
	body := c.ExpandStmt(c.FromBlockStmt(fnBody))
	params := append(c.ExpandParams(fnType), c.ExpandResults(fnType)...)
	return fn.NewBlock(`initial`, body, params)
}
