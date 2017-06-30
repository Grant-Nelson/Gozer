package types

import (
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
	Imports *PackageSet

	// Declarations gets the set of constant and variable declarations for this package.
	Declarations *DeclarationSet

	// Functions is the set of functions for this package.
	Functions *FunctionSet

	// Interfaces is the set of interfaces for this package.
	Interfaces *InterfaceSet

	// Classes is the set of classes for this package.
	Classes *ClassSet

	// Structures is the set of structure for this package.
	Structures *StructureSet

	// ReturnSets is the set of return sets for this package.
	ReturnSets *ReturnSetSet
}

// Package creates a new package description.
func Package() *PackageType {
	return &PackageType{
		Name:         "",
		Imports:      NewPackageSet(),
		Declarations: NewDeclarationSet(),
		Functions:    NewFunctionSet(),
		Interfaces:   NewInterfaceSet(),
		Classes:      NewClassSet(),
		Structures:   NewStructureSet(),
		ReturnSets:   NewReturnSetSet(),
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
	if t2, found := t.Imports.Find(name); found {
		return t2, true
	}
	if t2, exists := t.Declarations.Find(name); exists {
		return t2, true
	}
	if t2, found := t.Functions.Find(name); found {
		return t2, true
	}
	if t2, found := t.Interfaces.Find(name); found {
		return t2, true
	}
	if t2, found := t.Classes.Find(name); found {
		return t2, true
	}
	if t2, found := t.Structures.Find(name); found {
		return t2, true
	}
	if t2, found := t.ReturnSets.Find(name); found {
		return t2, true
	}
	return nil, false
}

// AddImport adds a package to this package.
func (t *PackageType) AddImport(name string) *PackageType {
	if t == nil {
		return nil
	}
	t2, _ := t.Imports.AddNew(name)
	return t2
}

// AddDeclaration adds a declaration to this package.
func (t *PackageType) AddDeclaration(name string, data Type) *DeclarationType {
	if t == nil {
		return nil
	}
	t2, _ := t.Declarations.AddNew(name, data)
	return t2
}

// AddFunction adds a function to this package.
func (t *PackageType) AddFunction(name string) *FunctionType {
	if t == nil {
		return nil
	}
	t2, _ := t.Functions.AddNew(name)
	return t2
}

// AddInterface adds a interface to this package.
func (t *PackageType) AddInterface(name string) *InterfaceType {
	if t == nil {
		return nil
	}
	t2, _ := t.Interfaces.AddNew(name)
	return t2
}

// AddClass adds a class to this package.
func (t *PackageType) AddClass(name string) *ClassType {
	if t == nil {
		return nil
	}
	t2, _ := t.Classes.AddNew(name)
	return t2
}

// AddStructure adds a structure to this package.
func (t *PackageType) AddStructure(name string) *StructureType {
	if t == nil {
		return nil
	}
	t2, _ := t.Structures.AddNew(name)
	return t2
}

// AddReturnSet adds a return set to this package.
func (t *PackageType) AddReturnSet(name string) *ReturnSet {
	if t == nil {
		return nil
	}
	t2, _ := t.ReturnSets.AddNew(name)
	return t2
}

// String gets the name for this type.
func (t *PackageType) String() string {
	return t.StringWithShort("")
}

// StringWithShort gets the name for this type with an optional short name.
func (t *PackageType) StringWithShort(short string) string {
	if t == nil {
		return nilStr
	}
	if len(short) > 0 {
		if len(t.Name) > 0 {
			return "import " + short + " = " + t.Name
		}
		return "import " + short
	}
	if len(t.Name) > 0 {
		return "import " + t.Name
	}
	return t.FullStringWithShort(short)
}

// FullString gets the full name for this type.
func (t *PackageType) FullString() string {
	return t.FullStringWithShort("")
}

// FullStringWithShort gets the full name for this type with an optional short name.
func (t *PackageType) FullStringWithShort(short string) string {
	if t == nil {
		return nilStr
	}
	name := "import"
	if len(short) > 0 {
		name = "import " + short
	} else if len(t.Name) > 0 {
		name = "import " + t.Name
	}

	result := ""
	if t.Imports.Len() > 0 {
		// Use String not FullString for imports
		result += "  " + common.Indent(t.Imports.String(), "  ") + "\n"
	}

	if t.Declarations.Len() > 0 {
		result += "  " + common.Indent(t.Declarations.FullString(), "  ") + "\n"
	}

	if t.Functions.Len() > 0 {
		result += "  " + common.Indent(t.Functions.FullBodyString(), "  ") + "\n"
	}

	if t.Interfaces.Len() > 0 {
		result += "  " + common.Indent(t.Interfaces.FullString(), "  ") + "\n"
	}

	if t.Classes.Len() > 0 {
		result += "  " + common.Indent(t.Classes.FullString(), "  ") + "\n"
	}

	if t.Structures.Len() > 0 {
		result += "  " + common.Indent(t.Structures.FullString(), "  ") + "\n"
	}

	if t.ReturnSets.Len() > 0 {
		result += "  " + common.Indent(t.ReturnSets.FullString(), "  ") + "\n"
	}

	if len(result) <= 0 {
		return name + "{}"
	}
	return name + "{\n" + result + "}"
}
