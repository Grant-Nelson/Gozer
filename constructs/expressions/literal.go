package expressions

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

var _ Expression = (*LiteralExp)(nil)

// LiteralExp is a constant typed value.
type LiteralExp struct {

	// Value is the string representation of the literal.
	Value string

	// Type is the value type of the literal.
	Type types.Type
}

// Literal creates new literal expression.
func Literal(value string, litType types.Type) *LiteralExp {
	return &LiteralExp{
		Value: value,
		Type:  litType,
	}
}

// ReturnType is the type of the literal.
func (e *LiteralExp) ReturnType() types.Type {
	if e == nil {
		return nil
	}
	return e.Type
}

// String gets the string for the literal.
func (e *LiteralExp) String() string {
	if e == nil {
		return nilStr
	}
	return e.Value
}
