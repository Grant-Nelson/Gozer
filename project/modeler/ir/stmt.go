package ir

import (
	"go/ast"
	"go/token"

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

	// A BlockStmt node represents a braced statement list.
	BlockStmt struct {
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
	_ Stmt = (*BlockStmt)(nil)
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

func (s *DeclStmt) String() string       { return `DeclStmt` }       // TODO: Implement
func (s *LabeledStmt) String() string    { return `LabeledStmt` }    // TODO: Implement
func (s *ExprStmt) String() string       { return `ExprStmt` }       // TODO: Implement
func (s *SendStmt) String() string       { return `SendStmt` }       // TODO: Implement
func (s *IncDecStmt) String() string     { return `IncDecStmt` }     // TODO: Implement
func (s *AssignStmt) String() string     { return `AssignStmt` }     // TODO: Implement
func (s *GoStmt) String() string         { return `GoStmt` }         // TODO: Implement
func (s *DeferStmt) String() string      { return `DeferStmt` }      // TODO: Implement
func (s *ReturnStmt) String() string     { return `ReturnStmt` }     // TODO: Implement
func (s *BranchStmt) String() string     { return `BranchStmt` }     // TODO: Implement
func (s *BlockStmt) String() string      { return `BlockStmt` }      // TODO: Implement
func (s *IfStmt) String() string         { return `IfStmt` }         // TODO: Implement
func (s *SwitchStmt) String() string     { return `SwitchStmt` }     // TODO: Implement
func (s *TypeSwitchStmt) String() string { return `TypeSwitchStmt` } // TODO: Implement
func (s *SelectStmt) String() string     { return `SelectStmt` }     // TODO: Implement
func (s *ForStmt) String() string        { return `ForStmt` }        // TODO: Implement
func (s *RangeStmt) String() string      { return `RangeStmt` }      // TODO: Implement
func (c *CaseClause) String() string     { return `CaseClause` }     // TODO: Implement
func (c *CommClause) String() string     { return `CommClause` }     // TODO: Implement

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
func (*BlockStmt) StmtNode()      {}
func (*IfStmt) StmtNode()         {}
func (*SwitchStmt) StmtNode()     {}
func (*TypeSwitchStmt) StmtNode() {}
func (*SelectStmt) StmtNode()     {}
func (*ForStmt) StmtNode()        {}
func (*RangeStmt) StmtNode()      {}

//===[from AST]=================================================================

func fromStmtSlice(ss []ast.Stmt) []Stmt {
	result := make([]Stmt, len(ss))
	for i, s := range ss {
		result[i] = fromStmt(s)
	}
	return result
}

func fromStmt(s ast.Stmt) Stmt {
	switch s := s.(type) {
	case nil, *ast.BadStmt, *ast.EmptyStmt:
		return nil
	case *ast.DeclStmt:
		return fromDeclStmt(s)
	case *ast.LabeledStmt:
		return fromLabeledStmt(s)
	case *ast.ExprStmt:
		return fromExprStmt(s)
	case *ast.SendStmt:
		return fromSendStmt(s)
	case *ast.IncDecStmt:
		return fromIncDecStmt(s)
	case *ast.AssignStmt:
		return fromAssignStmt(s)
	case *ast.GoStmt:
		return fromGoStmt(s)
	case *ast.DeferStmt:
		return fromDeferStmt(s)
	case *ast.ReturnStmt:
		return fromReturnStmt(s)
	case *ast.BranchStmt:
		return fromBranchStmt(s)
	case *ast.BlockStmt:
		return fromBlockStmt(s)
	case *ast.IfStmt:
		return fromIfStmt(s)
	case *ast.SwitchStmt:
		return fromSwitchStmt(s)
	case *ast.TypeSwitchStmt:
		return fromTypeSwitchStmt(s)
	case *ast.SelectStmt:
		return fromSelectStmt(s)
	case *ast.ForStmt:
		return fromForStmt(s)
	case *ast.RangeStmt:
		return fromRangeStmt(s)
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
	return &BranchStmt{Ast: s, Label: s.Label}
}

func fromBlockStmt(s *ast.BlockStmt) *BlockStmt {
	return &BlockStmt{Ast: s, List: fromStmtSlice(s.List)}
}

func fromIfStmt(s *ast.IfStmt) *IfStmt {
	return &IfStmt{Ast: s, Init: fromStmt(s.Init), Cond: s.Cond,
		Body: []Stmt{fromBlockStmt(s.Body)},
		Else: []Stmt{fromStmt(s.Else)}}
}

func fromCaseClause(s *ast.CaseClause) *CaseClause {
	return &CaseClause{Ast: s, List: s.List, Body: fromStmtSlice(s.Body)}
}

func fromCaseClauseSlice(s *ast.BlockStmt) []*CaseClause {
	ccs := make([]*CaseClause, len(s.List))
	for i, c := range s.List {
		cc, ok := c.(*ast.CaseClause)
		if !ok {
			panic(faults.New(`expected case clause`).
				WithF(`type`, `%T`, c).
				With(`index`, i))
		}
		ccs[i] = fromCaseClause(cc)
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
	ccs := make([]*CommClause, len(s.List))
	for i, c := range s.List {
		cc, ok := c.(*ast.CommClause)
		if !ok {
			panic(faults.New(`expected comm clause`).
				WithF(`type`, `%T`, c).
				With(`index`, i))
		}
		ccs[i] = fromCommClause(cc)
	}
	return ccs
}

func fromSelectStmt(s *ast.SelectStmt) *SelectStmt {
	return &SelectStmt{Ast: s, Body: fromCommClauseSlice(s.Body)}
}

func fromForStmt(s *ast.ForStmt) *ForStmt {
	return &ForStmt{Ast: s, Init: fromStmt(s.Init), Cond: s.Cond,
		Post: fromStmt(s.Post), Body: []Stmt{fromBlockStmt(s.Body)}}
}

func fromRangeStmt(s *ast.RangeStmt) *RangeStmt {
	return &RangeStmt{Ast: s, Key: s.Key, Value: s.Value, X: s.X,
		Body: []Stmt{fromBlockStmt(s.Body)}}
}
