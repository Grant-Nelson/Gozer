package expressions

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

var _ Expression = (*BinaryOpExp)(nil)

const (
	// AddOp is the binary operation for adding two values.
	AddOp = "+"

	// SubtractOp is the binary operation for subtracting two values.
	SubtractOp = "-"

	// MultiplyOp is the binary operation for multiplying two values.
	MultiplyOp = "*"

	// QuotentOp is the binary operation for dividing two values.
	QuotentOp = "/"

	// RemainderOp is the binary operation for getting the
	// remainder from the division of two values.
	RemainderOp = "%"

	// AndOp is the binary operation for AND of two values.
	AndOp = "&"

	// OrOp is the binary operation for OR of two values.
	OrOp = "|"

	// ExclusiveOrOp is the binary operation for XOR of two values.
	ExclusiveOrOp = "^"

	// LeftShiftOp is the binary operation for a value right shifted by another.
	LeftShiftOp = "<<"

	// RightShiftOp is the binary operation for a value left shifted by another.
	RightShiftOp = ">>"

	// AndNotOp is the binary operation for AND NOT of two values.
	AndNotOp = "&^"

	// LogicalAndOp is the binary operation for logical AND of two values.
	LogicalAndOp = "&&"

	// LogicalOrOp is the binary operation for logical OR of two values.
	LogicalOrOp = "||"

	// EqualOp is the binary operation for equality of two values.
	EqualOp = "=="

	// NotEqualOp is the binary operation for a value not equal to another.
	NotEqualOp = "!="

	// LessThanOp is the binary operation for a value is less than another.
	LessThanOp = "<"

	// LessThanEqualOp is the binary operation for a value is less than or equal to another.
	LessThanEqualOp = "<="

	// GreaterThanOp is the binary operation for a value is greater than another.
	GreaterThanOp = ">"

	// GreaterThanEqualOp is the binary operation for a value is greater than or equal to another.
	GreaterThanEqualOp = ">="
)

// BinaryOpExp is an expression which combines two values with some operation.
type BinaryOpExp struct {

	// Left is the left hand expression in the operation.
	Left Expression

	// Right is the right hand expression in the operation.
	Right Expression

	// Operand is the operations to apply to the two expressions.
	Operand string

	// Type is the resulting type after the operation.
	Type types.Type
}

// BinaryOp creates a new binary operation expression.
func BinaryOp(left Expression, right Expression, operand string, t types.Type) *BinaryOpExp {
	return &BinaryOpExp{
		Left:    left,
		Right:   right,
		Operand: operand,
		Type:    t,
	}
}

// ReturnType is the type this expression resolves to.
func (e *BinaryOpExp) ReturnType() types.Type {
	if e == nil {
		return nil
	}
	return e.Type
}

// String gets the string for this constuct.
func (e *BinaryOpExp) String() string {
	if e == nil {
		return nilStr
	}
	return "(" + ToString(e.Left) + " " + e.Operand + " " + ToString(e.Right) + ")"
}
