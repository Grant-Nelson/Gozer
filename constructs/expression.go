package constructs

var _ Statement = (Expression)(nil)

// Expression is a constuct for defining a combination of one
// or more literals, variables, operators, and functions.
type Expression interface {

	// ReturnType is the type this expression resolves to.
	ReturnType() Type

	// String gets the string for this constuct.
	String() string
}
