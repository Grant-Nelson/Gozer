package irc

import (
	"fmt"
	"go/token"
	"go/types"

	"github.com/Grant-Nelson/Gozer/project/enums/binaryOp"
	"github.com/Grant-Nelson/Gozer/project/enums/unaryOp"
)

type (
	Expr interface {
		Node

		// ResultType is the type this expression results in.
		ResultType() types.Type

		// exprNode is an empty method used to compile time type check that
		// only expressions duck-type to this interface.
		exprNode()
	}

	BasicLit struct {
		ValuePos token.Pos    // the position of the literal
		ValueEnd token.Pos    // position immediately after the literal
		Value    string       // the value of the constant
		Type     *types.Basic // the basic type of this value
	}

	// CallExpr is a method call that is non-breaking.
	CallExpr struct {
		Func      Expr       // expression for the function being called
		LeftParen token.Pos  // position of the left parenthesis of this call
		Args      []Expr     // function arguments in the call or empty
		Result    types.Type // the result type from this call
	}

	// UnaryExpr is a unary operation applied to an expression.
	UnaryExpr struct {
		OpPos  token.Pos       // position of the unary operator
		Op     unaryOp.UnaryOp // the unary operator kind
		Expr   Expr            // the operand
		Result types.Type      // the result type of the unary expr
	}

	BinaryExpr struct {
		Left   Expr              // left operand in the operation
		OpPos  token.Pos         // position of the binary operator
		Op     binaryOp.BinaryOp // the binary operator kind
		Right  Expr              // right operand in the operation
		Result types.Type        // the result type of the binary operator
	}

	TernaryExpr struct {
		OpPos  token.Pos  // the location of the `?` or the original if-statement
		Cond   Expr       // the conditional expression, must be a bool type
		Left   Expr       // the value to use when the conditional is true
		Right  Expr       // the value to use when the conditional is false
		Result types.Type // the result type, should match the left and right types
	}
)

var (
	_ Expr = (*BasicLit)(nil)
	_ Expr = (*CallExpr)(nil)
	_ Expr = (*UnaryExpr)(nil)
	_ Expr = (*BinaryExpr)(nil)
	_ Expr = (*TernaryExpr)(nil)
)

//====[String]==================================================================

func (e *BasicLit) String() string    { return fmt.Sprintf(`%v(%v)`, e.Type, e.Value) }
func (e *CallExpr) String() string    { return fmt.Sprintf(`%v(%s)`, e.Func, csvString(e.Args)) }
func (e *UnaryExpr) String() string   { return fmt.Sprintf(e.Op.Format(), e.Expr) }
func (e *BinaryExpr) String() string  { return fmt.Sprintf(e.Op.Format(), e.Left, e.Right) }
func (e *TernaryExpr) String() string { return fmt.Sprintf(`(%v?%v:%v)`, e.Cond, e.Left, e.Right) }

//====[Pos]=====================================================================

func (e *BasicLit) Pos() token.Pos    { return e.ValuePos }
func (e *CallExpr) Pos() token.Pos    { return e.LeftParen }
func (e *UnaryExpr) Pos() token.Pos   { return e.OpPos }
func (e *BinaryExpr) Pos() token.Pos  { return e.OpPos }
func (e *TernaryExpr) Pos() token.Pos { return e.OpPos }

//====[End]=====================================================================

func (e *BasicLit) End() token.Pos    { return e.ValueEnd }
func (e *CallExpr) End() token.Pos    { return }
func (e *UnaryExpr) End() token.Pos   { return }
func (e *BinaryExpr) End() token.Pos  { return }
func (e *TernaryExpr) End() token.Pos { return }

//====[ResultType]==============================================================

func (e *BasicLit) ResultType() types.Type    { return e.Type }
func (e *CallExpr) ResultType() types.Type    { return e.Result }
func (e *UnaryExpr) ResultType() types.Type   { return e.Result }
func (e *BinaryExpr) ResultType() types.Type  { return e.Result }
func (e *TernaryExpr) ResultType() types.Type { return e.Result }

//====[exprNode]================================================================

func (*BasicLit) exprNode()    {}
func (*CallExpr) exprNode()    {}
func (*UnaryExpr) exprNode()   {}
func (*BinaryExpr) exprNode()  {}
func (*TernaryExpr) exprNode() {}
