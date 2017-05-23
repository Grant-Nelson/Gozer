package expressions

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

var _ Expression = (*IdentifierExp)(nil)

// IdentifierExp is an expression which represents
// a variable, constant, or some other named construct.
type IdentifierExp struct {

	// Identifier the name of the identifier.
	Identifier string

	// Type is the resulting type of the identifier.
	Type types.Type
}

// Identifier creates a new identifier expression.
func Identifier(id string, t types.Type) *IdentifierExp {
	return &IdentifierExp{
		Identifier: id,
		Type:       t,
	}
}

// ReturnType is the type this expression resolves to.
func (e *IdentifierExp) ReturnType() types.Type {
	if e == nil {
		return nil
	}
	return e.Type
}

// String gets the string for this constuct.
func (e *IdentifierExp) String() string {
	if e == nil {
		return nilStr
	}
	return e.Identifier
}
