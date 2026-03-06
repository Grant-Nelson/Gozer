package irc

import (
	"fmt"
	"go/token"

	"github.com/Grant-Nelson/Gozer/project/enums/branchType"
)

type (
	// Stmt is a code statement that can be put inside of a function
	// but has no type value like an expression.
	//
	// Differences between AST Stmt and IRC Stmt:
	// - BadStmt, log as error
	// - DeclStmt, use decl types
	// - EmptyStmt, skip
	// - IncDecStmt, perform with a unary expr
	// - AssignStmt, perform with a binary expr
	// - BlockStmt, skip since scoping is already known
	Stmt interface {
		Node

		// String gets a simple representation for this statement.
		String() string

		// stmtNode is an empty method used to compile time type check that
		// only statements duck-type to this interface.
		stmtNode()
	}

	// ValueDeclStmt is a node for const or var declarations as a statement.
	ValueDeclStmt struct {
		Decl *ValueDecl
	}

	// LabeledStmt is a node for a labeled statement.
	LabeledStmt struct {
		Label Ident
		Colon token.Pos // position of ":"
		Stmt  Stmt
	}

	// ExprStmt is a node for a (stand-alone) expression as a statement.
	ExprStmt struct {
		Expr Expr // expression
	}

	// SendStmt is a node for a send statement.
	SendStmt struct {
		Chan  Expr
		Arrow token.Pos // position of "<-"
		Value Expr
	}

	// GoStmt is a node for a go statement.
	GoStmt struct {
		Go   token.Pos // position of "go" keyword
		Call *CallExpr
	}

	// DeferStmt is a node for a defer statement.
	DeferStmt struct {
		Defer token.Pos // position of "defer" keyword
		Call  *CallExpr
	}

	// ReturnStmt is a node for a return statement.
	ReturnStmt struct {
		Return  token.Pos // position of "return" keyword
		Results []Expr    // result expressions; or nil
	}

	// BranchStmt is a node for a break, continue, goto, or fallthrough statement.
	BranchStmt struct {
		BranchTypePos token.Pos // position of branch type
		BranchType    branchType.BranchType
		Label         Ident // label name; or nil
	}

	// IfStmt is a statement for an if-statement.
	// The original Init field is moved out since scoping is already known.
	IfStmt struct {
		If   token.Pos // position of "if" keyword
		Cond Expr      // condition of the if-statement
		Then []Stmt    // then branch; should not be empty
		Else []Stmt    // else branch; or empty for no else
	}

	// SwitchStmt is a node for an expression switch statement.
	// The original Init field is moved out since scoping is already known.
	SwitchStmt struct {
		Switch token.Pos // position of "switch" keyword
		Tag    Expr      // tag expression; or nil
		Cases  []*CaseClause
	}

	// TypeSwitchStmt is a node for a type switch statement.
	// The original Init field is moved out since scoping is already known.
	TypeSwitchStmt struct {
		Switch token.Pos // position of "switch" keyword
		Assign Stmt      // x := y.(type) or y.(type)
		Cases  []*CaseClause
	}

	// SelectStmt is a node for a select statement.
	SelectStmt struct {
		Select token.Pos // position of "select" keyword
		Cases  []*CommClause
	}

	// ForStmt is a node for a for-loop statement.
	ForStmt struct {
		For  token.Pos // position of "for" keyword
		Init Stmt      // initialization statement; or nil
		Cond Expr      // condition; or nil
		Post Stmt      // post iteration statement; or nil
		Body []Stmt
	}

	// RangeStmt is a node for a for-range statement.
	RangeStmt struct {
		For        token.Pos   // position of "for" keyword
		Key, Value Expr        // Key, Value may be nil
		TokPos     token.Pos   // position of Tok; invalid if Key == nil
		Tok        token.Token // ILLEGAL if Key == nil, ASSIGN, DEFINE
		Range      token.Pos   // position of "range" keyword
		X          Expr        // value to range over
		Body       []Stmt
	}
)

var (
	_ Stmt = (*ValueDeclStmt)(nil)
	_ Stmt = (*LabeledStmt)(nil)
	_ Stmt = (*ExprStmt)(nil)
	_ Stmt = (*SendStmt)(nil)
	_ Stmt = (*GoStmt)(nil)
	_ Stmt = (*DeferStmt)(nil)
	_ Stmt = (*ReturnStmt)(nil)
	_ Stmt = (*BranchStmt)(nil)
	_ Stmt = (*IfStmt)(nil)
	_ Stmt = (*SwitchStmt)(nil)
	_ Stmt = (*TypeSwitchStmt)(nil)
	_ Stmt = (*SelectStmt)(nil)
	_ Stmt = (*ForStmt)(nil)
	_ Stmt = (*RangeStmt)(nil)
)

//====[String]==================================================================

func (s *ValueDeclStmt) String() string { return s.Decl.String() }
func (s *LabeledStmt) String() string   { return fmt.Sprintf("%s:\n%v", s.Label, s.Stmt) }
func (s *ExprStmt) String() string      { return fmt.Sprintf(`%v`, s.Expr) }
func (s *SendStmt) String() string      { return fmt.Sprintf(`%v<-%v`, s.Chan, s.Value) }
func (s *GoStmt) String() string        { return fmt.Sprintf(`go %v`, s.Call) }
func (s *DeferStmt) String() string     { return fmt.Sprintf(`defer %v`, s.Call) }
func (s *ReturnStmt) String() string    { return fmt.Sprintf(`return %s`, csvString(s.Results)) }
func (s *BranchStmt) String() string    { return fmt.Sprintf(`%s %v`, s.BranchType, s.Label) }
func (s *IfStmt) String() string {
	str := fmt.Sprintf("if %v {\n%s\n}", s.Cond, linesString(s.Then, `  `))
	if len(s.Else) > 0 {
		str += fmt.Sprintf(" else {\n%s\n}", linesString(s.Else, `  `))
	}
	return str
}
func (s *SwitchStmt) String() string {
	return fmt.Sprintf(`switch %v {\n%s\n}`, s.Tag, linesString(s.Cases, `  `))
}
func (s *TypeSwitchStmt) String() string {
	return fmt.Sprintf(`switch %v {\n%s\n}`, s.Assign, linesString(s.Cases, `  `))
}
func (s *SelectStmt) String() string {
	return fmt.Sprintf(`select {\n%s\n}`, linesString(s.Cases, `  `))
}
func (s *ForStmt) String() string {
	return fmt.Sprintf(`for %v; %v; %v {\n%v\n}`, s.Init, s.Cond, s.Post, linesString(s.Body, `  `))
}
func (s *RangeStmt) String() string {
	return fmt.Sprintf(`for %v, %v %v range %v {\n%v\n}`, s.Key, s.Value, s.Tok.String(), s.X, linesString(s.Body, `  `))
}

//====[Pos]=====================================================================

func (s *ValueDeclStmt) Pos() token.Pos  { return s.Decl.Pos() }
func (s *LabeledStmt) Pos() token.Pos    { return s.Label.Pos() }
func (s *ExprStmt) Pos() token.Pos       { return s.Expr.Pos() }
func (s *SendStmt) Pos() token.Pos       { return s.Chan.Pos() }
func (s *GoStmt) Pos() token.Pos         { return s.Go }
func (s *DeferStmt) Pos() token.Pos      { return s.Defer }
func (s *ReturnStmt) Pos() token.Pos     { return s.Return }
func (s *BranchStmt) Pos() token.Pos     { return s.BranchTypePos }
func (s *IfStmt) Pos() token.Pos         { return s.If }
func (s *SwitchStmt) Pos() token.Pos     { return s.Switch }
func (s *TypeSwitchStmt) Pos() token.Pos { return s.Switch }
func (s *SelectStmt) Pos() token.Pos     { return s.Select }
func (s *ForStmt) Pos() token.Pos        { return s.For }
func (s *RangeStmt) Pos() token.Pos      { return s.For }

//====[End]=====================================================================

func (s *ValueDeclStmt) End() token.Pos { return s.Decl.End() }
func (s *LabeledStmt) End() token.Pos   { return s.Stmt.End() }
func (s *ExprStmt) End() token.Pos      { return s.Expr.End() }
func (s *SendStmt) End() token.Pos      { return s.Value.End() }
func (s *GoStmt) End() token.Pos        { return s.Call.End() }
func (s *DeferStmt) End() token.Pos     { return s.Call.End() }
func (s *ReturnStmt) End() token.Pos {
	if p := endOfSlice(s.Results); p.IsValid() {
		return p
	}
	return s.Return + 6 // len("return")
}
func (s *BranchStmt) End() token.Pos {
	if s.Label != nil {
		return s.Label.End()
	}
	return token.Pos(int(s.BranchTypePos) + len(s.BranchType.String()))
}
func (s *IfStmt) End() token.Pos {
	if p := endOfSlice(s.Else); p.IsValid() {
		return p
	}
	return endOfSlice(s.Then)
}
func (s *SwitchStmt) End() token.Pos     { return endOfSlice(s.Cases) }
func (s *TypeSwitchStmt) End() token.Pos { return endOfSlice(s.Cases) }
func (s *SelectStmt) End() token.Pos     { return endOfSlice(s.Cases) }
func (s *ForStmt) End() token.Pos        { return endOfSlice(s.Body) }
func (s *RangeStmt) End() token.Pos      { return endOfSlice(s.Body) }

//====[stmtNode]================================================================

func (*ValueDeclStmt) stmtNode()  {}
func (*LabeledStmt) stmtNode()    {}
func (*ExprStmt) stmtNode()       {}
func (*SendStmt) stmtNode()       {}
func (*GoStmt) stmtNode()         {}
func (*DeferStmt) stmtNode()      {}
func (*ReturnStmt) stmtNode()     {}
func (*BranchStmt) stmtNode()     {}
func (*IfStmt) stmtNode()         {}
func (*SwitchStmt) stmtNode()     {}
func (*TypeSwitchStmt) stmtNode() {}
func (*SelectStmt) stmtNode()     {}
func (*ForStmt) stmtNode()        {}
func (*RangeStmt) stmtNode()      {}
