package constructs

var _ Construct = (Statement)(nil)

// Statement is an action a program can perform.
type Statement interface {

	// Equals determines if this statement is the same as the other statement.
	Equals(other interface{}) bool

	// String gets the string for this statement.
	String() string
}
