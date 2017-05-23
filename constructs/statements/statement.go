package statements

import (
	"github.com/grant-nelson/Gozer/constructs/expressions"
)

var _ Statement = (expressions.Expression)(nil)

// nilStr is the string to use for nil values.
const nilStr = "nil"

// Statement is an action a program can perform.
type Statement interface {

	// String gets the string for this statement.
	String() string
}

// ToString creates a string for the given construct.
func ToString(s Statement) string {
	if s == nil {
		return nilStr
	}
	return s.String()
}
