package constructs

import "strings"

var _ Type = (*FunctionType)(nil)

// FunctionType for storing the types of function.
type FunctionType struct {

	// Comment is the developer comment for this function.
	Comment string

	// ParamNames is the names of all the parameters.
	ParamNames []string

	// ParamTypes is the types of all the parameters.
	ParamTypes []Type

	// Ellipsis indicates if the last parameter is variable length.
	Ellipsis bool

	// ResultNames is the names of all the return types.
	ResultNames []string

	// ResultTypes is the types of all the return types.
	ResultTypes []Type

	// ReceiverName is the name for the optional receiver class or empty.
	ReceiverName string

	// ReceiverClass is the class receiver or nil.
	ReceiverClass *ClassType
}

// Function creates a new function type description with the given information.
func Function() *FunctionType {
	return &FunctionType{
		Comment:       "",
		ParamNames:    []string{},
		ParamTypes:    []Type{},
		Ellipsis:      false,
		ResultNames:   []string{},
		ResultTypes:   []Type{},
		ReceiverName:  "",
		ReceiverClass: nil,
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

// AddResult adds a return type to the function.
func (t *FunctionType) AddResult(name string, resultType Type) *FunctionType {
	t.ResultNames = append(t.ResultNames, name)
	t.ResultTypes = append(t.ResultTypes, resultType)
	return t
}

// Equals determins if the given type is the same as this type.
func (t *FunctionType) Equals(other Type) bool {
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
	if len(t.ResultNames) != len(tother.ResultNames) {
		return false
	}
	if len(t.ResultTypes) != len(tother.ResultTypes) {
		return false
	}
	for i, name := range t.ParamNames {
		if name != tother.ParamNames[i] {
			return false
		}
	}
	for i, paramType := range t.ParamTypes {
		if !paramType.Equals(tother.ParamTypes[i]) {
			return false
		}
	}
	for i, name := range t.ResultNames {
		if name != tother.ResultNames[i] {
			return false
		}
	}
	for i, resultType := range t.ResultTypes {
		if !resultType.Equals(tother.ResultTypes[i]) {
			return false
		}
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
			paramType := t.ParamTypes[i].String()
			if t.Ellipsis && (i == paramsCount-1) {
				paramStrs[i] = name + " ..." + paramType
			} else {
				paramStrs[i] = name + " " + paramType
			}
		}
		params = "(" + strings.Join(paramStrs, ", ") + ")"
	}

	result := ""
	if resultCount := len(t.ResultNames); resultCount > 0 {
		resultStrs := make([]string, resultCount)
		for i, name := range t.ResultNames {
			resultType := t.ParamTypes[i].String()
			if len(name) > 0 {
				resultStrs[i] = name + " " + resultType
			} else {
				resultStrs[i] = resultType
			}
		}
		result = "(" + strings.Join(resultStrs, ", ") + ")"
	}

	return "func" + params + result
}
