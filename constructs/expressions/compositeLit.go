package expressions

import (
	"strings"

	"github.com/grant-nelson/Gozer/constructs/types"
)

var _ Expression = (*CompositeLiteralExp)(nil)

// CompositeLiteralExp is a list or map definition typed value.
type CompositeLiteralExp struct {

	// Elements is the set of values to initialize.
	Elements []Expression

	// Type is the value type of the literal.
	Type types.Type
}

// CompositeLiteral creates new composite literal expression.
func CompositeLiteral(elems []Expression, litType types.Type) *CompositeLiteralExp {
	return &CompositeLiteralExp{
		Elements: elems,
		Type:     litType,
	}
}

// ReturnType is the type of the composite literal.
func (e *CompositeLiteralExp) ReturnType() types.Type {
	if e == nil {
		return nil
	}
	return e.Type
}

// String gets the string for the composite literal.
func (e *CompositeLiteralExp) String() string {
	if e == nil {
		return nilStr
	}
	parts := make([]string, len(e.Elements))
	for i, elem := range e.Elements {
		parts[i] = ToString(elem)
	}
	return types.ToString(e.Type) + "{" + strings.Join(parts, ", ") + "}"
}
