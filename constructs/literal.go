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

// String gets the string for the literal.
func (e *LiteralExp) String() string {
	return e.Value
}
