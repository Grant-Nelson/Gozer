package types

// nilStr is the string to use for nil values.
const nilStr = "nil"

// Type is the interface for all types.
type Type interface {

	// String gets the identifier name for this type.
	String() string
}

// ToString creates a string for the given construct.
func ToString(t Type) string {
	if t == nil {
		return nilStr
	}
	return t.String()
}
