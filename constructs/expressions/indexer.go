package expressions

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

var _ Expression = (*IndexerExp)(nil)

// IndexerExp is an expression which gets the index from an expression.
type IndexerExp struct {

	// Expression is the expression to index from.
	Expression Expression

	// IndexExp the name of the index.
	IndexExp Expression
}

// Indexer creates a new index expression.
func Indexer(exp Expression, index Expression) *IndexerExp {
	return &IndexerExp{
		Expression: exp,
		IndexExp:   index,
	}
}

// ReturnType is the type this expression resolves to.
func (e *IndexerExp) ReturnType() types.Type {
	if (e == nil) || (e.Expression == nil) {
		return nil
	}
	t, _ := types.GetIndexableType(e.Expression.ReturnType())
	return t
}

// String gets the string for this constuct.
func (e *IndexerExp) String() string {
	if e == nil {
		return nilStr
	}
	return ToString(e.Expression) + "[" + ToString(e.IndexExp) + "]"
}
