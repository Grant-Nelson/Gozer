package irc

import (
	"fmt"
	"go/token"

	"github.com/Grant-Nelson/Gozer/project/enums/binaryOp"
	"github.com/Grant-Nelson/Gozer/project/enums/unaryOp"
)

type (
	Expr interface {
		// String gets a simple representation for this expression.
		String() string

		// Pos is the position for this expression in the source file.
		// This may be [token.NoPos] if not specifically in the source.
		//
		// This may be able to be used to look up the [ast.Node] that
		// became this IRC statement, however, this may not match the
		// [Pos()] for the [ast.Expr]. For example, the binary operator
		// expression will return the position of the operator symbol
		// instead of the position of the left operand.
		Pos() token.Pos

		// ResultType is the type this expression results in.
		ResultType() Type

		// expr is an empty method used to compile time type check that
		// only expressions are used for this interface.
		expr()
	}

	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
		Type    Type      // the type of the identifier
	}

	BasicLit struct {
		ValuePos token.Pos  // the position of the literal
		Value    string     // the value of te constant
		Type     *BasicType // the basic type of this value
	}

	// CallExpr is a method call that is non-breaking.
	CallExpr struct {
		Func      Expr      // expression for the function being called
		LeftParen token.Pos // position of the left parenthesis of this call
		Args      []Expr    // function arguments in the call or empty
		Result    Type      // the result type from this call
	}

	// UnaryExpr is a unary operation applied to an expression.
	UnaryExpr struct {
		OpPos  token.Pos       // position of the unary operator
		Op     unaryOp.UnaryOp // the unary operator kind
		Expr   Expr            // the operand
		Result Type            // the result type of the unary expr
	}

	BinaryExpr struct {
		Left   Expr              // left operand in the operation
		OpPos  token.Pos         // position of the binary operator
		Op     binaryOp.BinaryOp // the binary operator kind
		Right  Expr              // right operand in the operation
		Result Type              // the result type of the binary operator
	}

	TernaryExpr struct {
		OpPos  token.Pos // the location of the `?` or the original if-statement
		Cond   Expr      // the conditional expression, must be a bool type
		Left   Expr      // the value to use when the conditional is true
		Right  Expr      // the value to use when the conditional is false
		Result Type      // the result type, should match the left and right types
	}

	TupleExpr struct {
		OpenPos token.Pos  // the location of the `(` for the group
		Type    *TupleType // the result type of this tuple
		Values  []Expr     // the values in the tuple
	}
)

var (
	_ Expr = (*Ident)(nil)
	_ Expr = (*BasicLit)(nil)
	_ Expr = (*CallExpr)(nil)
	_ Expr = (*UnaryExpr)(nil)
	_ Expr = (*BinaryExpr)(nil)
	_ Expr = (*TernaryExpr)(nil)
	_ Expr = (*TupleExpr)(nil)
)

// IsExported reports whether id starts with an upper-case letter.
func (id *Ident) IsExported() bool { return token.IsExported(id.Name) }

func (id *Ident) String() string {
	if id != nil {
		return id.Name
	}
	return `<nil>`
}
func (e *BasicLit) String() string    { return fmt.Sprintf(`%s(%s)`, e.Type, e.Value) }
func (e *CallExpr) String() string    { return fmt.Sprintf(`%s(%s)`, e.Func, csvString(e.Args)) }
func (e *UnaryExpr) String() string   { return fmt.Sprintf(e.Op.Format(), e.Expr) }
func (e *BinaryExpr) String() string  { return fmt.Sprintf(e.Op.Format(), e.Left, e.Right) }
func (e *TernaryExpr) String() string { return fmt.Sprintf(`(%s?%s:%s)`, e.Cond, e.Left, e.Right) }
func (e *TupleExpr) String() string   { return fmt.Sprintf(`(%s)`, csvString(e.Values)) }

func (e *Ident) Pos() token.Pos       { return e.NamePos }
func (e *BasicLit) Pos() token.Pos    { return e.ValuePos }
func (e *CallExpr) Pos() token.Pos    { return e.LeftParen }
func (e *UnaryExpr) Pos() token.Pos   { return e.OpPos }
func (e *BinaryExpr) Pos() token.Pos  { return e.OpPos }
func (e *TernaryExpr) Pos() token.Pos { return e.OpPos }
func (e *TupleExpr) Pos() token.Pos   { return e.OpenPos }

func (e *Ident) ResultType() Type       { return e.Type }
func (e *BasicLit) ResultType() Type    { return e.Type }
func (e *CallExpr) ResultType() Type    { return e.Result }
func (e *UnaryExpr) ResultType() Type   { return e.Result }
func (e *BinaryExpr) ResultType() Type  { return e.Result }
func (e *TernaryExpr) ResultType() Type { return e.Result }
func (e *TupleExpr) ResultType() Type   { return e.Type }

func (*Ident) expr()       {}
func (*BasicLit) expr()    {}
func (*CallExpr) expr()    {}
func (*UnaryExpr) expr()   {}
func (*BinaryExpr) expr()  {}
func (*TernaryExpr) expr() {}
func (*TupleExpr) expr()   {}
