package types

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
	Functions []*FunctionType
}

// Interface creates a new interface type.
func Interface() *InterfaceType {
	return &InterfaceType{
		Parent:    nil,
		Name:      "",
		Functions: []*FunctionType{},
	}
}

// Find looks up a subtype to this interface.
func (t *InterfaceType) Find(name string) (Type, bool) {
	for _, tfunc := range t.Functions {
		if tfunc.Name == name {
			return tfunc, true
		}
	}
	return nil, false
}

// AddFunction adds a function to this interface.
// If a function by that name already exists, that function is returned.
func (t *InterfaceType) AddFunction(name string) *FunctionType {
	if tfunc, found := t.Find(name); found {
		return tfunc.(*FunctionType)
	}
	tfunc := Function()
	tfunc.Name = name
	t.Functions = append(t.Functions, tfunc)
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
	return t.FullString()
}

// FullString gets a string for the structure for this type.
func (t *InterfaceType) FullString() string {
	if t == nil {
		return nilStr
	}
	name := "interface"
	if len(t.Name) > 0 {
		name = t.Name
	}
	if len(t.Functions) <= 0 {
		return name + "{}"
	}
	i := 0
	parts := make([]string, len(t.Functions))
	for _, tfunc := range t.Functions {
		parts[i] = tfunc.FullString()
		i++
	}
	sort.Strings(parts)
	return name + "{\n  " + strings.Join(parts, "\n  ") + "\n}"
}
