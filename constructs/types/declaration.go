package types

import "github.com/grant-nelson/Gozer/common"

var _ Type = (*DeclarationType)(nil)
var _ NamedType = (*DeclarationType)(nil)
var _ SubtypableType = (*DeclarationType)(nil)

// DeclarationType for storing the types of declaration.
type DeclarationType struct {

	// Name is the name of the declaration.
	Name string

	// Data is the type of the data in this declaration.
	Data Type
}

// Declaration creates a new declaration type.
func Declaration() *DeclarationType {
	c := &DeclarationType{
		Name: "",
		Data: nil,
	}
	return c
}

// GetName gets the name of the type.
// May be empty if this type is unnamed.
func (t *DeclarationType) GetName() string {
	if t == nil {
		return ""
	}
	return t.Name
}

// Find looks up a subtype to this declaration.
func (t *DeclarationType) Find(name string) (Type, bool) {
	if t == nil {
		return nil, false
	}
	if structType, ok := t.Data.(*StructureType); ok {
		t2, exists := structType.Find(name)
		return t2, exists
	}
	return nil, false
}

// String gets the name for this type.
func (t *DeclarationType) String() string {
	if t == nil {
		return nilStr
	}
	if len(t.Name) > 0 {
		return t.Name
	}
	return t.FullString()
}

// FullString gets a string for the structure for this type.
func (t *DeclarationType) FullString() string {
	if t == nil {
		return nilStr
	}
	name := "decl"
	if len(t.Name) > 0 {
		name = t.Name
	}
	data := ""
	if str := ToString(t.Data); len(str) > 0 {
		data = common.Indent(str, "  ") + " "
	}
	return data + name
}
