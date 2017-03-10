package constructs

var _ Construct = (Statement)(nil)

// Statement is an action a program can perform.
type Statement interface {

	// String gets the string for this statement.
	String() string
}
