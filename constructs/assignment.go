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

// Assignment creates a new a assignment statement.
func Assignment(left []Expression, right Expression) *AssignmentStat {
	return &AssignmentStat{
		LeftExp:  left,
		RightExp: right,
	}
}

// ReturnTypes are the types this expression resolves to.
func (s *AssignmentStat) ReturnTypes() []Type {
	return s.RightExp.ReturnTypes()
}

// String gets the string for this constuct.
func (s *AssignmentStat) String() string {
	if s == nil {
		return nilStr
	}
	parts := make([]string, len(s.LeftExp))
	for i, exp := range s.LeftExp {
		parts[i] = ToString(exp)
	}
	return strings.Join(parts, ", ") + " = " + ToString(s.RightExp)
}
