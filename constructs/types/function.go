package types

import (
	"fmt"
	"strings"
)

var _ Type = (*FunctionType)(nil)
var _ NamedType = (*FunctionType)(nil)

// FunctionType for storing the types of function.
type FunctionType struct {

	// Comment is the developer comment for this function.
	Comment string

	// Name is the name of the function.
	Name string

	// ParamNames is the names of all the parameters.
	ParamNames []string

	// ParamTypes is the types of all the parameters.
	ParamTypes []Type

	// Ellipsis indicates if the last parameter is variable length.
	Ellipsis bool

	// ReturnType is the type returned from the function.
	// Void for no return type.
	ReturnType Type

	// ReceiverName is the name for the optional receiver class or empty.
	ReceiverName string

	// ReceiverClass is the class receiver or nil.
	ReceiverClass *ClassType

	// Body is the statement block for the function
	// or nil if the function is a forward declaration.
	// This will be of type *BlockStatement.
	Body interface{}
}

// Function creates a new function type description with the given information.
func Function() *FunctionType {
	return &FunctionType{
		Comment:       "",
		Name:          "",
		ParamNames:    []string{},
		ParamTypes:    []Type{},
		Ellipsis:      false,
		ReturnType:    Void(),
		ReceiverName:  "",
		ReceiverClass: nil,
		Body:          nil,
	}
}

// GetName gets the name of the type.
// May be empty if this type is unnamed.
func (t *FunctionType) GetName() string {
	if t == nil {
		return ""
	}
	return t.Name
}

// SetName sets the name of the function.
func (t *FunctionType) SetName(name string) *FunctionType {
	t.Name = name
	return t
}

// AddParam adds a parameter to the function.
func (t *FunctionType) AddParam(name string, paramType Type) *FunctionType {
	t.ParamNames = append(t.ParamNames, name)
	t.ParamTypes = append(t.ParamTypes, paramType)
	return t
}

// SetEllipse sets the parameter ellipse for the function.
func (t *FunctionType) SetEllipse(ellipsis bool) *FunctionType {
	t.Ellipsis = ellipsis
	return t
}

// SetReturn sets the return type to the function.
func (t *FunctionType) SetReturn(returnType Type) *FunctionType {
	t.ReturnType = returnType
	return t
}

// String gets the name or structure for this type.
func (t *FunctionType) String() string {
	if t == nil {
		return nilStr
	}
	if len(t.Name) > 0 {
		return t.Name
	}
	return t.FullString()
}

// FullString gets a string for the signature and body of this function.
func (t *FunctionType) FullString() string {
	if t == nil {
		return nilStr
	}
	params := "()"
	if paramsCount := len(t.ParamNames); paramsCount > 0 {
		paramStrs := make([]string, paramsCount)
		for i, name := range t.ParamNames {
			paramType := ToString(t.ParamTypes[i])
			if t.Ellipsis && (i == paramsCount-1) {
				paramStrs[i] = paramType + "... " + name
			} else {
				paramStrs[i] = paramType + " " + name
			}
		}
		params = "(" + strings.Join(paramStrs, ", ") + ")"
	}
	name := "func"
	if len(t.Name) > 0 {
		name = t.Name
	}
	bodyStr := ""
	if t.Body != nil {
		bodyStr = " " + fmt.Sprint(t.Body)
	}
	return ToString(t.ReturnType) + " " + name + params + bodyStr
}
