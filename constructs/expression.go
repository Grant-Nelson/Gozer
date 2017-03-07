package constructs

var _ Statement = (Expression)(nil)

// Expression is a constuct for defining a combination of one
// or more literals, variables, operators, and functions.
type Expression interface {

	// ReturnTypes are the typea this expression resolves to.
	ReturnTypes() []Type

	// Equals determines if this expression is the same as the other expression.
	Equals(other interface{}) bool

	// String gets the string for this constuct.
	String() string
}
