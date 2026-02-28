package irc

import (
	"fmt"
	"go/token"
)

type (
	Stmt interface {
		// String gets a simple representation for this statement.
		String() string

		// Pos is the position for this statement in the source file.
		// This may be [token.NoPos] if not specifically in the source.
		// This may be able to be used to look up the [ast.Node] that
		// became this IRC statement.
		Pos() token.Pos

		// stmt is an empty method used to compile time type check that
		// only statements are used for this interface.
		stmt()
	}

	// ExprStmt is a statement for an expression.
	ExprStmt struct {
		Expr Expr
	}

	// GotoStmt is a breaking block control that jumps to another block.
	GotoStmt struct {
		KeyPos token.Pos // the position of the keyword, e.g. `goto`, `return`
		Goto   *BlockRef // block to goto
	}

	// CallStmt is a breaking block control that call a function.
	CallStmt struct {
		KeyPos   token.Pos // the position of the keyword, e.g. `goto`, `return`
		Call     Expr      // the expression for the function to call, e.g. `fmt.Println`
		CallArgs []Expr    // the arguments to pass onto the call
		Follow   *BlockRef // block to goto when this call returns
	}

	// PanicStmt is a breaking block control that emits a panic.
	PanicStmt struct {
		PanicPos token.Pos // the position of the keyword, e.g. `panic`
		Value    Expr      // the value to panic
	}

	// SendStmt is a breaking block control that sends a value to a channel.
	SendStmt struct {
		ArrowPos token.Pos // the position of the channel arrow
		Channel  Expr      // the channel to send a value to
		Value    Expr      // the value to send
		Follow   *BlockRef // block to goto when this call returns
	}

	// ReceiveStmt is a breaking block control that receives a value from a channel.
	ReceiveStmt struct {
		ArrowPos token.Pos // the position of the channel arrow
		Channel  Expr      // the channel to receive a value from
		NeedOk   bool      // indicates the follow block should also receive the channel closed boolean
		Follow   *BlockRef // block to goto when this receive returns
	}

	// TODO: SelectStmt
	// TODO: LockStmt
	// TODO: SleepStmt
	// TODO: SuspendStmt
	// TODO: ExitThreadStmt
	// TODO: ExitAppStmt
	// TODO: MainSuspend

	// IfStmt is a statement for an if-statement.
	IfStmt struct {
		IfPos token.Pos // position of "if" keyword
		Cond  Expr      // condition of the if-statement
		Then  []Stmt    // then branch; should not be empty
		Else  []Stmt    // else branch; or empty for no else
	}

	// TODO: Add SwitchStmt and CaseClause
	// TODO: Add TypeSwitchStmt
)

var (
	_ Stmt = (*ExprStmt)(nil)
	_ Stmt = (*GotoStmt)(nil)
	_ Stmt = (*CallStmt)(nil)
	_ Stmt = (*PanicStmt)(nil)
	_ Stmt = (*SendStmt)(nil)
	_ Stmt = (*ReceiveStmt)(nil)
	_ Stmt = (*IfStmt)(nil)
)

func (s *ExprStmt) String() string { return fmt.Sprintf(`%s`, s.Expr) }
func (s *GotoStmt) String() string { return fmt.Sprintf(`goto(%s)`, s.Goto) }
func (s *CallStmt) String() string {
	return fmt.Sprintf(`%s(%s)->%s`, s.Call, csvString(s.CallArgs), s.Follow)
}
func (s *PanicStmt) String() string { return fmt.Sprintf(`panic(%s)`, s.Value) }
func (s *SendStmt) String() string  { return fmt.Sprintf(`(%s<-%s)->%s`, s.Channel, s.Value, s.Follow) }
func (s *ReceiveStmt) String() string {
	return fmt.Sprintf(`(%t<-%s)->%s`, s.NeedOk, s.Channel, s.Follow)
}
func (s *IfStmt) String() string {
	str := fmt.Sprintf("if %s {\n%s\n}", s.Cond, linesString(s.Then, `  `))
	if len(s.Else) > 0 {
		str += fmt.Sprintf(" else {\n%s\n}", linesString(s.Else, `  `))
	}
	return str
}

func (s *ExprStmt) Pos() token.Pos    { return s.Expr.Pos() }
func (s *GotoStmt) Pos() token.Pos    { return s.KeyPos }
func (s *CallStmt) Pos() token.Pos    { return s.KeyPos }
func (s *PanicStmt) Pos() token.Pos   { return s.PanicPos }
func (s *SendStmt) Pos() token.Pos    { return s.ArrowPos }
func (s *ReceiveStmt) Pos() token.Pos { return s.ArrowPos }
func (s *IfStmt) Pos() token.Pos      { return s.IfPos }

func (*ExprStmt) stmt()    {}
func (*GotoStmt) stmt()    {}
func (*CallStmt) stmt()    {}
func (*PanicStmt) stmt()   {}
func (*SendStmt) stmt()    {}
func (*ReceiveStmt) stmt() {}
func (*IfStmt) stmt()      {}
