package expressions

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

var _ Expression = (*AssignmentExp)(nil)

// AssignmentExp is an assignment expression.
type AssignmentExp struct {

	// LeftExp is the left hand expression of the assignment.
	LeftExp Expression

	// RightExp the right hand expression of the assignment.
	RightExp Expression
}

// Assignment creates a new assignment statement.
func Assignment(left Expression, right Expression) *AssignmentExp {
	return &AssignmentExp{
		LeftExp:  left,
		RightExp: right,
	}
}

// ReturnType is the type this expression resolves to.
func (e *AssignmentExp) ReturnType() types.Type {
	if (e == nil) || (e.RightExp == nil) {
		return nil
	}
	return e.RightExp.ReturnType()
}

// String gets the string for this constuct.
func (e *AssignmentExp) String() string {
	if e == nil {
		return nilStr
	}
	return ToString(e.LeftExp) + " = " + ToString(e.RightExp)
}
