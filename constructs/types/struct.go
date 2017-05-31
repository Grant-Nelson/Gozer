package types

import "strings"

var _ Type = (*StructureType)(nil)

// StructureType for storing the types of structure.
type StructureType struct {

	// Parent is the parent type for this structure.
	// It is a class, package, or nil.
	Parent Type

	// Name is the name of the structure.
	Name string

	// MemeberNames is the names of members for this structure.
	MemeberNames []string

	// MemeberTypes is the types of members for this structure.
	MemeberTypes []Type
}

// Structure creates a new struct type.
func Structure() *StructureType {
	return &StructureType{
		Parent:       nil,
		Name:         "",
		MemeberNames: []string{},
		MemeberTypes: []Type{},
	}
}

// Find looks up a subtype to this structure.
func (t *StructureType) Find(name string) (Type, bool) {
	for i, other := range t.MemeberNames {
		if name == other {
			return t.MemeberTypes[i], true
		}
	}
	return nil, false
}

// AddMember adds a member to this structure.
// If the member already exists it will be overwritten with the new type.
func (t *StructureType) AddMember(name string, tmem Type) *StructureType {
	for i, other := range t.MemeberNames {
		if name == other {
			t.MemeberTypes[i] = tmem
			return t
		}
	}
	t.MemeberNames = append(t.MemeberNames, name)
	t.MemeberTypes = append(t.MemeberTypes, tmem)
	return t
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
	if len(t.MemeberNames) <= 0 {
		return name + "{}"
	}
	parts := make([]string, len(t.MemeberNames))
	for i, name := range t.MemeberNames {
		parts[i] = ToString(t.MemeberTypes[i]) + " " + name
		i++
	}
	return name + "{\n  " + strings.Join(parts, "\n  ") + "\n}"
}
