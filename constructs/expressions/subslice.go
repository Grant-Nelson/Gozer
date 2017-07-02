package expressions

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

var _ Expression = (*SubsliceExp)(nil)

// SubsliceExp is an expression which gets the index from an expression.
type SubsliceExp struct {

	// Expression is the expression to index from.
	Expression Expression

	// Low the low index for the subslice.
	Low Expression

	// High the high index for the subslice.
	High Expression

	// Maximum the maximum capacity for the subslice.
	Maximum Expression
}

// Subslice creates a new subslice expression.
func Subslice(exp Expression, low Expression, high Expression, max Expression) *SubsliceExp {
	return &SubsliceExp{
		Expression: exp,
		Low:        low,
		High:       high,
		Maximum:    max,
	}
}

// ReturnType is the type this expression resolves to.
func (e *SubsliceExp) ReturnType() types.Type {
	if (e == nil) || (e.Expression == nil) {
		return nil
	}
	return e.Expression.ReturnType()
}

// String gets the string for this constuct.
func (e *SubsliceExp) String() string {
	if e == nil {
		return nilStr
	}
	lowStr := ""
	if e.Low != nil {
		lowStr = ToString(e.Low)
	}
	highStr := ""
	if e.High != nil {
		highStr = ToString(e.High)
	}
	maxStr := ""
	if e.Maximum != nil {
		maxStr = ":" + ToString(e.Maximum)
	}
	return ToString(e.Expression) + "[" + lowStr + ":" + highStr + maxStr + "]"
}
