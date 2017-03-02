package constructs

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

// Equals determins if the given type is the same as this type.
func (t *PointerType) Equals(other Type) bool {
	if t == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	tother, ok := other.(*PointerType)
	if !ok {
		return false
	}
	if !t.Inner.Equals(tother.Inner) {
		return false
	}
	return true
}

// String gets the name for this type.
func (t *PointerType) String() string {
	if t == nil {
		return nilStr
	}
	return "*" + t.Inner.String()
}
