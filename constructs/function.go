package constructs

import "strings"

var _ Type = (*FunctionType)(nil)

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

	// ReturnNames is the names of all the return types.
	ReturnNames []string

	// ReturnTypes is the types of all the return types.
	ReturnTypes []Type

	// ReceiverName is the name for the optional receiver class or empty.
	ReceiverName string

	// ReceiverClass is the class receiver or nil.
	ReceiverClass *ClassType

	// Body is the statement block for the function
	// or nil if the function is a forward declaration.
	Body *BlockStatement
}

// Function creates a new function type description with the given information.
func Function() *FunctionType {
	return &FunctionType{
		Comment:       "",
		Name:          "",
		ParamNames:    []string{},
		ParamTypes:    []Type{},
		Ellipsis:      false,
		ReturnNames:   []string{},
		ReturnTypes:   []Type{},
		ReceiverName:  "",
		ReceiverClass: nil,
		Body:          nil,
	}
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

// AddReturn adds a return type to the function.
func (t *FunctionType) AddReturn(name string, returnType Type) *FunctionType {
	t.ReturnNames = append(t.ReturnNames, name)
	t.ReturnTypes = append(t.ReturnTypes, returnType)
	return t
}

// Equals determins if the given type is the same as this type.
func (t *FunctionType) Equals(other interface{}) bool {
	if t == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	tother, ok := other.(*FunctionType)
	if !ok {
		return false
	}
	if t.ReceiverName != tother.ReceiverName {
		return false
	}
	if len(t.ParamNames) != len(tother.ParamNames) {
		return false
	}
	if len(t.ParamTypes) != len(tother.ParamTypes) {
		return false
	}
	if len(t.ReturnNames) != len(tother.ReturnNames) {
		return false
	}
	if len(t.ReturnTypes) != len(tother.ReturnTypes) {
		return false
	}
	for i, name := range t.ParamNames {
		if name != tother.ParamNames[i] {
			return false
		}
	}
	for i, paramType := range t.ParamTypes {
		if !Equals(paramType, tother.ParamTypes[i]) {
			return false
		}
	}
	for i, name := range t.ReturnNames {
		if name != tother.ReturnNames[i] {
			return false
		}
	}
	for i, returnType := range t.ReturnTypes {
		if !Equals(returnType, tother.ReturnTypes[i]) {
			return false
		}
	}
	if !Equals(t.Body, tother.Body) {
		return false
	}
	// Don't check ReceiverClass
	return true
}

// String gets the name for this type.
func (t *FunctionType) String() string {
	if t == nil {
		return nilStr
	}
	params := "()"
	if paramsCount := len(t.ParamNames); paramsCount > 0 {
		paramStrs := make([]string, paramsCount)
		for i, name := range t.ParamNames {
			paramType := ToString(t.ParamTypes[i])
			if t.Ellipsis && (i == paramsCount-1) {
				paramStrs[i] = name + " ..." + paramType
			} else {
				paramStrs[i] = name + " " + paramType
			}
		}
		params = "(" + strings.Join(paramStrs, ", ") + ")"
	}

	returns := "()"
	if returnCount := len(t.ReturnNames); returnCount > 0 {
		returnStrs := make([]string, returnCount)
		for i, name := range t.ReturnNames {
			returnType := ToString(t.ParamTypes[i])
			if len(name) > 0 {
				returnStrs[i] = name + " " + returnType
			} else {
				returnStrs[i] = returnType
			}
		}
		returns = "(" + strings.Join(returnStrs, ", ") + ")"
	}

	return "func " + t.Name + params + returns + ToString(t.Body)
}
