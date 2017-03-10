package constructs

var _ Statement = (Expression)(nil)

// Expression is a constuct for defining a combination of one
// or more literals, variables, operators, and functions.
type Expression interface {

	// ReturnTypes are the typea this expression resolves to.
	ReturnTypes() []Type

	// String gets the string for this constuct.
	String() string
}
