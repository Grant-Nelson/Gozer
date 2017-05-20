package constructs

var _ Statement = (*DefinitionExp)(nil)

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
func (s *DefinitionExp) ReturnType() Type {
	return s.Identifier.ReturnType()
}

// String gets the string for this constuct.
func (s *DefinitionExp) String() string {
	if s == nil {
		return nilStr
	}
	result := nilStr
	if s.Identifier != nil {
		result = ToString(s.Identifier.Type) + " " + s.Identifier.Identifier
	}
	if s.InitialValue != nil {
		result += " = " + ToString(s.InitialValue)
	}
	return result
}
