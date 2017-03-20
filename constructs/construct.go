package constructs

const nilStr = "nil"

// Construct is the base interface for all objects representing parts of codes.
type Construct interface {

	// String gets the identifier name for this constructs.
	String() string
}

// ToString creates a string for the given construct.
func ToString(c Construct) string {
	if c == nil {
		return nilStr
	}
	return c.String()
}
