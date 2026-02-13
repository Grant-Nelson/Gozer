package project

import (
	"strings"

	"golang.org/x/tools/go/packages"
)

type Package struct {
	*packages.Package

	// TempTypeFile is set when the [types.Package] and related data was
	// serialized into a temporary file to be used as part of the cache.
	// This will be set to the path of that temporary file.
	TempTypeFile string
}

func (p *Package) IsTest() bool {
	return len(p.ForTest) > 0
}

func (p *Package) IsXTest() bool {
	return strings.HasSuffix(p.Name, `_test`)
}
