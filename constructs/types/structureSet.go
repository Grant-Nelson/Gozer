package types

import (
	"sort"
	"strings"
)

// StructureSet for storing a set of structures.
type StructureSet struct {

	// Structs is the set of structures.
	Structs []*StructureType
}

// NewStructureSet creates a new set of structures.
func NewStructureSet() *StructureSet {
	return &StructureSet{
		Structs: []*StructureType{},
	}
}

// AddNew adds a structure to this set.
// Returns the new interface and true if new, false if already exists.
func (set *StructureSet) AddNew(name string) (*StructureType, bool) {
	if set == nil {
		return nil, false
	}
	if st, found := set.Find(name); found {
		return st, false
	}
	st := Structure()
	st.Name = name
	set.Add(st)
	return st, true
}

// Add will append all non-nil structure to this set.
func (set *StructureSet) Add(sts ...*StructureType) {
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
func (set *StructureSet) Find(name string) (*StructureType, bool) {
	if set != nil {
		for _, st := range set.Structs {
			if (st != nil) && (st.Name == name) {
				return st, true
			}
		}
	}
	return nil, false
}

// Sort will sort the structures by name.
func (set *StructureSet) Sort() {
	sort.Sort(set)
}

// Len get the number of structures in this set.
func (set *StructureSet) Len() int {
	if set != nil {
		return len(set.Structs)
	}
	return 0
}

// Swap swaps the structures at the two given indices.
func (set *StructureSet) Swap(aIndex, bIndex int) {
	if (set != nil) && (aIndex >= 0) && (bIndex >= 0) && (aIndex != bIndex) {
		if length := set.Len(); (aIndex < length) && (bIndex < length) {
			set.Structs[aIndex], set.Structs[bIndex] = set.Structs[bIndex], set.Structs[aIndex]
		}
	}
}

// Less determines if the structure's name at the first index
// is less than the structure's name at the second index.
func (set *StructureSet) Less(aIndex, bIndex int) bool {
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
func (set *StructureSet) String() string {
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
func (set *StructureSet) FullString() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, st := range set.Structs {
		parts[i] = st.FullString()
	}
	return strings.Join(parts, "\n")
}
