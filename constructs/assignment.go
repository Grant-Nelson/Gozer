package constructs

import "strings"

var _ Expression = (*AssignmentStat)(nil)

// AssignmentStat is an assignment statement.
type AssignmentStat struct {

	// LeftExp is the left hand expression of the assignment.
	LeftExp []Expression

	// RightExp the right hand expression of the assignment.
	RightExp Expression
}

// Assignment creates a new assignment statement.
func Assignment(left []Expression, right Expression) *AssignmentStat {
	return &AssignmentStat{
		LeftExp:  left,
		RightExp: right,
	}
}

// ReturnTypes are the types this expression resolves to.
func (e *AssignmentStat) ReturnTypes() []Type {
	return e.RightExp.ReturnTypes()
}

// String gets the string for this constuct.
func (e *AssignmentStat) String() string {
	if e == nil {
		return nilStr
	}
	parts := make([]string, len(e.LeftExp))
	for i, exp := range e.LeftExp {
		parts[i] = ToString(exp)
	}
	return strings.Join(parts, ", ") + " = " + ToString(e.RightExp)
}
