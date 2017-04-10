package constructs

import "fmt"

var _ Statement = (*IfStat)(nil)

// IfStat is an if statement.
type IfStat struct {

	// Cond is the expession for the conditional of this if.
	Cond Expression

	// Body is the statement to call when Cond evaluates to true.
	Body Statement

	// Else is the statement to call when Cond evaluates to false.
	// Maybe nil if there is no else part.
	Else Statement
}

// If creates a new if-statement.
func If(cond Expression, bodyStat Statement, elseStat Statement) *IfStat {
	return &IfStat{
		Cond: cond,
		Body: bodyStat,
		Else: elseStat,
	}
}

// String gets the string for this constuct.
func (s *IfStat) String() string {
	if s == nil {
		return nilStr
	}
	result := fmt.Sprint("if ", ToString(s.Cond), " ", ToString(s.Body))
	if s.Else != nil {
		result = fmt.Sprint(result, " else ", ToString(s.Else))
	}
	return result
}
