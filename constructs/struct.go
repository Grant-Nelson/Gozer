package constructs

import (
	"sort"
	"strings"
)

var _ Type = (*StructureType)(nil)

// StructureType for storing the types of structure.
type StructureType struct {

	// Parent is the parent type for this function.
	// It is a class, package, or nil.
	Parent Type

	// Name is the name of the structure.
	Name string

	// Members is the set of members for this structure.
	Members map[string]Type
}

// Structure creates a new struct type.
func Structure() *StructureType {
	return &StructureType{
		Name:    "",
		Members: map[string]Type{},
	}
}

// Find looks up a subtype to this structure.
func (t *StructureType) Find(name string) (Type, bool) {
	t2, exists := t.Members[name]
	return t2, exists
}

// AddMember adds a member to this structure.
func (t *StructureType) AddMember(name string, tmem Type) *StructureType {
	t.Members[name] = tmem
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
	if len(t.Members) <= 0 {
		return "struct{}"
	}
	i := 0
	parts := make([]string, len(t.Members))
	for name, member := range t.Members {
		parts[i] = name + " " + ToString(member)
		i++
	}
	sort.Strings(parts)
	return "struct{\n  " + strings.Join(parts, "\n  ") + "\n}"
}
