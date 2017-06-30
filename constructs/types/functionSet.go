package types

import (
	"sort"
	"strings"
)

// FunctionSet for storing a set of functions.
type FunctionSet struct {

	// Functions is the set of structures.
	Functions []*FunctionType
}

// NewFunctionSet creates a new set of functions.
func NewFunctionSet() *FunctionSet {
	return &FunctionSet{
		Functions: []*FunctionType{},
	}
}

// AddNew adds a function to this set.
// Returns the new function and true if new, false if already exists.
func (set *FunctionSet) AddNew(name string) (*FunctionType, bool) {
	if set == nil {
		return nil, false
	}
	if fn, found := set.Find(name); found {
		return fn, false
	}
	fn := Function()
	fn.Name = name
	set.Add(fn)
	return fn, true
}

// Add will append all non-nil functions to this set.
func (set *FunctionSet) Add(sts ...*FunctionType) *FunctionSet {
	if set != nil {
		for _, fn := range sts {
			if fn != nil {
				set.Functions = append(set.Functions, fn)
			}
		}
		set.Sort()
	}
	return set
}

// Find searches the set of functions to find a functions with the given name.
// If no functions by the given name is found, nil is returned.
func (set *FunctionSet) Find(name string) (*FunctionType, bool) {
	if set != nil {
		for _, fn := range set.Functions {
			if (fn != nil) && (fn.Name == name) {
				return fn, true
			}
		}
	}
	return nil, false
}

// Sort will sort the functions by name.
func (set *FunctionSet) Sort() {
	sort.Sort(set)
}

// Len get the number of functions in this set.
func (set *FunctionSet) Len() int {
	if set != nil {
		return len(set.Functions)
	}
	return 0
}

// Swap swaps the functions at the two given indices.
func (set *FunctionSet) Swap(aIndex, bIndex int) {
	if (set != nil) && (aIndex >= 0) && (bIndex >= 0) && (aIndex != bIndex) {
		if length := set.Len(); (aIndex < length) && (bIndex < length) {
			set.Functions[aIndex], set.Functions[bIndex] = set.Functions[bIndex], set.Functions[aIndex]
		}
	}
}

// Less determines if the functions's name at the first index
// is less than the functions's name at the second index.
func (set *FunctionSet) Less(aIndex, bIndex int) bool {
	aName, bName := "", ""
	if set != nil {
		if a := set.Functions[aIndex]; a != nil {
			aName = a.Name
		}
		if b := set.Functions[bIndex]; b != nil {
			bName = b.Name
		}
	}
	return aName < bName
}

// String gets the string of all the functions in this set.
func (set *FunctionSet) String() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, fn := range set.Functions {
		parts[i] = fn.String()
	}
	return strings.Join(parts, "\n")
}

// FullString gets the full string of all the functions in this set.
func (set *FunctionSet) FullString() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, fn := range set.Functions {
		parts[i] = fn.FullString()
	}
	return strings.Join(parts, "\n")
}

// FullBodyString gets the full body string of all the functions in this set.
func (set *FunctionSet) FullBodyString() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, fn := range set.Functions {
		parts[i] = fn.FullBodyString()
	}
	return strings.Join(parts, "\n")
}
