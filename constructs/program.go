package constructs

import "strings"

var _ Type = (*ProgramType)(nil)

// ProgramType for storing the types of program.
type ProgramType struct {

	// Packages are the list of all packages used by this program.
	// The key is the path to the package.
	Packages map[string]*PackageType
}

// Program creates a new program description.
func Program() *ProgramType {
	return &ProgramType{
		Packages: map[string]*PackageType{},
	}
}

// Contains checks if a package with the given name is in this program.
func (t *ProgramType) Contains(name string) bool {
	_, exists := t.Packages[name]
	return exists
}

// AddPackage adds a package to this program.
func (t *ProgramType) AddPackage(name string, pack *PackageType) *ProgramType {
	t.Packages[name] = pack
	return t
}

// Equals determins if the given type is the same as this type.
func (t *ProgramType) Equals(other Type) bool {
	if t == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	tother, ok := other.(*ProgramType)
	if !ok {
		return false
	}
	if len(t.Packages) != len(tother.Packages) {
		return false
	}
	for name, pack := range t.Packages {
		tpack, exists := tother.Packages[name]
		if !exists {
			return false
		}
		if !pack.Equals(tpack) {
			return false
		}
	}
	return true
}

// String gets the name for this type.
func (t *ProgramType) String() string {
	if t == nil {
		return nilStr
	}
	parts := make([]string, len(t.Packages))
	i := 0
	for name := range t.Packages {
		parts[i] = "import " + name
		i++
	}
	return "{\n" +
		indent(strings.Join(parts, "\n"), "  ") + "\n" +
		"}"
}
