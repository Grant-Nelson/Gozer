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

// AddReturn adds a return type to the function.
func (t *FunctionType) AddReturn(name string, returnType Type) *FunctionType {
	t.ReturnNames = append(t.ReturnNames, name)
	t.ReturnTypes = append(t.ReturnTypes, returnType)
	return t
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

	returns := ""
	if returnCount := len(t.ReturnNames); returnCount > 0 {
		returnStrs := make([]string, returnCount)
		for i, name := range t.ReturnNames {
			returnType := ToString(t.ReturnTypes[i])
			if len(name) > 0 {
				returnStrs[i] = name + " " + returnType
			} else {
				returnStrs[i] = returnType
			}
		}
		returns = "(" + strings.Join(returnStrs, ", ") + ")"
	}

	bodyStr := ""
	if t.Body != nil {
		bodyStr = " " + ToString(t.Body)
	}
	name := ""
	if len(t.Name) > 0 {
		name = " " + t.Name
	}
	return "func" + name + params + returns + bodyStr
}
