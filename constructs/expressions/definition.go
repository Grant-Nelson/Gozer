package expressions

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

var _ Expression = (*DefinitionExp)(nil)

// DefinitionExp is an expression which represents
// the creation of a variable, constant, or some other named construct.
type DefinitionExp struct {

	// Identifier the name and type of the definition.
	Identifier *IdentifierExp

	// InitialValue is the expression for the initial value of this definition.
	// If nil use the define value for the variable.
	InitialValue Expression
}

// Definition creates a new identifier definition.
func Definition(id *IdentifierExp, exp Expression) *DefinitionExp {
	return &DefinitionExp{
		Identifier:   id,
		InitialValue: exp,
	}
}

// ReturnType is the type this expression resolves to.
func (e *DefinitionExp) ReturnType() types.Type {
	if (e == nil) || (e.Identifier == nil) {
		return nil
	}
	return e.Identifier.ReturnType()
}

// String gets the string for this constuct.
func (e *DefinitionExp) String() string {
	if e == nil {
		return nilStr
	}
	result := nilStr
	if e.Identifier != nil {
		result = types.ToString(e.Identifier.Type) + " " + e.Identifier.Identifier
	}
	if e.InitialValue != nil {
		result += " = " + ToString(e.InitialValue)
	}
	return result
}
