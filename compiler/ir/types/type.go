package types

// Type represents a description of type information.
type Type interface {

	// Underlying return the underlying type of this type.
	//
	// See https://go.dev/ref/spec#Underlying_types.
	Underlying() Type

	// String gets a string representation of this type.
	String() string
}
