package constructs

var _ Construct = (Type)(nil)

// Type is the interface for all types.
type Type interface {

	// String gets the identifier name for this type.
	String() string
}
