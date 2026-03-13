package ir

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type (
	// Stmt is a code statement that can be put inside of a function
	// but has no type value like an expression.
	Stmt interface {

		// String returns a human-readable text representing the statement
		// for debugging and testing. The output must be consistent.
		String() string

		// StmtNode is an empty method used to compile time type check that
		// only statements duck-type to this interface.
		StmtNode()
	}

	//===[Flow Control]=========================================================

	// GotoBlockStmt is a statement for a goto jump to another block.
	// This is unique to IR and not a mirror of an AST node.
	GotoBlockStmt struct {
		Block *BlockRef
	}

	//===[AST Mirrors]==========================================================

	// A DeclStmt node represents a declaration in a statement list.
	DeclStmt struct {
		Ast  *ast.DeclStmt
		Decl ast.Decl // *GenDecl with CONST, TYPE, or VAR token
	}

	// A LabeledStmt node represents a labeled statement.
	LabeledStmt struct {
		Ast   *ast.LabeledStmt
		Label *ast.Ident // same as AST so works with types.Info
		Stmt  Stmt
	}

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.
	ExprStmt struct {
		Ast *ast.ExprStmt
		X   ast.Expr
	}

	// A SendStmt node represents a send statement.
	SendStmt struct {
		Ast   *ast.SendStmt
		Chan  ast.Expr
		Value ast.Expr
	}

	// An IncDecStmt node represents an increment or decrement statement.
	IncDecStmt struct {
		Ast *ast.IncDecStmt
		X   ast.Expr
		Tok token.Token // INC or DEC
	}

	// An AssignStmt node represents an assignment or
	// a short variable declaration.
	AssignStmt struct {
		Ast *ast.AssignStmt
		Lhs []ast.Expr
		Tok token.Token // assignment token, DEFINE
		Rhs []ast.Expr
	}

	// A GoStmt node represents a go statement.
	GoStmt struct {
		Ast  *ast.GoStmt
		Call *ast.CallExpr
	}

	// A DeferStmt node represents a defer statement.
	DeferStmt struct {
		Ast  *ast.DeferStmt
		Call *ast.CallExpr
	}

	// A ReturnStmt node represents a return statement.
	ReturnStmt struct {
		Ast     *ast.ReturnStmt
		Results []ast.Expr // result expressions; or nil
	}

	// A BranchStmt node represents a break, continue, goto,
	// or fallthrough statement.
	BranchStmt struct {
		Ast   *ast.BranchStmt
		Tok   token.Token // keyword token (BREAK, CONTINUE, GOTO, FALLTHROUGH)
		Label *ast.Ident  // label name; or nil
	}

	// A StmtListStmt node represents a braced statement list.
	StmtListStmt struct {
		Ast  *ast.BlockStmt
		List []Stmt
	}

	// An IfStmt node represents an if statement.
	IfStmt struct {
		Ast  *ast.IfStmt
		Init Stmt     // initialization statement; or nil
		Cond ast.Expr // condition
		Body []Stmt
		Else []Stmt
	}

	// A SwitchStmt node represents an expression switch statement.
	SwitchStmt struct {
		Ast  *ast.SwitchStmt
		Init Stmt     // initialization statement; or nil
		Tag  ast.Expr // tag expression; or nil
		Body []*CaseClause
	}

	// A TypeSwitchStmt node represents a type switch statement.
	TypeSwitchStmt struct {
		Ast    *ast.TypeSwitchStmt
		Init   Stmt // initialization statement; or nil
		Assign Stmt // x := y.(type) or y.(type)
		Body   []*CaseClause
	}

	// A SelectStmt node represents a select statement.
	SelectStmt struct {
		Ast  *ast.SelectStmt
		Body []*CommClause
	}

	// A ForStmt represents a for statement.
	ForStmt struct {
		Ast  *ast.ForStmt
		Init Stmt     // initialization statement; or nil
		Cond ast.Expr // condition; or nil
		Post Stmt     // post iteration statement; or nil
		Body []Stmt
	}

	// A RangeStmt represents a for statement with a range clause.
	RangeStmt struct {
		Ast   *ast.RangeStmt
		Key   ast.Expr    // Key may be nil
		Value ast.Expr    // Value may be nil
		Tok   token.Token // ILLEGAL if Key == nil, ASSIGN, DEFINE
		X     ast.Expr    // value to range over
		Body  []Stmt
	}

	//===[Clauses]==============================================================

	// A CaseClause represents a case of an expression or type switch statement.
	CaseClause struct {
		Ast  *ast.CaseClause
		List []ast.Expr // list of expressions or types; nil means default case
		Body []Stmt     // statement list; or nil
	}

	// A CommClause node represents a case of a select statement.
	CommClause struct {
		Ast  *ast.CommClause
		Comm Stmt   // send or receive statement; nil means default case
		Body []Stmt // statement list; or nil
	}
)

var (
	_ Stmt = (*GotoBlockStmt)(nil)
	_ Stmt = (*DeclStmt)(nil)
	_ Stmt = (*LabeledStmt)(nil)
	_ Stmt = (*ExprStmt)(nil)
	_ Stmt = (*SendStmt)(nil)
	_ Stmt = (*IncDecStmt)(nil)
	_ Stmt = (*AssignStmt)(nil)
	_ Stmt = (*GoStmt)(nil)
	_ Stmt = (*DeferStmt)(nil)
	_ Stmt = (*ReturnStmt)(nil)
	_ Stmt = (*BranchStmt)(nil)
	_ Stmt = (*StmtListStmt)(nil)
	_ Stmt = (*IfStmt)(nil)
	_ Stmt = (*SwitchStmt)(nil)
	_ Stmt = (*TypeSwitchStmt)(nil)
	_ Stmt = (*SelectStmt)(nil)
	_ Stmt = (*ForStmt)(nil)
	_ Stmt = (*RangeStmt)(nil)
)

func (s *GotoBlockStmt) String() string {
	return `goto(` + s.Block.String() + `)`
}
func (s *DeclStmt) String() string    { return nodeString(s.Decl) }
func (s *LabeledStmt) String() string { return fmt.Sprintf("%s:\n%v", s.Label, s.Stmt) }
func (s *ExprStmt) String() string    { return nodeString(s.X) }
func (s *SendStmt) String() string {
	return fmt.Sprintf(`%s<-%s`, nodeString(s.Chan), nodeString(s.Value))
}
func (s *IncDecStmt) String() string { return fmt.Sprintf(`%s%s`, nodeString(s.X), s.Tok.String()) }
func (s *AssignStmt) String() string {
	return fmt.Sprintf(`%s%s%s`, csvString(s.Lhs), s.Tok.String(), csvString(s.Rhs))
}
func (s *GoStmt) String() string       { return fmt.Sprintf(`go %s`, nodeString(s.Call)) }
func (s *DeferStmt) String() string    { return fmt.Sprintf(`defer %s`, nodeString(s.Call)) }
func (s *ReturnStmt) String() string   { return fmt.Sprintf(`return %s`, csvString(s.Results)) }
func (s *BranchStmt) String() string   { return fmt.Sprintf(`%s %v`, s.Tok.String(), s.Label) }
func (s *StmtListStmt) String() string { return fmt.Sprintf("{\n%s\n}", linesString(s.List)) }
func (s *IfStmt) String() string {
	str := fmt.Sprintf("if %s {\n%s\n}", nodeString(s.Cond), linesString(s.Body))
	if len(s.Else) > 0 {
		str += fmt.Sprintf(" else {\n%s\n}", linesString(s.Else))
	}
	return str
}
func (s *SwitchStmt) String() string {
	return fmt.Sprintf(`switch %s {\n%s\n}`, nodeString(s.Tag), linesString(s.Body))
}
func (s *TypeSwitchStmt) String() string {
	return fmt.Sprintf(`switch %v {\n%s\n}`, s.Assign, linesString(s.Body))
}
func (s *SelectStmt) String() string {
	return fmt.Sprintf(`select {\n%s\n}`, linesString(s.Body))
}
func (s *ForStmt) String() string {
	return fmt.Sprintf(`for %v; %s; %v {\n%v\n}`, s.Init, nodeString(s.Cond), s.Post, linesString(s.Body))
}
func (s *RangeStmt) String() string {
	return fmt.Sprintf(`for %v, %v %v range %v {\n%v\n}`, s.Key, s.Value, s.Tok.String(), s.X, linesString(s.Body))
}
func (s *CaseClause) String() string {
	str := `default:`
	if len(s.List) > 0 {
		str = fmt.Sprintf(`case %v:`, csvString(s.List))
	}
	if len(s.Body) > 0 {
		str += "\n" + linesString(s.Body)
	}
	return str
}
func (s *CommClause) String() string {
	str := `default:`
	if s.Comm != nil {
		str = fmt.Sprintf(`case %v:`, s.Comm)
	}
	if len(s.Body) > 0 {
		str += "\n" + linesString(s.Body)
	}
	return str
}

func (*GotoBlockStmt) StmtNode()  {}
func (*DeclStmt) StmtNode()       {}
func (*LabeledStmt) StmtNode()    {}
func (*ExprStmt) StmtNode()       {}
func (*SendStmt) StmtNode()       {}
func (*IncDecStmt) StmtNode()     {}
func (*AssignStmt) StmtNode()     {}
func (*GoStmt) StmtNode()         {}
func (*DeferStmt) StmtNode()      {}
func (*ReturnStmt) StmtNode()     {}
func (*BranchStmt) StmtNode()     {}
func (*StmtListStmt) StmtNode()   {}
func (*IfStmt) StmtNode()         {}
func (*SwitchStmt) StmtNode()     {}
func (*TypeSwitchStmt) StmtNode() {}
func (*SelectStmt) StmtNode()     {}
func (*ForStmt) StmtNode()        {}
func (*RangeStmt) StmtNode()      {}

//===[from AST converters]======================================================

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

func fromStmtSlice(ss []ast.Stmt) []Stmt {
	result := make([]Stmt, 0, len(ss))
	for _, s := range ss {
		result = append(result, expandStmt(fromStmt(s))...)
	}
	return result
}

func fromStmt(s ast.Stmt) Stmt {
	switch s := s.(type) {
	case nil, *ast.BadStmt, *ast.EmptyStmt:
		return nil
	case *ast.DeclStmt:
		return fromNilSafeStmt(s, fromDeclStmt)
	case *ast.LabeledStmt:
		return fromNilSafeStmt(s, fromLabeledStmt)
	case *ast.ExprStmt:
		return fromNilSafeStmt(s, fromExprStmt)
	case *ast.SendStmt:
		return fromNilSafeStmt(s, fromSendStmt)
	case *ast.IncDecStmt:
		return fromNilSafeStmt(s, fromIncDecStmt)
	case *ast.AssignStmt:
		return fromNilSafeStmt(s, fromAssignStmt)
	case *ast.GoStmt:
		return fromNilSafeStmt(s, fromGoStmt)
	case *ast.DeferStmt:
		return fromNilSafeStmt(s, fromDeferStmt)
	case *ast.ReturnStmt:
		return fromNilSafeStmt(s, fromReturnStmt)
	case *ast.BranchStmt:
		return fromNilSafeStmt(s, fromBranchStmt)
	case *ast.BlockStmt:
		return fromNilSafeStmt(s, fromBlockStmt)
	case *ast.IfStmt:
		return fromNilSafeStmt(s, fromIfStmt)
	case *ast.SwitchStmt:
		return fromNilSafeStmt(s, fromSwitchStmt)
	case *ast.TypeSwitchStmt:
		return fromNilSafeStmt(s, fromTypeSwitchStmt)
	case *ast.SelectStmt:
		return fromNilSafeStmt(s, fromSelectStmt)
	case *ast.ForStmt:
		return fromNilSafeStmt(s, fromForStmt)
	case *ast.RangeStmt:
		return fromNilSafeStmt(s, fromRangeStmt)
	default:
		panic(faults.New(`unexpected AST statement type`).
			WithF(`type`, `%T`, s))
	}
}

func fromDeclStmt(s *ast.DeclStmt) *DeclStmt {
	return &DeclStmt{Ast: s, Decl: s.Decl}
}

func fromLabeledStmt(s *ast.LabeledStmt) *LabeledStmt {
	return &LabeledStmt{Ast: s, Label: s.Label, Stmt: fromStmt(s.Stmt)}
}

func fromExprStmt(s *ast.ExprStmt) *ExprStmt {
	return &ExprStmt{Ast: s, X: s.X}
}

func fromSendStmt(s *ast.SendStmt) *SendStmt {
	return &SendStmt{Ast: s, Chan: s.Chan, Value: s.Value}
}

func fromIncDecStmt(s *ast.IncDecStmt) *IncDecStmt {
	return &IncDecStmt{Ast: s, X: s.X, Tok: s.Tok}
}

func fromAssignStmt(s *ast.AssignStmt) *AssignStmt {
	return &AssignStmt{Ast: s, Lhs: s.Lhs, Tok: s.Tok, Rhs: s.Rhs}
}

func fromGoStmt(s *ast.GoStmt) *GoStmt {
	return &GoStmt{Ast: s, Call: s.Call}
}

func fromDeferStmt(s *ast.DeferStmt) *DeferStmt {
	return &DeferStmt{Ast: s, Call: s.Call}
}

func fromReturnStmt(s *ast.ReturnStmt) *ReturnStmt {
	return &ReturnStmt{Ast: s, Results: s.Results}
}

func fromBranchStmt(s *ast.BranchStmt) *BranchStmt {
	return &BranchStmt{Ast: s, Tok: s.Tok, Label: s.Label}
}

func fromBlockStmt(s *ast.BlockStmt) *StmtListStmt {
	return &StmtListStmt{Ast: s, List: fromStmtSlice(s.List)}
}

func expandStmtSlice(ss []Stmt) []Stmt {
	st := make([]Stmt, 0, len(ss))
	for _, s := range ss {
		st = append(st, expandStmt(s)...)
	}
	return st
}

func expandStmt(s Stmt) []Stmt {
	if s == nil {
		return []Stmt{}
	}
	if b, ok := s.(*StmtListStmt); ok {
		return expandStmtSlice(b.List)
	}
	return []Stmt{s}
}

func fromIfStmt(s *ast.IfStmt) *IfStmt {
	return &IfStmt{Ast: s, Init: fromStmt(s.Init), Cond: s.Cond,
		Body: expandStmt(fromBlockStmt(s.Body)),
		Else: expandStmt(fromStmt(s.Else))}
}

func fromCaseClause(s *ast.CaseClause) *CaseClause {
	return &CaseClause{Ast: s, List: s.List, Body: fromStmtSlice(s.Body)}
}

func fromCaseClauseSlice(s *ast.BlockStmt) []*CaseClause {
	if s == nil {
		return []*CaseClause{}
	}
	ccs := make([]*CaseClause, 0, len(s.List))
	for i, c := range s.List {
		if c == nil {
			continue
		}
		cc, ok := c.(*ast.CaseClause)
		if !ok {
			panic(faults.New(`expected case clause`).
				WithF(`type`, `%T`, c).
				With(`index`, i))
		}
		ccs = append(ccs, fromCaseClause(cc))
	}
	return ccs
}

func fromSwitchStmt(s *ast.SwitchStmt) *SwitchStmt {
	return &SwitchStmt{Ast: s, Init: fromStmt(s.Init), Tag: s.Tag,
		Body: fromCaseClauseSlice(s.Body)}
}

func fromTypeSwitchStmt(s *ast.TypeSwitchStmt) *TypeSwitchStmt {
	return &TypeSwitchStmt{Ast: s, Init: fromStmt(s.Init),
		Assign: fromStmt(s.Assign),
		Body:   fromCaseClauseSlice(s.Body)}
}

func fromCommClause(s *ast.CommClause) *CommClause {
	return &CommClause{Ast: s, Comm: fromStmt(s.Comm),
		Body: fromStmtSlice(s.Body)}
}

func fromCommClauseSlice(s *ast.BlockStmt) []*CommClause {
	if s == nil {
		return []*CommClause{}
	}
	ccs := make([]*CommClause, 0, len(s.List))
	for i, c := range s.List {
		if c == nil {
			continue
		}
		cc, ok := c.(*ast.CommClause)
		if !ok {
			panic(faults.New(`expected comm clause`).
				WithF(`type`, `%T`, c).
				With(`index`, i))
		}
		ccs = append(ccs, fromCommClause(cc))
	}
	return ccs
}

func fromSelectStmt(s *ast.SelectStmt) *SelectStmt {
	return &SelectStmt{Ast: s, Body: fromCommClauseSlice(s.Body)}
}

func fromForStmt(s *ast.ForStmt) *ForStmt {
	return &ForStmt{Ast: s, Init: fromStmt(s.Init), Cond: s.Cond,
		Post: fromStmt(s.Post), Body: expandStmt(fromBlockStmt(s.Body))}
}

func fromRangeStmt(s *ast.RangeStmt) *RangeStmt {
	return &RangeStmt{Ast: s, Key: s.Key, Value: s.Value, X: s.X,
		Body: expandStmt(fromBlockStmt(s.Body))}
}
