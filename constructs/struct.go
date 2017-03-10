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
