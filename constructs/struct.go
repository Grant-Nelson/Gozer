package constructs

import (
	"sort"
	"strings"
)

var _ Type = (*StructureType)(nil)

// StructureType for storing the types of structure.
type StructureType struct {

	// Members is the set of members for this structure.
	Members map[string]Type
}

// Structure creates a new struct type.
func Structure() *StructureType {
	return &StructureType{
		Members: map[string]Type{},
	}
}

// Equals determins if the given type is the same as this type.
func (t *StructureType) Equals(other interface{}) bool {
	if t == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	tother, ok := other.(*StructureType)
	if !ok {
		return false
	}
	if len(t.Members) != len(tother.Members) {
		return false
	}
	for name, member := range t.Members {
		tmem, exists := tother.Members[name]
		if !exists {
			return false
		}
		if !Equals(member, tmem) {
			return false
		}
	}
	return true
}

// String gets the name for this type.
func (t *StructureType) String() string {
	if t == nil {
		return nilStr
	}
	i := 0
	parts := make([]string, len(t.Members))
	for name, member := range t.Members {
		parts[i] = name + " " + ToString(member)
		i++
	}
	sort.Strings(parts)
	return "{\n  " + strings.Join(parts, "\n  ") + "\n}"
}
