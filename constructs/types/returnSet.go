package types

import (
	"github.com/grant-nelson/Gozer/common"
)

var _ Type = (*ReturnSet)(nil)
var _ NamedType = (*ReturnSet)(nil)
var _ SubtypableType = (*ReturnSet)(nil)

// ReturnSet for storing the types of return sets.
// Designed to store multiple types for a multiple return set.
type ReturnSet struct {

	// Name is the name of the return set.
	Name string

	// Members is the data members for this return set.
	Members *DeclarationSet
}

// NewReturnSet creates a new return set type.
func NewReturnSet() *ReturnSet {
	mem := NewDeclarationSet()
	mem.AutoSort = false
	return &ReturnSet{
		Name:    "",
		Members: mem,
	}
}

// GetName gets the name of the type.
// May be empty if this type is unnamed.
func (t *ReturnSet) GetName() string {
	if t == nil {
		return ""
	}
	return t.Name
}

// Find looks up a subtype to this return set.
func (t *ReturnSet) Find(name string) (Type, bool) {
	if (t == nil) || (t.Members == nil) {
		return nil, false
	}
	return t.Members.Find(name)
}

// AddMember adds a member to this return set.
// If the member already exists it will be overwritten with the new type.
func (t *ReturnSet) AddMember(name string, data Type) *DeclarationType {
	t2, _ := t.Members.AddNew(name, data)
	return t2
}

// String gets the name for this type.
func (t *ReturnSet) String() string {
	if t == nil {
		return nilStr
	}
	if len(t.Name) > 0 {
		return t.Name
	}
	return t.FullString()
}

// FullString gets a string for the return set for this type.
func (t *ReturnSet) FullString() string {
	if t == nil {
		return nilStr
	}
	name := "returns"
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
