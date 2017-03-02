package constructs

var _ Type = (*ConstantType)(nil)

// ConstantType contains the information for constant types.
type ConstantType struct {

	// Inner gets the inner type for this contant.
	Inner Type
}

// Constant creates a new constant type.
func Constant(inner Type) *ConstantType {
	return &ConstantType{
		Inner: inner,
	}
}

// Equals determins if the given type is the same as this type.
func (t *ConstantType) Equals(other Type) bool {
	if t == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	tother, ok := other.(*ConstantType)
	if !ok {
		return false
	}
	if !t.Inner.Equals(tother.Inner) {
		return false
	}
	return true
}

// String gets the string for this constant.
func (t *ConstantType) String() string {
	if t == nil {
		return nilStr
	}
	return "const " + t.Inner.String()
}
