package types

import (
	"sort"
	"strings"
)

// InterfaceSet for storing a set of interfaces.
type InterfaceSet struct {

	// Interfaces is the set of interfaces.
	Interfaces []*InterfaceType
}

// NewInterfaceSet creates a new set of interfaces.
func NewInterfaceSet() *InterfaceSet {
	return &InterfaceSet{
		Interfaces: []*InterfaceType{},
	}
}

// AddNew adds a interface to this set.
// Returns the new interface and true if new, false if already exists.
func (set *InterfaceSet) AddNew(name string) (*InterfaceType, bool) {
	if set == nil {
		return nil, false
	}
	if class, found := set.Find(name); found {
		return class, false
	}
	inter := Interface()
	inter.Name = name
	set.Add(inter)
	return inter, true
}

// Add will append all non-nil interfaces to this set.
func (set *InterfaceSet) Add(inters ...*InterfaceType) *InterfaceSet {
	if set != nil {
		for _, inter := range inters {
			if inter != nil {
				set.Interfaces = append(set.Interfaces, inter)
			}
		}
		set.Sort()
	}
	return set
}

// Find searches the set of interfaces to find a interface with the given name.
// If no interface by the given name is found, nil is returned.
func (set *InterfaceSet) Find(name string) (*InterfaceType, bool) {
	if set != nil {
		for _, inter := range set.Interfaces {
			if (inter != nil) && (inter.Name == name) {
				return inter, true
			}
		}
	}
	return nil, false
}

// Sort will sort the interfaces by name.
func (set *InterfaceSet) Sort() {
	sort.Sort(set)
}

// Len get the number of interfaces in this set.
func (set *InterfaceSet) Len() int {
	if set != nil {
		return len(set.Interfaces)
	}
	return 0
}

// Swap swaps the interfaces at the two given indices.
func (set *InterfaceSet) Swap(aIndex, bIndex int) {
	if (set != nil) && (aIndex >= 0) && (bIndex >= 0) && (aIndex != bIndex) {
		if length := set.Len(); (aIndex < length) && (bIndex < length) {
			set.Interfaces[aIndex], set.Interfaces[bIndex] = set.Interfaces[bIndex], set.Interfaces[aIndex]
		}
	}
}

// Less determines if the interface's name at the first index
// is less than the interface's name at the second index.
func (set *InterfaceSet) Less(aIndex, bIndex int) bool {
	aName, bName := "", ""
	if set != nil {
		if a := set.Interfaces[aIndex]; a != nil {
			aName = a.Name
		}
		if b := set.Interfaces[bIndex]; b != nil {
			bName = b.Name
		}
	}
	return aName < bName
}

// String gets the string of all the interfaces in this set.
func (set *InterfaceSet) String() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, inter := range set.Interfaces {
		parts[i] = inter.String()
	}
	return strings.Join(parts, "\n")
}

// FullString gets the full string of all the interfaces in this set.
func (set *InterfaceSet) FullString() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, inter := range set.Interfaces {
		parts[i] = inter.FullString()
	}
	return strings.Join(parts, "\n")
}
