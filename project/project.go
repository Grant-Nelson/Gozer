package project

import (
	"go/token"

	"golang.org/x/tools/go/packages"
)

type Project struct {
	fSet     *token.FileSet
	packages []*packages.Package
}

func New(fSet *token.FileSet, packages []*packages.Package) *Project {
	return &Project{fSet: fSet, packages: packages}
}
