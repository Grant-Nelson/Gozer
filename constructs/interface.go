package constructs

import (
	"sort"
	"strings"
)

var _ Type = (*InterfaceType)(nil)

// InterfaceType for storing the types of interfaces.
type InterfaceType struct {

	// Functions is the set of functions for this interface.
	Functions map[string]*FunctionType
}

// Interface creates a new interface type.
func Interface() *InterfaceType {
	return &InterfaceType{
		Functions: map[string]*FunctionType{},
	}
}

// AddFunction adds a function to this interface.
func (t *InterfaceType) AddFunction(name string) *FunctionType {
	tfunc := Function()
	t.Functions[name] = tfunc
	return tfunc
}

// String gets the name for this type.
func (t *InterfaceType) String() string {
	if t == nil {
		return nilStr
	}
	i := 0
	parts := make([]string, len(t.Functions))
	for name, tfunc := range t.Functions {
		parts[i] = name + " " + ToString(tfunc)
		i++
	}
	sort.Strings(parts)
	return "{\n  " + strings.Join(parts, "\n  ") + "\n}"
}
