package types

import (
	"sort"
	"strings"
)

// ClassSet for storing a set of classes.
type ClassSet struct {

	// Classes is the set of classes.
	Classes []*ClassType
}

// NewClassSet creates a new set of classes.
func NewClassSet() *ClassSet {
	return &ClassSet{
		Classes: []*ClassType{},
	}
}

// AddNew adds a class to this set.
// Returns the new class and true if new, false if already exists.
func (set *ClassSet) AddNew(name string) (*ClassType, bool) {
	if set == nil {
		return nil, false
	}
	if class, found := set.Find(name); found {
		return class, false
	}
	class := Class()
	class.Name = name
	set.Add(class)
	return class, true
}

// Add will append all non-nil classes to this set.
func (set *ClassSet) Add(classes ...*ClassType) *ClassSet {
	if set != nil {
		for _, class := range classes {
			if class != nil {
				set.Classes = append(set.Classes, class)
			}
		}
		set.Sort()
	}
	return set
}

// Find searches the set of classes to find a class with the given name.
// If no class by the given name is found, nil is returned.
func (set *ClassSet) Find(name string) (*ClassType, bool) {
	if set != nil {
		for _, class := range set.Classes {
			if (class != nil) && (class.Name == name) {
				return class, true
			}
		}
	}
	return nil, false
}

// Sort will sort the classes by name.
func (set *ClassSet) Sort() {
	sort.Sort(set)
}

// Len get the number of classes in this set.
func (set *ClassSet) Len() int {
	if set != nil {
		return len(set.Classes)
	}
	return 0
}

// Swap swaps the classes at the two given indices.
func (set *ClassSet) Swap(aIndex, bIndex int) {
	if (set != nil) && (aIndex >= 0) && (bIndex >= 0) && (aIndex != bIndex) {
		if length := set.Len(); (aIndex < length) && (bIndex < length) {
			set.Classes[aIndex], set.Classes[bIndex] = set.Classes[bIndex], set.Classes[aIndex]
		}
	}
}

// Less determines if the class's name at the first index
// is less than the class's name at the second index.
func (set *ClassSet) Less(aIndex, bIndex int) bool {
	aName, bName := "", ""
	if set != nil {
		if a := set.Classes[aIndex]; a != nil {
			aName = a.Name
		}
		if b := set.Classes[bIndex]; b != nil {
			bName = b.Name
		}
	}
	return aName < bName
}

// String gets the string of all the classes in this set.
func (set *ClassSet) String() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, class := range set.Classes {
		parts[i] = class.String()
	}
	return strings.Join(parts, "\n")
}

// FullString gets the full string of all the classes in this set.
func (set *ClassSet) FullString() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, class := range set.Classes {
		parts[i] = class.FullString()
	}
	return strings.Join(parts, "\n")
}
