package constructs

import "strings"

const nilStr = "nil"

// Construct is the base interface for all objects representing parts of codes.
type Construct interface {

	// Equals determines if this type is the same as the other constructs.
	Equals(other interface{}) bool

	// String gets the identifier name for this constructs.
	String() string
}

// indent returns the given text indented.
func indent(text string, indent string) string {
	return indent + strings.Replace(text, "\n", "\n"+indent, -1)
}

// Equals checks the equality of the two given constructs.
func Equals(c Construct, other interface{}) bool {
	if c == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	return c.Equals(other)
}

// ToString creates a string for the given construct.
func ToString(c Construct) string {
	if c == nil {
		return nilStr
	}
	return c.String()
}
