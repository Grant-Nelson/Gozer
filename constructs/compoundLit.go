package constructs

import (
	"strings"
)

var _ Expression = (*CompoundLiteralExp)(nil)

// CompoundLiteralExp is a list definition typed value.
type CompoundLiteralExp struct {

	// Elements is the set of values to initialize this list with.
	Elements []Expression

	// Type is the value type of the literal.
	Type Type
}

// CompoundLiteral creates new composite literal expression.
func CompoundLiteral(elems []Expression, litType Type) *CompoundLiteralExp {
	return &CompoundLiteralExp{
		Elements: elems,
		Type:     litType,
	}
}

// ReturnType is the type of the composite literal.
func (e *CompoundLiteralExp) ReturnType() Type {
	return e.Type
}

// String gets the string for the composite literal.
func (e *CompoundLiteralExp) String() string {
	if e == nil {
		return ""
	}
	parts := make([]string, len(e.Elements))
	for i, elem := range e.Elements {
		parts[i] = ToString(elem)
	}
	return ToString(e.Type) + "{" + strings.Join(parts, ", ") + "}"
}
