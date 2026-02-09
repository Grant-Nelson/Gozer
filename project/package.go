package project

import (
	"strings"

	"golang.org/x/tools/go/packages"
)

type Package struct {
	*packages.Package
}

func (p *Package) IsTest() bool {
	return len(p.ForTest) > 0
}

func (p *Package) IsXTest() bool {
	return strings.HasSuffix(p.Name, `_test`)
}
