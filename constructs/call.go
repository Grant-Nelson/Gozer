package constructs

import "strings"

var _ Expression = (*CallExp)(nil)

// CallExp is a function call expression.
type CallExp struct {

	// Function is the function that is being called
	Function *FunctionType

	// Receiver is the experssion to get the function target or nil is static.
	Receiver Expression

	// Parameters are the expressions for the parameters of the call.
	Parameters []Expression
}

// Call creates new function call expression.
func Call(function *FunctionType, receiver Expression, parameters []Expression) *CallExp {
	return &CallExp{
		Function:   function,
		Receiver:   receiver,
		Parameters: parameters,
	}
}

// ReturnTypes is the list of return types from the called function.
func (e *CallExp) ReturnTypes() []Type {
	if e.Function == nil {
		return []Type{}
	}
	return e.Function.ReturnTypes
}

// Equals determines if this call is the same as the other expression.
func (e *CallExp) Equals(other interface{}) bool {
	if e == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	eother, ok := other.(*CallExp)
	if !ok {
		return false
	}
	if len(e.Parameters) != len(eother.Parameters) {
		return false
	}
	if !Equals(e.Function, eother.Function) {
		return false
	}
	if !Equals(e.Receiver, eother.Receiver) {
		return false
	}
	for i, param := range e.Parameters {
		if !Equals(param, eother.Parameters[i]) {
			return false
		}
	}
	return true
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
	return "(" + ToString(e.Receiver) + ")" + ToString(e.Function) + "(" + strings.Join(parts, ", ") + ")"
}
