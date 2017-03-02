package constructs

import "strings"

const nilStr = "nil"

// Type is the interface for all types.
type Type interface {

	// Equals determines if this type is the same as the other type.
	Equals(other Type) bool

	// String gets the identifier name for this type.
	String() string
}

// indent returns the given text indented.
func indent(text string, indent string) string {
	return indent + strings.Replace(text, "\n", "\n"+indent, -1)
}
