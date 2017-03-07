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

// SelectIdentifieror creates a new identifier expression.
func Identifier(id string) *IdentifierExp {
	return &IdentifierExp{
		Identifier: id,
		Type:       nil,
	}
}

// ReturnTypes are the typea this expression resolves to.
func (e *IdentifierExp) ReturnTypes() []Type {
	return []Type{e.Type}
}

// Equals determines if this expression is the same as the other expression.
func (e *IdentifierExp) Equals(other interface{}) bool {
	if e == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	eother, ok := other.(*SelectorExp)
	if !ok {
		return false
	}
	if e.Identifier != eother.Identifier {
		return false
	}
	if !Equals(e.Type, eother.Type) {
		return false
	}
	return true
}

// String gets the string for this constuct.
func (e *IdentifierExp) String() string {
	if e == nil {
		return nilStr
	}
	return e.Identifier + "(" + ToString(e.Type) + ")"
}
