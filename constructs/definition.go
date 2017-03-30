package constructs

var _ Statement = (*DefinitionExp)(nil)

// DefinitionExp is an expression which represents
// the creation of a variable, constant, or some other named construct.
type DefinitionExp struct {

	// Identifier the name of the definition.
	Identifier string

	// Type is the resulting type of the definition.
	Type Type

	// InitialValue is the expression for the initial value of this definition.
	// If nil use the define value for the variable.
	InitialValue Expression
}

// Definition creates a new identifier definition.
func Definition(id string, t Type, exp Expression) *DefinitionExp {
	return &DefinitionExp{
		Identifier:   id,
		Type:         t,
		InitialValue: exp,
	}
}

// ReturnTypes are the types this expression resolves to.
func (e *DefinitionExp) ReturnTypes() []Type {
	return []Type{e.Type}
}

// String gets the string for this constuct.
func (e *DefinitionExp) String() string {
	if e == nil {
		return nilStr
	}
	if e.InitialValue == nil {
		return ToString(e.Type) + " " + e.Identifier
	}
	return ToString(e.Type) + " " + e.Identifier + " = " + ToString(e.InitialValue)
}
