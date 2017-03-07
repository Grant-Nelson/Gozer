package constructs

var _ Construct = (Type)(nil)

// Type is the interface for all types.
type Type interface {

	// Equals determines if this type is the same as the other type.
	Equals(other interface{}) bool

	// String gets the identifier name for this type.
	String() string
}
