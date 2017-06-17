package types

import (
	"sort"
	"strings"
)

// DeclarationSet for storing a set of declarations.
type DeclarationSet struct {

	// Declarations is the set of declarations.
	Declarations []*DeclarationType
}

// NewDeclarationSet creates a new set of declarations.
func NewDeclarationSet() *DeclarationSet {
	return &DeclarationSet{
		Declarations: []*DeclarationType{},
	}
}

// AddNew adds a declaration to this set.
// Returns the new declaration and true if new, false if already exists.
func (set *DeclarationSet) AddNew(name string, data Type) (*DeclarationType, bool) {
	if set == nil {
		return nil, false
	}
	if class, found := set.Find(name); found {
		return class, false
	}
	inter := Declaration()
	inter.Name = name
	inter.Data = data
	set.Add(inter)
	return inter, true
}

// Add will append all non-nil declarations to this set.
func (set *DeclarationSet) Add(inters ...*DeclarationType) *DeclarationSet {
	if set != nil {
		for _, inter := range inters {
			if inter != nil {
				set.Declarations = append(set.Declarations, inter)
			}
		}
		set.Sort()
	}
	return set
}

// Find searches the set of declarations to find a declaration with the given name.
// If no declaration by the given name is found, nil is returned.
func (set *DeclarationSet) Find(name string) (*DeclarationType, bool) {
	if set != nil {
		for _, inter := range set.Declarations {
			if (inter != nil) && (inter.Name == name) {
				return inter, true
			}
		}
	}
	return nil, false
}

// Sort will sort the declarations by name.
func (set *DeclarationSet) Sort() {
	sort.Sort(set)
}

// Len get the number of declarations in this set.
func (set *DeclarationSet) Len() int {
	if set != nil {
		return len(set.Declarations)
	}
	return 0
}

// Swap swaps the declarations at the two given indices.
func (set *DeclarationSet) Swap(aIndex, bIndex int) {
	if (set != nil) && (aIndex >= 0) && (bIndex >= 0) && (aIndex != bIndex) {
		if length := set.Len(); (aIndex < length) && (bIndex < length) {
			set.Declarations[aIndex], set.Declarations[bIndex] = set.Declarations[bIndex], set.Declarations[aIndex]
		}
	}
}

// Less determines if the declaration's name at the first index
// is less than the declaration's name at the second index.
func (set *DeclarationSet) Less(aIndex, bIndex int) bool {
	aName, bName := "", ""
	if set != nil {
		if a := set.Declarations[aIndex]; a != nil {
			aName = a.Name
		}
		if b := set.Declarations[bIndex]; b != nil {
			bName = b.Name
		}
	}
	return aName < bName
}

// String gets the string of all the declarations in this set.
func (set *DeclarationSet) String() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, inter := range set.Declarations {
		parts[i] = inter.String()
	}
	return strings.Join(parts, "\n")
}

// FullString gets the full string of all the declarations in this set.
func (set *DeclarationSet) FullString() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, inter := range set.Declarations {
		parts[i] = inter.FullString()
	}
	return strings.Join(parts, "\n")
}
