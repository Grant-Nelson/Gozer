package constructs

import (
	"sort"
	"strings"
)

var _ Type = (*PackageType)(nil)

// PackageType for storing the types of package.
type PackageType struct {

	// Imports are the list of packages used by this package.
	// The key is the path to the package.
	Imports map[string]*PackageType

	// Declarations gets the set of constant and variable declarations for this package.
	Declarations map[string]Type

	// Functions is the set of functions for this package.
	Functions map[string]*FunctionType

	// Interfaces is the set of interfaces for this package.
	Interfaces map[string]*InterfaceType

	// Classes is the set of classes for this package.
	Classes map[string]*ClassType
}

// Package creates a new package description.
func Package() *PackageType {
	return &PackageType{
		Imports:      map[string]*PackageType{},
		Declarations: map[string]Type{},
		Functions:    map[string]*FunctionType{},
		Interfaces:   map[string]*InterfaceType{},
		Classes:      map[string]*ClassType{},
	}
}

// Find looks up a subtype to this package.
func (t *PackageType) Find(name string) (Type, bool) {
	if t2, exists := t.Imports[name]; exists {
		return t2, true
	}
	if t2, exists := t.Declarations[name]; exists {
		return t2, true
	}
	if t2, exists := t.Functions[name]; exists {
		return t2, true
	}
	if t2, exists := t.Interfaces[name]; exists {
		return t2, true
	}
	if t2, exists := t.Classes[name]; exists {
		return t2, true
	}
	if t2, exists := t.Imports[name]; exists {
		return t2, true
	}
	return nil, false
}

// AddFunction adds a function to this package.
func (t *PackageType) AddFunction(name string) *FunctionType {
	tfunc := Function()
	t.Functions[name] = tfunc
	return tfunc
}

// AddInterface adds a interface to this package.
func (t *PackageType) AddInterface(name string) *InterfaceType {
	inter := Interface()
	t.Interfaces[name] = inter
	return inter
}

// AddClass adds a class to this package.
func (t *PackageType) AddClass(name string) *ClassType {
	class := Class()
	t.Classes[name] = class
	return class
}

// String gets the name for this type.
func (t *PackageType) String() string {
	if t == nil {
		return nilStr
	}
	result := ""

	if len(t.Imports) > 0 {
		i := 0
		parts1 := make([]string, len(t.Imports))
		for name := range t.Imports {
			parts1[i] = "import " + name
			i++
		}
		sort.Strings(parts1)
		result += indent(strings.Join(parts1, "\n"), "  ") + "\n"
	}

	if len(t.Declarations) > 0 {
		i := 0
		parts2 := make([]string, len(t.Declarations))
		for name, decl := range t.Declarations {
			parts2[i] = name + " " + ToString(decl)
			i++
		}
		sort.Strings(parts2)
		result += indent(strings.Join(parts2, "\n"), "  ") + "\n"
	}

	if len(t.Functions) > 0 {
		i := 0
		parts3 := make([]string, len(t.Functions))
		for name, tfunc := range t.Functions {
			parts3[i] = name + " " + ToString(tfunc)
			i++
		}
		sort.Strings(parts3)
		result += indent(strings.Join(parts3, "\n"), "  ") + "\n"
	}

	if len(t.Interfaces) > 0 {
		i := 0
		parts4 := make([]string, len(t.Interfaces))
		for name, inter := range t.Interfaces {
			parts4[i] = name + " " + ToString(inter)
			i++
		}
		sort.Strings(parts4)
		result += indent(strings.Join(parts4, "\n"), "  ") + "\n"
	}

	if len(t.Classes) > 0 {
		i := 0
		parts5 := make([]string, len(t.Classes))
		for name, class := range t.Classes {
			parts5[i] = name + " " + ToString(class)
			i++
		}
		sort.Strings(parts5)
		result += indent(strings.Join(parts5, "\n"), "  ") + "\n"
	}

	if len(result) > 0 {
		return "{}"
	}
	return "{\n" + result + "}"
}
