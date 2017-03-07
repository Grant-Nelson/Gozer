package constructs

var _ Expression = (*LiteralExp)(nil)

// LiteralExp is a constant typed value.
type LiteralExp struct {

	// Value is the string representation of the literal.
	Value string

	// Type is the value type of the literal.
	Type Type
}

// Literal creates new literal expression.
func Literal(value string, litType Type) *LiteralExp {
	return &LiteralExp{
		Value: value,
		Type:  litType,
	}
}

// ReturnTypes is an array containing the type of the literal.
func (e *LiteralExp) ReturnTypes() []Type {
	return []Type{e.Type}
}

// Equals determines if this literal is the same as the other expression.
func (e *LiteralExp) Equals(other interface{}) bool {
	if e == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	eother, ok := other.(*LiteralExp)
	if !ok {
		return false
	}
	if e.Value != eother.Value {
		return false
	}
	if !Equals(e.Type, eother.Type) {
		return false
	}
	return true
}

// String gets the string for the literal.
func (e *LiteralExp) String() string {
	return ToString(e.Type) + "(" + e.Value + ")"
}
