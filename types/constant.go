package types

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

// String gets the string for this constant.
func (t *ConstantType) String() string {
	if t == nil {
		return nilStr
	}
	return "const " + ToString(t.Inner)
}
