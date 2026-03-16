package ir

import (
	"go/ast"
	"go/types"
	"reflect"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type Converter struct {
	Info *types.Info
}

func isNotNil[T any](t T) bool {
	v := reflect.ValueOf(t)
	return v.IsValid() && !v.IsZero()
}

func fromNilSafeStmt[TIn ast.Stmt, TOut Stmt](s TIn, fn func(TIn) TOut) Stmt {
	if isNotNil(s) {
		if t := fn(s); isNotNil(t) {
			return t
		}
	}
	return nil
}

func (c *Converter) FromStmtSlice(ss []ast.Stmt) []Stmt {
	result := make([]Stmt, 0, len(ss))
	for _, s := range ss {
		result = append(result, c.ExpandStmt(c.FromStmt(s))...)
	}
	return result
}

func (c *Converter) FromStmt(s ast.Stmt) Stmt {
	switch s := s.(type) {
	case nil, *ast.BadStmt, *ast.EmptyStmt:
		return nil
	case *ast.DeclStmt:
		return fromNilSafeStmt(s, c.FromDeclStmt)
	case *ast.LabeledStmt:
		return fromNilSafeStmt(s, c.FromLabeledStmt)
	case *ast.ExprStmt:
		return fromNilSafeStmt(s, c.FromExprStmt)
	case *ast.SendStmt:
		return fromNilSafeStmt(s, c.FromSendStmt)
	case *ast.IncDecStmt:
		return fromNilSafeStmt(s, c.FromIncDecStmt)
	case *ast.AssignStmt:
		return fromNilSafeStmt(s, c.FromAssignStmt)
	case *ast.GoStmt:
		return fromNilSafeStmt(s, c.FromGoStmt)
	case *ast.DeferStmt:
		return fromNilSafeStmt(s, c.FromDeferStmt)
	case *ast.ReturnStmt:
		return fromNilSafeStmt(s, c.FromReturnStmt)
	case *ast.BranchStmt:
		return fromNilSafeStmt(s, c.FromBranchStmt)
	case *ast.BlockStmt:
		return fromNilSafeStmt(s, c.FromBlockStmt)
	case *ast.IfStmt:
		return fromNilSafeStmt(s, c.FromIfStmt)
	case *ast.SwitchStmt:
		return fromNilSafeStmt(s, c.FromSwitchStmt)
	case *ast.TypeSwitchStmt:
		return fromNilSafeStmt(s, c.FromTypeSwitchStmt)
	case *ast.SelectStmt:
		return fromNilSafeStmt(s, c.FromSelectStmt)
	case *ast.ForStmt:
		return fromNilSafeStmt(s, c.FromForStmt)
	case *ast.RangeStmt:
		return fromNilSafeStmt(s, c.FromRangeStmt)
	default:
		panic(faults.New(`unexpected AST statement type`).
			WithF(`type`, `%T`, s))
	}
}

func (c *Converter) FromDeclStmt(s *ast.DeclStmt) *DeclStmt {
	return &DeclStmt{Ast: s, Decl: s.Decl}
}

func (c *Converter) FromLabeledStmt(s *ast.LabeledStmt) *LabeledStmt {
	return &LabeledStmt{Ast: s, Label: s.Label, Stmt: c.FromStmt(s.Stmt)}
}

func (c *Converter) FromExprStmt(s *ast.ExprStmt) *ExprStmt {
	return &ExprStmt{Ast: s, X: s.X}
}

func (c *Converter) FromSendStmt(s *ast.SendStmt) *SendStmt {
	return &SendStmt{Ast: s, Chan: s.Chan, Value: s.Value}
}

func (c *Converter) FromIncDecStmt(s *ast.IncDecStmt) *ExprStmt {
	exp := &ast.UnaryExpr{OpPos: s.Pos(), Op: s.Tok, X: s.X}
	c.Info.Types[exp] = c.Info.Types[s.X]
	return &ExprStmt{Ast: s, X: exp}
}

func (c *Converter) FromAssignStmt(s *ast.AssignStmt) *AssignStmt {
	return &AssignStmt{Ast: s, Lhs: s.Lhs, Tok: s.Tok, Rhs: s.Rhs}
}

func (c *Converter) FromGoStmt(s *ast.GoStmt) *GoStmt {
	return &GoStmt{Ast: s, Call: s.Call}
}

func (c *Converter) FromDeferStmt(s *ast.DeferStmt) *DeferStmt {
	return &DeferStmt{Ast: s, Call: s.Call}
}

func (c *Converter) FromReturnStmt(s *ast.ReturnStmt) *ReturnStmt {
	return &ReturnStmt{Ast: s, Results: s.Results}
}

func (c *Converter) FromBranchStmt(s *ast.BranchStmt) *BranchStmt {
	return &BranchStmt{Ast: s, Tok: s.Tok, Label: s.Label}
}

func (c *Converter) FromBlockStmt(s *ast.BlockStmt) *StmtListStmt {
	return &StmtListStmt{Ast: s, List: c.FromStmtSlice(s.List)}
}

func (c *Converter) ExpandStmtSlice(ss []Stmt) []Stmt {
	st := make([]Stmt, 0, len(ss))
	for _, s := range ss {
		st = append(st, c.ExpandStmt(s)...)
	}
	return st
}

func (c *Converter) ExpandStmt(s Stmt) []Stmt {
	if s == nil {
		return []Stmt{}
	}
	if b, ok := s.(*StmtListStmt); ok {
		return c.ExpandStmtSlice(b.List)
	}
	return []Stmt{s}
}

func (c *Converter) FromIfStmt(s *ast.IfStmt) *IfStmt {
	return &IfStmt{
		Ast:  s,
		Init: c.FromStmt(s.Init),
		Cond: s.Cond,
		Body: c.ExpandStmt(c.FromBlockStmt(s.Body)),
		Else: c.ExpandStmt(c.FromStmt(s.Else)),
	}
}

func (c *Converter) FromCaseClause(s *ast.CaseClause) *CaseClause {
	return &CaseClause{
		Ast:  s,
		List: s.List,
		Body: c.FromStmtSlice(s.Body),
	}
}

func (c *Converter) FromCaseClauseSlice(s *ast.BlockStmt) []*CaseClause {
	if s == nil {
		return []*CaseClause{}
	}
	ccs := make([]*CaseClause, 0, len(s.List))
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

func (c *Converter) FromSwitchStmt(s *ast.SwitchStmt) *SwitchStmt {
	return &SwitchStmt{
		Ast:  s,
		Init: c.FromStmt(s.Init),
		Tag:  s.Tag,
		Body: c.FromCaseClauseSlice(s.Body),
	}
}

func (c *Converter) FromTypeSwitchStmt(s *ast.TypeSwitchStmt) *TypeSwitchStmt {
	return &TypeSwitchStmt{
		Ast:    s,
		Init:   c.FromStmt(s.Init),
		Assign: c.FromStmt(s.Assign),
		Body:   c.FromCaseClauseSlice(s.Body),
	}
}

func (c *Converter) FromCommClause(s *ast.CommClause) *CommClause {
	return &CommClause{
		Ast:  s,
		Comm: c.FromStmt(s.Comm),
		Body: c.FromStmtSlice(s.Body),
	}
}

func (c *Converter) FromCommClauseSlice(s *ast.BlockStmt) []*CommClause {
	if s == nil {
		return []*CommClause{}
	}
	ccs := make([]*CommClause, 0, len(s.List))
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

func (c *Converter) FromSelectStmt(s *ast.SelectStmt) *SelectStmt {
	return &SelectStmt{
		Ast:  s,
		Body: c.FromCommClauseSlice(s.Body),
	}
}

func (c *Converter) FromForStmt(s *ast.ForStmt) *ForStmt {
	return &ForStmt{
		Ast:  s,
		Init: c.FromStmt(s.Init),
		Cond: s.Cond,
		Post: c.FromStmt(s.Post),
		Body: c.ExpandStmt(c.FromBlockStmt(s.Body)),
	}
}

func (c *Converter) FromRangeStmt(s *ast.RangeStmt) *RangeStmt {
	return &RangeStmt{
		Ast:   s,
		Key:   s.Key,
		Value: s.Value,
		X:     s.X,
		Body:  c.ExpandStmt(c.FromBlockStmt(s.Body)),
	}
}
