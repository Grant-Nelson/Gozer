package types

var _ Type = (NamedType)(nil)

// NamedType is a type which has named definition,
// such as interfaces, functions, or structures
type NamedType interface {

	// GetName gets the name of the type.
	// May be empty if this type is unnamed.
	GetName() string

	// String gets the identifier name for this type.
	// If the name isn't set then the type is considered
	// unnamed and FullString will be returned.
	String() string

	// Gets the full definition of this type,
	// for example the signature and body of a function.
	FullString() string
}
