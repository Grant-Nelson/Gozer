package types

import (
	"github.com/grant-nelson/Gozer/common"
)

var _ Type = (*ProgramType)(nil)

// ProgramType for storing the types of program.
type ProgramType struct {

	// Packages are the list of all packages used by this program.
	// The key is the path to the package.
	Packages *PackageSet
}

// Program creates a new program description.
func Program() *ProgramType {
	return &ProgramType{
		Packages: NewPackageSet(),
	}
}

// Contains checks if a package with the given name is in this program.
func (t *ProgramType) Contains(name string) bool {
	if t == nil {
		return false
	}
	_, exists := t.Packages.Find(name)
	return exists
}

// AddPackage adds a package to this program.
func (t *ProgramType) AddPackage(pack *PackageType) *ProgramType {
	if t != nil {
		t.Packages.Add(pack)
	}
	return t
}

// AddPackageWithShort adds a package to this program with a short name.
func (t *ProgramType) AddPackageWithShort(short string, pack *PackageType) *ProgramType {
	if t != nil {
		t.Packages.AddWithShort(short, pack)
	}
	return t
}

// String gets the name for this type.
func (t *ProgramType) String() string {
	if t == nil {
		return nilStr
	}
	if t.Packages.Len() <= 0 {
		return "{}"
	}
	return "{\n  " + common.Indent(t.Packages.String(), "  ") + "\n}"
}
