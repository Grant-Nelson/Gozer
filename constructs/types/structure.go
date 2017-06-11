package types

import (
	"github.com/grant-nelson/Gozer/common"
)

var _ Type = (*StructureType)(nil)
var _ NamedType = (*StructureType)(nil)
var _ SubtypableType = (*StructureType)(nil)

// StructureType for storing the types of structure.
type StructureType struct {

	// Name is the name of the structure.
	Name string

	// Members is the data members for this structure.
	Members *DeclarationSet
}

// Structure creates a new struct type.
func Structure() *StructureType {
	return &StructureType{
		Name:    "",
		Members: NewDeclarationSet(),
	}
}

// GetName gets the name of the type.
// May be empty if this type is unnamed.
func (t *StructureType) GetName() string {
	if t == nil {
		return ""
	}
	return t.Name
}

// Find looks up a subtype to this structure.
func (t *StructureType) Find(name string) (Type, bool) {
	if (t == nil) || (t.Members == nil) {
		return nil, false
	}
	return t.Members.Find(name)
}

// AddMember adds a member to this structure.
// If the member already exists it will be overwritten with the new type.
func (t *StructureType) AddMember(name string, data Type) *DeclarationType {
	t2, _ := t.Members.AddNew(name, data)
	return t2
}

// String gets the name for this type.
func (t *StructureType) String() string {
	if t == nil {
		return nilStr
	}
	if len(t.Name) > 0 {
		return t.Name
	}
	return t.FullString()
}

// FullString gets a string for the structure for this type.
func (t *StructureType) FullString() string {
	if t == nil {
		return nilStr
	}
	name := "struct"
	if len(t.Name) > 0 {
		name = t.Name
	}
	if t.Members.Len() <= 0 {
		return name + "{}"
	}
	return name + "{\n" +
		"  " + common.Indent(t.Members.FullString(), "  ") + "\n" +
		"}"
}
