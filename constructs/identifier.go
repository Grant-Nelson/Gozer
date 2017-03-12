package constructs

var _ Expression = (*IdentifierExp)(nil)

// IdentifierExp is an expression which represents
// a variable, constant, or some other named construct.
type IdentifierExp struct {

	// Identifier the name of the identifier.
	Identifier string

	// Type is the resulting type of the identifier.
	Type Type
}

// Identifier creates a new identifier expression.
func Identifier(id string, t Type) *IdentifierExp {
	return &IdentifierExp{
		Identifier: id,
		Type:       t,
	}
}

// ReturnTypes are the typea this expression resolves to.
func (e *IdentifierExp) ReturnTypes() []Type {
	return []Type{e.Type}
}

// String gets the string for this constuct.
func (e *IdentifierExp) String() string {
	if e == nil {
		return nilStr
	}
	return e.Identifier
}
