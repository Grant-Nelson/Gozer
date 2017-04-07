package constructs

import (
	"sort"
	"strings"
)

var _ Type = (*InterfaceType)(nil)

// InterfaceType for storing the types of interfaces.
type InterfaceType struct {

	// Parent is the parent type for this function.
	// It is a class, package, or nil.
	Parent Type

	// Name is the name of the interface.
	Name string

	// Functions is the set of functions for this interface.
	Functions map[string]*FunctionType
}

// Interface creates a new interface type.
func Interface() *InterfaceType {
	return &InterfaceType{
		Name:      "",
		Functions: map[string]*FunctionType{},
	}
}

// Find looks up a subtype to this interface.
func (t *InterfaceType) Find(name string) (Type, bool) {
	t2, exists := t.Functions[name]
	return t2, exists
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
	if len(t.Name) > 0 {
		return t.Name
	}
	if len(t.Functions) <= 0 {
		return "interface{}"
	}
	i := 0
	parts := make([]string, len(t.Functions))
	for name, tfunc := range t.Functions {
		parts[i] = name + " " + ToString(tfunc)
		i++
	}
	sort.Strings(parts)
	return "interface{\n  " + strings.Join(parts, "\n  ") + "\n}"
}
