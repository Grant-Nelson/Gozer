package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

type (
	// Stmt is a code statement that can be put inside of a function
	// but has no type value like an expression.
	Stmt interface {
		Node

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
		SrcPos token.Pos
		Block  *BlockRef
	}

	// FuncCallStmt is a statement for jumping to another block or named function.
	// This is unique to IR and not a mirror of an AST node.
	FuncCallStmt struct {
		Ast    *ast.CallExpr
		Fun    ast.Expr // function expression
		Args   []ast.Expr
		Follow *BlockRef
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
		Ast ast.Stmt
		X   ast.Expr
	}

	// A SendStmt node represents a send statement.
	SendStmt struct {
		Ast   *ast.SendStmt
		Chan  ast.Expr
		Value ast.Expr
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
	_ Stmt = (*FuncCallStmt)(nil)
	_ Stmt = (*DeclStmt)(nil)
	_ Stmt = (*LabeledStmt)(nil)
	_ Stmt = (*ExprStmt)(nil)
	_ Stmt = (*SendStmt)(nil)
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

func (s *FuncCallStmt) String() string {
	return `call(` + nodeString(s.Fun) + `, ` + csvString(s.Args) + `|` + s.Follow.String() + `)`
}
func (s *DeclStmt) String() string    { return nodeString(s.Decl) }
func (s *LabeledStmt) String() string { return fmt.Sprintf("%s:\n%v", s.Label, s.Stmt) }
func (s *ExprStmt) String() string    { return nodeString(s.X) }
func (s *SendStmt) String() string {
	return fmt.Sprintf(`%s<-%s`, nodeString(s.Chan), nodeString(s.Value))
}

func (s *AssignStmt) String() string {
	return fmt.Sprintf(`%s%s%s`, csvString(s.Lhs), s.Tok.String(), csvString(s.Rhs))
}
func (s *GoStmt) String() string    { return fmt.Sprintf(`go %s`, nodeString(s.Call)) }
func (s *DeferStmt) String() string { return fmt.Sprintf(`defer %s`, nodeString(s.Call)) }
func (s *ReturnStmt) String() string {
	if len(s.Results) <= 0 {
		return `return`
	}
	return fmt.Sprintf(`return %s`, csvString(s.Results))
}
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
	return fmt.Sprintf("switch %s {\n%s\n}", nodeString(s.Tag), linesString(s.Body))
}

func (s *TypeSwitchStmt) String() string {
	return fmt.Sprintf("switch %v {\n%s\n}", s.Assign, linesString(s.Body))
}

func (s *SelectStmt) String() string {
	return fmt.Sprintf("select {\n%s\n}", linesString(s.Body))
}

func (s *ForStmt) String() string {
	return fmt.Sprintf("for %v; %s; %v {\n%v\n}",
		s.Init, nodeString(s.Cond), s.Post, linesString(s.Body))
}

func (s *RangeStmt) String() string {
	return fmt.Sprintf("for %v, %v %v range %v {\n%v\n}",
		s.Key, s.Value, s.Tok.String(), s.X, linesString(s.Body))
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

func (s *GotoBlockStmt) Pos() token.Pos  { return s.SrcPos }
func (s *FuncCallStmt) Pos() token.Pos   { return astPos(s.Ast) }
func (s *DeclStmt) Pos() token.Pos       { return astPos(s.Ast) }
func (s *LabeledStmt) Pos() token.Pos    { return astPos(s.Ast) }
func (s *ExprStmt) Pos() token.Pos       { return astPos(s.Ast) }
func (s *SendStmt) Pos() token.Pos       { return astPos(s.Ast) }
func (s *AssignStmt) Pos() token.Pos     { return astPos(s.Ast) }
func (s *GoStmt) Pos() token.Pos         { return astPos(s.Ast) }
func (s *DeferStmt) Pos() token.Pos      { return astPos(s.Ast) }
func (s *ReturnStmt) Pos() token.Pos     { return astPos(s.Ast) }
func (s *BranchStmt) Pos() token.Pos     { return astPos(s.Ast) }
func (s *StmtListStmt) Pos() token.Pos   { return astPos(s.Ast) }
func (s *IfStmt) Pos() token.Pos         { return astPos(s.Ast) }
func (s *SwitchStmt) Pos() token.Pos     { return astPos(s.Ast) }
func (s *TypeSwitchStmt) Pos() token.Pos { return astPos(s.Ast) }
func (s *SelectStmt) Pos() token.Pos     { return astPos(s.Ast) }
func (s *ForStmt) Pos() token.Pos        { return astPos(s.Ast) }
func (s *RangeStmt) Pos() token.Pos      { return astPos(s.Ast) }

func (*GotoBlockStmt) StmtNode()  {}
func (*FuncCallStmt) StmtNode()   {}
func (*DeclStmt) StmtNode()       {}
func (*LabeledStmt) StmtNode()    {}
func (*ExprStmt) StmtNode()       {}
func (*SendStmt) StmtNode()       {}
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

//===[IR helper methods]========================================================

func NewGotoBlockStmt(pos token.Pos, nextBlk *Block, args ...ast.Expr) *GotoBlockStmt {
	return &GotoBlockStmt{
		SrcPos: pos,
		Block: &BlockRef{
			Block: nextBlk,
			Args:  args,
		},
	}
}

// IsFlowControlStatement determines if the given statement
// is a flow control statement such as a return or branch.
func IsFlowControlStatement(s Stmt) bool {
	switch s.(type) {
	case *GotoBlockStmt, *ReturnStmt, *BranchStmt:
		return true
	}
	return false
}
