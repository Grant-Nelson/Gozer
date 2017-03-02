package constructs

import "strings"

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

// Equals determins if the given type is the same as this type.
func (t *PackageType) Equals(other Type) bool {
	if t == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	tother, ok := other.(*PackageType)
	if !ok {
		return false
	}
	if len(t.Imports) != len(tother.Imports) {
		return false
	}
	if len(t.Declarations) != len(tother.Declarations) {
		return false
	}
	if len(t.Functions) != len(tother.Functions) {
		return false
	}
	if len(t.Interfaces) != len(tother.Interfaces) {
		return false
	}
	if len(t.Classes) != len(tother.Classes) {
		return false
	}
	for name, imp := range t.Imports {
		timp, exists := tother.Imports[name]
		if !exists {
			return false
		}
		if !imp.Equals(timp) {
			return false
		}
	}
	for name, decl := range t.Declarations {
		tdecl, exists := tother.Declarations[name]
		if !exists {
			return false
		}
		if !decl.Equals(tdecl) {
			return false
		}
	}
	for name, tfunc := range t.Functions {
		ttfunc, exists := tother.Functions[name]
		if !exists {
			return false
		}
		if !tfunc.Equals(ttfunc) {
			return false
		}
	}
	for name, inter := range t.Interfaces {
		tinter, exists := tother.Interfaces[name]
		if !exists {
			return false
		}
		if !inter.Equals(tinter) {
			return false
		}
	}
	for name, class := range t.Classes {
		tclass, exists := tother.Classes[name]
		if !exists {
			return false
		}
		if !class.Equals(tclass) {
			return false
		}
	}
	return true
}

// String gets the name for this type.
func (t *PackageType) String() string {
	if t == nil {
		return nilStr
	}
	parts := make([]string, len(t.Imports)+len(t.Declarations)+len(t.Functions)+len(t.Interfaces)+len(t.Classes))
	i := 0
	for name := range t.Imports {
		parts[i] = "import " + name
		i++
	}
	for name, decl := range t.Declarations {
		parts[i] = name + " " + decl.String()
		i++
	}
	for name, tfunc := range t.Functions {
		parts[i] = name + " " + tfunc.String()
		i++
	}
	for name, inter := range t.Interfaces {
		parts[i] = name + " " + inter.String()
		i++
	}
	for name, class := range t.Classes {
		parts[i] = name + " " + class.String()
		i++
	}
	return "{\n" +
		indent(strings.Join(parts, "\n"), "  ") + "\n" +
		"}"
}
