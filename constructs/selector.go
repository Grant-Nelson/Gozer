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
func Selector(exp Expression, id string) *SelectorExp {
	return &SelectorExp{
		Expression: exp,
		Identifier: id,
		Type:       nil,
	}
}

// ReturnTypes are the typea this expression resolves to.
func (e *SelectorExp) ReturnTypes() []Type {
	return []Type{e.Type}
}

// Equals determines if this expression is the same as the other expression.
func (e *SelectorExp) Equals(other interface{}) bool {
	if e == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	eother, ok := other.(*SelectorExp)
	if !ok {
		return false
	}
	if !Equals(e.Expression, eother.Expression) {
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
func (e *SelectorExp) String() string {
	if e == nil {
		return nilStr
	}
	return ToString(e.Expression) + "." + e.Identifier + "(" + ToString(e.Type) + ")"
}
