package irc

import (
	"fmt"
	"go/token"
)

type (
	FlowCtrl interface {
		Stmt

		// flowCtrl is an empty method used to compile time type
		// check that only flow controls duck-type to this interface.
		flowCtrl()
	}

	// GotoFlowCtrl is a flow control that jumps to another block.
	GotoFlowCtrl struct {
		KeyPos token.Pos // the position of the keyword, e.g. `goto`, `return`
		Goto   *BlockRef // block to goto
	}

	// CallFlowCtrl is a flow control that calls a function.
	CallFlowCtrl struct {
		KeyPos   token.Pos // the position of the keyword, e.g. `goto`
		Call     Expr      // the expression for the function to call, e.g. `fmt.Println`
		CallArgs []Expr    // the arguments to pass onto the call
		Follow   *BlockRef // block to goto when this call returns
	}

	// RetFlowCtrl is a flow control that returns from a function.
	RetFlowCtrl struct {
		KeyPos token.Pos // the position of the keyword, e.g. `return`
		Follow *BlockRef // block to goto when this receive returns
	}

	// PanicFlowCtrl is a flow control that emits a panic.
	PanicFlowCtrl struct {
		PanicPos token.Pos // the position of the keyword, e.g. `panic`
		Value    Expr      // the value to panic
	}

	// SendFlowCtrl is a flow control that sends a value to a channel.
	SendFlowCtrl struct {
		ArrowPos token.Pos // the position of the channel arrow
		Channel  Expr      // the channel to send a value to
		Value    Expr      // the value to send
		Follow   *BlockRef // block to goto when this call returns
	}

	// ReceiveFlowCtrl is a flow control that receives a value from a channel.
	ReceiveFlowCtrl struct {
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
)

var (
	_ Stmt = (*GotoFlowCtrl)(nil)
	_ Stmt = (*CallFlowCtrl)(nil)
	_ Stmt = (*RetFlowCtrl)(nil)
	_ Stmt = (*PanicFlowCtrl)(nil)
	_ Stmt = (*SendFlowCtrl)(nil)
	_ Stmt = (*ReceiveFlowCtrl)(nil)
)

//====[String]==================================================================

func (s *GotoFlowCtrl) String() string { return fmt.Sprintf(`goto(%v)`, s.Goto) }
func (s *CallFlowCtrl) String() string {
	return fmt.Sprintf(`%s(%s)->%s`, s.Call, csvString(s.CallArgs), s.Follow)
}
func (s *RetFlowCtrl) String() string   { return fmt.Sprintf(`ret(%v)`, csvString(s.Results)) }
func (s *PanicFlowCtrl) String() string { return fmt.Sprintf(`panic(%v)`, s.Value) }
func (s *SendFlowCtrl) String() string {
	return fmt.Sprintf(`(%v<-%v)->%v`, s.Channel, s.Value, s.Follow)
}
func (s *ReceiveFlowCtrl) String() string {
	return fmt.Sprintf(`(%t<-%v)->%v`, s.NeedOk, s.Channel, s.Follow)
}

//====[Pos]=====================================================================

func (s *GotoFlowCtrl) Pos() token.Pos    { return s.KeyPos }
func (s *CallFlowCtrl) Pos() token.Pos    { return s.KeyPos }
func (s *RetFlowCtrl) Pos() token.Pos     { return s.KeyPos }
func (s *PanicFlowCtrl) Pos() token.Pos   { return s.PanicPos }
func (s *SendFlowCtrl) Pos() token.Pos    { return s.ArrowPos }
func (s *ReceiveFlowCtrl) Pos() token.Pos { return s.ArrowPos }

//====[End]=====================================================================

func (s *GotoFlowCtrl) End() token.Pos    { return }
func (s *CallFlowCtrl) End() token.Pos    { return }
func (s *RetFlowCtrl) End() token.Pos     { return }
func (s *PanicFlowCtrl) End() token.Pos   { return }
func (s *SendFlowCtrl) End() token.Pos    { return }
func (s *ReceiveFlowCtrl) End() token.Pos { return }

//====[stmt]====================================================================

func (*GotoFlowCtrl) stmt()    {}
func (*CallFlowCtrl) stmt()    {}
func (*RetFlowCtrl) stmt()     {}
func (*PanicFlowCtrl) stmt()   {}
func (*SendFlowCtrl) stmt()    {}
func (*ReceiveFlowCtrl) stmt() {}

//====[flowCtrl]================================================================

func (*GotoFlowCtrl) flowCtrl()    {}
func (*CallFlowCtrl) flowCtrl()    {}
func (*RetFlowCtrl) flowCtrl()     {}
func (*PanicFlowCtrl) flowCtrl()   {}
func (*SendFlowCtrl) flowCtrl()    {}
func (*ReceiveFlowCtrl) flowCtrl() {}

//==============================================================================

func NewGotoBlock(block *Block) *GotoFlowCtrl {
	return &GotoFlowCtrl{Goto: &BlockRef{Block: block}}
}

func IsFlowCtrl(stmt Stmt) bool {
	_, ok := stmt.(FlowCtrl)
	return ok
}
