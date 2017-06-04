package types

import (
	"sort"
	"strings"

	"github.com/grant-nelson/Gozer/common"
)

var _ Type = (*PackageType)(nil)
var _ NamedType = (*PackageType)(nil)
var _ SubtypableType = (*PackageType)(nil)

// PackageType for storing the types of package.
type PackageType struct {

	// Name is the name of the package.
	Name string

	// Imports are the list of packages used by this package.
	// The key is the path to the package.
	Imports map[string]*PackageType

	// Declarations gets the set of constant and variable declarations for this package.
	Declarations map[string]Type

	// Functions is the set of functions for this package.
	Functions map[string]*FunctionType

	// Interfaces is the set of interfaces for this package.
	Interfaces *InterfaceSet

	// Classes is the set of classes for this package.
	Classes *ClassSet

	// Structures is the set of structure for this package.
	Structures map[string]*StructureType
}

// Package creates a new package description.
func Package() *PackageType {
	return &PackageType{
		Name:         "",
		Imports:      map[string]*PackageType{},
		Declarations: map[string]Type{},
		Functions:    map[string]*FunctionType{},
		Interfaces:   &InterfaceSet{},
		Classes:      &ClassSet{},
		Structures:   map[string]*StructureType{},
	}
}

// GetName gets the name of the type.
// May be empty if this type is unnamed.
func (t *PackageType) GetName() string {
	if t == nil {
		return ""
	}
	return t.Name
}

// Find looks up a subtype to this package.
func (t *PackageType) Find(name string) (Type, bool) {
	if t == nil {
		return nil, false
	}
	if t2, exists := t.Imports[name]; exists {
		return t2, true
	}
	if t2, exists := t.Declarations[name]; exists {
		return t2, true
	}
	if t2, exists := t.Functions[name]; exists {
		return t2, true
	}
	if t2 := t.Interfaces.Find(name); t2 != nil {
		return t2, true
	}
	if t2 := t.Classes.Find(name); t2 != nil {
		return t2, true
	}
	if t2, exists := t.Structures[name]; exists {
		return t2, true
	}
	return nil, false
}

// AddDeclaration adds a declaration to this package.
func (t *PackageType) AddDeclaration(name string, decl Type) *PackageType {
	if t == nil {
		return nil
	}
	t.Declarations[name] = decl
	return t
}

// AddFunction adds a function to this package.
func (t *PackageType) AddFunction(name string) *FunctionType {
	if t == nil {
		return nil
	}
	tfunc := Function()
	tfunc.Name = name
	t.Functions[name] = tfunc
	return tfunc
}

// AddInterface adds a interface to this package.
func (t *PackageType) AddInterface(name string) *InterfaceType {
	if t == nil {
		return nil
	}
	return t.Interfaces.AddNew(name)
}

// AddClass adds a class to this package.
func (t *PackageType) AddClass(name string) *ClassType {
	if t == nil {
		return nil
	}
	return t.Classes.AddNew(name)
}

// AddStructure adds a structure to this package.
func (t *PackageType) AddStructure(name string) *StructureType {
	if t == nil {
		return nil
	}
	st := Structure()
	st.Name = name
	t.Structures[name] = st
	return st
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
		result += "  " + common.Indent(strings.Join(parts1, "\n"), "  ") + "\n"
	}

	if len(t.Declarations) > 0 {
		i := 0
		parts2 := make([]string, len(t.Declarations))
		for name, decl := range t.Declarations {
			parts2[i] = name + " " + ToString(decl)
			i++
		}
		sort.Strings(parts2)
		result += "  " + common.Indent(strings.Join(parts2, "\n"), "  ") + "\n"
	}

	if len(t.Functions) > 0 {
		i := 0
		parts3 := make([]string, len(t.Functions))
		for _, tfunc := range t.Functions {
			parts3[i] = tfunc.FullString()
			i++
		}
		sort.Strings(parts3)
		result += "  " + common.Indent(strings.Join(parts3, "\n"), "  ") + "\n"
	}

	if t.Interfaces.Len() > 0 {
		result += "  " + common.Indent(t.Interfaces.FullString(), "  ") + "\n"
	}

	if t.Classes.Len() > 0 {
		result += "  " + common.Indent(t.Classes.FullString(), "  ") + "\n"
	}

	if len(t.Structures) > 0 {
		i := 0
		parts6 := make([]string, len(t.Structures))
		for _, strt := range t.Structures {
			parts6[i] = strt.FullString()
			i++
		}
		sort.Strings(parts6)
		result += "  " + common.Indent(strings.Join(parts6, "\n"), "  ") + "\n"
	}

	if len(result) <= 0 {
		return "{}"
	}
	return "{\n" + result + "}"
}
