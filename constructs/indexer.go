package constructs

var _ Expression = (*IndexerExp)(nil)

// IndexerExp is an expression which gets the index from an expression.
type IndexerExp struct {

	// Expression is the expression to index from.
	Expression Expression

	// IndexExp the name of the index.
	IndexExp Expression

	// Type is the resulting type of the index.
	Type Type
}

// Indexer creates a new index expression.
func Indexer(exp Expression, index Expression, t Type) *IndexerExp {
	return &IndexerExp{
		Expression: exp,
		IndexExp:   index,
		Type:       t,
	}
}

// ReturnTypes are the types this expression resolves to.
func (e *IndexerExp) ReturnTypes() []Type {
	return []Type{e.Type}
}

// String gets the string for this constuct.
func (e *IndexerExp) String() string {
	if e == nil {
		return nilStr
	}
	return ToString(e.Expression) + "[" + ToString(e.IndexExp) + "]"
}
