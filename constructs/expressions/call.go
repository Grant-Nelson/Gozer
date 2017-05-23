package expressions

import (
	"strings"

	"github.com/grant-nelson/Gozer/constructs/types"
)

var _ Expression = (*CallExp)(nil)

// CallExp is a function call expression.
type CallExp struct {

	// Function is the function that is being called
	Function *types.FunctionType

	// Receiver is the experssion to get the function target or nil is static.
	Receiver Expression

	// Parameters are the expressions for the parameters of the call.
	Parameters []Expression
}

// Call creates new function call expression.
func Call(function *types.FunctionType, receiver Expression, parameters []Expression) *CallExp {
	return &CallExp{
		Function:   function,
		Receiver:   receiver,
		Parameters: parameters,
	}
}

// ReturnType is the return type from the called function.
func (e *CallExp) ReturnType() types.Type {
	if (e == nil) || (e.Function == nil) {
		return nil
	}
	return e.Function.ReturnType
}

// String gets the string for the call.
func (e *CallExp) String() string {
	if e == nil {
		return nilStr
	}
	parts := make([]string, len(e.Parameters))
	for i, param := range e.Parameters {
		parts[i] = ToString(param)
	}
	return ToString(e.Receiver) + "(" + strings.Join(parts, ", ") + ")"
}
