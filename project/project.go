package project

import (
	"go/token"

	"golang.org/x/tools/go/packages"
)

type Project struct {
	fSet     *token.FileSet
	packages []*packages.Package
}
