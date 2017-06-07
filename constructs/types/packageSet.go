package types

import (
	"sort"
	"strings"
)

// PackageSet for storing a set of packages.
type PackageSet struct {

	// Packages is the set of packages.
	Packages []*PackageType
}

// NewPackageSet creates a new set of packages.
func NewPackageSet() *PackageSet {
	return &PackageSet{
		Packages: []*PackageType{},
	}
}

// AddNew adds a package to this set.
// Returns the new package and true if new, false if already exists.
func (set *PackageSet) AddNew(name string) (*PackageType, bool) {
	if set == nil {
		return nil, false
	}
	if pack, found := set.Find(name); found {
		return pack, false
	}
	pack := Package()
	pack.Name = name
	set.Add(pack)
	return pack, true
}

// Add will append all non-nil packages to this set.
func (set *PackageSet) Add(packages ...*PackageType) {
	if set != nil {
		for _, pack := range packages {
			if pack != nil {
				set.Packages = append(set.Packages, pack)
			}
		}
		set.Sort()
	}
}

// Find searches the set of packages to find a package with the given name.
// If no package by the given name is found, nil is returned.
func (set *PackageSet) Find(name string) (*PackageType, bool) {
	if set != nil {
		for _, pack := range set.Packages {
			if (pack != nil) && (pack.Name == name) {
				return pack, true
			}
		}
	}
	return nil, false
}

// Sort will sort the packages by name.
func (set *PackageSet) Sort() {
	sort.Sort(set)
}

// Len get the number of packages in this set.
func (set *PackageSet) Len() int {
	if set != nil {
		return len(set.Packages)
	}
	return 0
}

// Swap swaps the packages at the two given indices.
func (set *PackageSet) Swap(aIndex, bIndex int) {
	if (set != nil) && (aIndex >= 0) && (bIndex >= 0) && (aIndex != bIndex) {
		if length := set.Len(); (aIndex < length) && (bIndex < length) {
			set.Packages[aIndex], set.Packages[bIndex] = set.Packages[bIndex], set.Packages[aIndex]
		}
	}
}

// Less determines if the package's name at the first index
// is less than the package's name at the second index.
func (set *PackageSet) Less(aIndex, bIndex int) bool {
	aName, bName := "", ""
	if set != nil {
		if a := set.Packages[aIndex]; a != nil {
			aName = a.Name
		}
		if b := set.Packages[bIndex]; b != nil {
			bName = b.Name
		}
	}
	return aName < bName
}

// String gets the string of all the packages in this set.
func (set *PackageSet) String() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, pack := range set.Packages {
		parts[i] = pack.String()
	}
	return strings.Join(parts, "\n")
}

// FullString gets the full string of all the packages in this set.
func (set *PackageSet) FullString() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, pack := range set.Packages {
		parts[i] = pack.FullString()
	}
	return strings.Join(parts, "\n")
}
