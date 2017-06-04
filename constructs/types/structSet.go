package types

import (
	"sort"
	"strings"
)

// StructSet for storing a set of structures.
type StructSet struct {

	// Structs is the set of structures.
	Structs []*StructureType
}

// NewStructureSet creates a new set of structures.
func NewStructureSet() *StructSet {
	return &StructSet{
		Structs: []*StructureType{},
	}
}

// AddNew adds a structure to this set.
func (set *StructSet) AddNew(name string) *StructureType {
	if set == nil {
		return nil
	}
	st := Structure()
	st.Name = name
	set.Add(switch)
	return st
}

// Add will append all non-nil structure to this set.
func (set *StructSet) Add(sts ...*StructureType) {
	if set != nil {
		for _, st := range sts {
			if st != nil {
				set.Structs = append(set.Structs, st)
			}
		}
		set.Sort()
	}
}

// Find searches the set of structures to find a structure with the given name.
// If no structure by the given name is found, nil is returned.
func (set *StructSet) Find(name string) *StructureType {
	if set != nil {
		for _, st := range set.Structs {
			if (st != nil) && (st.Name == name) {
				return st
			}
		}
	}
	return nil
}

// Sort will sort the structures by name.
func (set *StructSet) Sort() {
	sort.Sort(set)
}

// Len get the number of structures in this set.
func (set *StructSet) Len() int {
	if set != nil {
		return len(set.Structs)
	}
	return 0
}

// Swap swaps the structures at the two given indices.
func (set *StructSet) Swap(aIndex, bIndex int) {
	if (set != nil) && (aIndex >= 0) && (bIndex >= 0) && (aIndex != bIndex) {
		if length := set.Len(); (aIndex < length) && (bIndex < length) {
			set.Structs[aIndex], set.Structs[bIndex] = set.Structs[bIndex], set.Structs[aIndex]
		}
	}
}

// Less determines if the structure's name at the first index
// is less than the structure's name at the second index.
func (set *StructSet) Less(aIndex, bIndex int) bool {
	aName, bName := "", ""
	if set != nil {
		if a := set.Structs[aIndex]; a != nil {
			aName = a.Name
		}
		if b := set.Structs[bIndex]; b != nil {
			bName = b.Name
		}
	}
	return aName < bName
}

// String gets the string of all the structures in this set.
func (set *StructSet) String() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, st := range set.Structs {
		parts[i] = st.String()
	}
	return strings.Join(parts, "\n")
}

// FullString gets the full string of all the structures in this set.
func (set *StructSet) FullString() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, st := range set.Structs {
		parts[i] = st.FullString()
	}
	return strings.Join(parts, "\n")
}
