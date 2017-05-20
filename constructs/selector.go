package constructs

var _ Expression = (*SelectorExp)(nil)

// SelectorExp is an expression which gets the member
// or function handle from an expression.
type SelectorExp struct {

	// Expression is the expression to select from.
	Expression Expression

	// Identifier the name of the selection.
	Identifier string

	// Type is the resulting type of the selection.
	Type Type
}

// Selector creates a new selector expression.
func Selector(exp Expression, id string, t Type) *SelectorExp {
	return &SelectorExp{
		Expression: exp,
		Identifier: id,
		Type:       t,
	}
}

// ReturnType is the type this expression resolves to.
func (e *SelectorExp) ReturnType() Type {
	return e.Type
}

// String gets the string for this constuct.
func (e *SelectorExp) String() string {
	if e == nil {
		return nilStr
	}
	return ToString(e.Expression) + "." + e.Identifier
}
