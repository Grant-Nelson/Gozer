package types

import (
	"sort"
	"strings"
)

// ReturnSetSet for storing a set of return sets.
type ReturnSetSet struct {

	// Sets is the set of return sets.
	Sets []*ReturnSet
}

// NewReturnSetSet creates a new set of return sets.
func NewReturnSetSet() *ReturnSetSet {
	return &ReturnSetSet{
		Sets: []*ReturnSet{},
	}
}

// AddNew adds a structure to this set.
// Returns the new interface and true if new, false if already exists.
func (set *ReturnSetSet) AddNew(name string) (*ReturnSet, bool) {
	if set == nil {
		return nil, false
	}
	if st, found := set.Find(name); found {
		return st, false
	}
	st := NewReturnSet()
	st.Name = name
	set.Add(st)
	return st, true
}

// Add will append all non-nil structure to this set.
func (set *ReturnSetSet) Add(sts ...*ReturnSet) *ReturnSetSet {
	if set != nil {
		for _, st := range sts {
			if st != nil {
				set.Sets = append(set.Sets, st)
			}
		}
		set.Sort()
	}
	return set
}

// Find searches the set of structures to find a structure with the given name.
// If no structure by the given name is found, nil is returned.
func (set *ReturnSetSet) Find(name string) (*ReturnSet, bool) {
	if set != nil {
		for _, st := range set.Sets {
			if (st != nil) && (st.Name == name) {
				return st, true
			}
		}
	}
	return nil, false
}

// Sort will sort the structures by name.
func (set *ReturnSetSet) Sort() {
	sort.Sort(set)
}

// Len get the number of structures in this set.
func (set *ReturnSetSet) Len() int {
	if set != nil {
		return len(set.Sets)
	}
	return 0
}

// Swap swaps the structures at the two given indices.
func (set *ReturnSetSet) Swap(aIndex, bIndex int) {
	if (set != nil) && (aIndex >= 0) && (bIndex >= 0) && (aIndex != bIndex) {
		if length := set.Len(); (aIndex < length) && (bIndex < length) {
			set.Sets[aIndex], set.Sets[bIndex] = set.Sets[bIndex], set.Sets[aIndex]
		}
	}
}

// Less determines if the structure's name at the first index
// is less than the structure's name at the second index.
func (set *ReturnSetSet) Less(aIndex, bIndex int) bool {
	aName, bName := "", ""
	if set != nil {
		if a := set.Sets[aIndex]; a != nil {
			aName = a.Name
		}
		if b := set.Sets[bIndex]; b != nil {
			bName = b.Name
		}
	}
	return aName < bName
}

// String gets the string of all the structures in this set.
func (set *ReturnSetSet) String() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, st := range set.Sets {
		parts[i] = st.String()
	}
	return strings.Join(parts, "\n")
}

// FullString gets the full string of all the structures in this set.
func (set *ReturnSetSet) FullString() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, st := range set.Sets {
		parts[i] = st.FullString()
	}
	return strings.Join(parts, "\n")
}
