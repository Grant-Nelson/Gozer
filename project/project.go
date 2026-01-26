package project

import (
	"go/token"

	"golang.org/x/tools/go/packages"
)

type Project struct {
	FileSet  *token.FileSet
	Packages []*packages.Package
}

func New(fSet *token.FileSet, packages []*packages.Package) *Project {
	return &Project{FileSet: fSet, Packages: packages}
}
