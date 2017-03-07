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
func (t *PackageType) Equals(other interface{}) bool {
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
		if !Equals(imp, timp) {
			return false
		}
	}
	for name, decl := range t.Declarations {
		tdecl, exists := tother.Declarations[name]
		if !exists {
			return false
		}
		if !Equals(decl, tdecl) {
			return false
		}
	}
	for name, tfunc := range t.Functions {
		ttfunc, exists := tother.Functions[name]
		if !exists {
			return false
		}
		if !Equals(tfunc, ttfunc) {
			return false
		}
	}
	for name, inter := range t.Interfaces {
		tinter, exists := tother.Interfaces[name]
		if !exists {
			return false
		}
		if !Equals(inter, tinter) {
			return false
		}
	}
	for name, class := range t.Classes {
		tclass, exists := tother.Classes[name]
		if !exists {
			return false
		}
		if !Equals(class, tclass) {
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

	i := 0
	parts1 := make([]string, len(t.Imports))
	for name := range t.Imports {
		parts1[i] = "import " + name
		i++
	}
	sort.Strings(parts1)

	i = 0
	parts2 := make([]string, len(t.Declarations))
	for name, decl := range t.Declarations {
		parts2[i] = name + " " + ToString(decl)
		i++
	}
	sort.Strings(parts2)

	i = 0
	parts3 := make([]string, len(t.Functions))
	for name, tfunc := range t.Functions {
		parts3[i] = name + " " + ToString(tfunc)
		i++
	}
	sort.Strings(parts3)

	i = 0
	parts4 := make([]string, len(t.Interfaces))
	for name, inter := range t.Interfaces {
		parts4[i] = name + " " + ToString(inter)
		i++
	}
	sort.Strings(parts4)

	i = 0
	parts5 := make([]string, len(t.Classes))
	for name, class := range t.Classes {
		parts5[i] = name + " " + ToString(class)
		i++
	}
	sort.Strings(parts5)

	return "{\n" +
		indent(strings.Join(parts1, "\n"), "  ") + "\n" +
		indent(strings.Join(parts2, "\n"), "  ") + "\n" +
		indent(strings.Join(parts3, "\n"), "  ") + "\n" +
		indent(strings.Join(parts4, "\n"), "  ") + "\n" +
		indent(strings.Join(parts5, "\n"), "  ") + "\n" +
		"}"
}
