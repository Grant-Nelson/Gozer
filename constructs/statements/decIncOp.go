package statements

import (
	"github.com/grant-nelson/Gozer/constructs/expressions"
)

var _ Statement = (*DecIncOpStat)(nil)

// DecIncOpStat is an expression which increments or decrements an expression.
type DecIncOpStat struct {

	// Exp is the expression in the operation.
	Exp expressions.Expression

	// Increment indicates if this operation is an increment or decrement.
	Increment bool
}

// IncDecOp creates a new expression increment or decrement operator.
func IncDecOp(exp expressions.Expression, inc bool) *DecIncOpStat {
	return &DecIncOpStat{
		Exp:       exp,
		Increment: inc,
	}
}

// String gets the string for this constuct.
func (s *DecIncOpStat) String() string {
	if s == nil {
		return nilStr
	}
	if s.Increment {
		return expressions.ToString(s.Exp) + "++"
	}
	return expressions.ToString(s.Exp) + "--"
}
