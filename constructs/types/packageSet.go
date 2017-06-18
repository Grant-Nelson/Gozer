package types

import (
	"sort"
	"strings"
)

// PackageSet for storing a set of packages.
type PackageSet struct {

	// shorts is the short name for an import.
	// If the short name isn't given the entry will be an empty string.
	shorts []string

	// packages is the set of packages.
	packages []*PackageType
}

// NewPackageSet creates a new set of packages.
func NewPackageSet() *PackageSet {
	return &PackageSet{
		shorts:   []string{},
		packages: []*PackageType{},
	}
}

// Shorts gets the set of short names for packages.
func (set *PackageSet) Shorts() []string {
	if set == nil {
		return []string{}
	}
	return set.shorts
}

// Packages gets the set of packages.
func (set *PackageSet) Packages() []*PackageType {
	if set == nil {
		return []*PackageType{}
	}
	return set.packages
}

// SetShort
func (set *PackageSet) SetShort(short string, name string) bool {
	if set != nil {
		for i, pack := range set.packages {
			if (pack != nil) && (pack.Name == name) {
				set.shorts[i] = short
				return true
			}
		}
	}
	return false
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
	set.AddWithShort("", pack)
	return pack, true
}

// Add will append all non-nil packages to this set.
func (set *PackageSet) Add(packages ...*PackageType) *PackageSet {
	if set != nil {
		for _, pack := range packages {
			if pack != nil {
				set.shorts = append(set.shorts, "")
				set.packages = append(set.packages, pack)
			}
		}
		set.Sort()
	}
	return set
}

// Add will append all non-nil packages to this set.
func (set *PackageSet) AddWithShort(short string, pack *PackageType) {
	set.shorts = append(set.shorts, short)
	set.packages = append(set.packages, pack)
	set.Sort()
}

// Find searches the set of packages to find a package with the given name.
// If no package by the given name is found, nil is returned.
func (set *PackageSet) Find(name string) (*PackageType, bool) {
	if set != nil {
		for i, pack := range set.packages {
			if (pack != nil) && ((pack.Name == name) || (set.shorts[i] == name)) {
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
		return len(set.packages)
	}
	return 0
}

// Swap swaps the packages at the two given indices.
func (set *PackageSet) Swap(aIndex, bIndex int) {
	if (set != nil) && (aIndex >= 0) && (bIndex >= 0) && (aIndex != bIndex) {
		if length := set.Len(); (aIndex < length) && (bIndex < length) {
			set.shorts[aIndex], set.shorts[bIndex] = set.shorts[bIndex], set.shorts[aIndex]
			set.packages[aIndex], set.packages[bIndex] = set.packages[bIndex], set.packages[aIndex]
		}
	}
}

// Less determines if the package's name at the first index
// is less than the package's name at the second index.
func (set *PackageSet) Less(aIndex, bIndex int) bool {
	aName, bName := "", ""
	if set != nil {
		if aShort := set.shorts[aIndex]; len(aShort) > 0 {
			aName = aShort
		} else if a := set.packages[aIndex]; a != nil {
			aName = a.Name
		}
		if bShort := set.shorts[bIndex]; len(bShort) > 0 {
			bName = bShort
		} else if b := set.packages[bIndex]; b != nil {
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
	for i, pack := range set.packages {
		parts[i] = pack.StringWithShort(set.shorts[i])
	}
	return strings.Join(parts, "\n")
}

// FullString gets the full string of all the packages in this set.
func (set *PackageSet) FullString() string {
	if set == nil {
		return nilStr
	}
	parts := make([]string, set.Len())
	for i, pack := range set.packages {
		parts[i] = pack.FullStringWithShort(set.shorts[i])
	}
	return strings.Join(parts, "\n")
}
