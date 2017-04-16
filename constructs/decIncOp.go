package constructs

var _ Statement = (*DecIncOpExp)(nil)

// DecIncOpExp is an expression which increments or decrements an expression.
type DecIncOpExp struct {

	// Left is the left hand expression in the operation.
	Exp Expression

	// Increment indicates if this operation is an increment or decrement.
	Increment bool
}

// IncDec creates a new expression increment or decrement operator.
func IncDec(exp Expression, inc bool) *DecIncOpExp {
	return &DecIncOpExp{
		Exp:       exp,
		Increment: inc,
	}
}

// String gets the string for this constuct.
func (s *DecIncOpExp) String() string {
	if s == nil {
		return nilStr
	}
	if s.Increment {
		return ToString(s.Exp) + "++"
	}
	return ToString(s.Exp) + "--"
}
