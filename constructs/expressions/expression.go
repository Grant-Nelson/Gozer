package expressions

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

// nilStr is the string to use for nil values.
const nilStr = "nil"

// Expression is a constuct for defining a combination of one
// or more literals, variables, operators, and functions.
type Expression interface {

	// ReturnType is the type this expression resolves to.
	ReturnType() types.Type

	// String gets the string for this constuct.
	String() string
}

// ToString creates a string for the given expression.
func ToString(e Expression) string {
	if e == nil {
		return nilStr
	}
	return e.String()
}
