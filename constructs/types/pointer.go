package types

var _ Type = (*PointerType)(nil)

// PointerType for storing the types of pointers.
type PointerType struct {

	// Inner gets the inner type for this pointer.
	Inner Type
}

// Pointer creates a new pointer type for the given type.
func Pointer(inner Type) *PointerType {
	return &PointerType{
		Inner: inner,
	}
}

// String gets the name for this type.
func (t *PointerType) String() string {
	if t == nil {
		return nilStr
	}
	return "*" + ToString(t.Inner)
}
