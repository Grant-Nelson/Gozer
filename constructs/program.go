package constructs

import (
	"sort"
	"strings"

	"github.com/grant-nelson/Gozer/common"
)

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

// String gets the name for this type.
func (t *ProgramType) String() string {
	if t == nil {
		return nilStr
	}
	if len(t.Packages) <= 0 {
		return "{}"
	}
	i := 0
	parts := make([]string, len(t.Packages))
	for name := range t.Packages {
		parts[i] = "import " + name
		i++
	}
	sort.Strings(parts)
	return "{\n" +
		"  " + common.Indent(strings.Join(parts, "\n"), "  ") + "\n" +
		"}"
}
