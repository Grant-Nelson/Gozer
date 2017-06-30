package expressions

import (
	"github.com/grant-nelson/Gozer/constructs/types"
)

var _ Expression = (*LiteralExp)(nil)

// LambdaExp is a constant function definition.
type LambdaExp struct {

	// Function is the function literal.
	Function *types.FunctionType
}

// Lambda creates new lambda expression.
func Lambda(function *types.FunctionType) *LambdaExp {
	return &LambdaExp{
		Function: function,
	}
}

// ReturnType is the type of the lambda.
func (e *LambdaExp) ReturnType() types.Type {
	if e == nil {
		return nil
	}
	return e.Function
}

// String gets the string for the lambda.
func (e *LambdaExp) String() string {
	if e == nil {
		return nilStr
	}
	return e.Function.FullBodyString()
}
