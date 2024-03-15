package reader

import (
	"golang.org/x/tools/go/packages"
)

type Project struct {
	Packages []*packages.Package
}

func (p *Project) PreOrder() []*packages.Package {
	all := []*packages.Package{}
	packages.Visit(p.Packages, func(dep *packages.Package) bool {
		all = append(all, dep)
		return true
	}, nil)
	return all
}

func (p *Project) PostOrder() []*packages.Package {
	all := []*packages.Package{}
	packages.Visit(p.Packages, nil, func(dep *packages.Package) {
		all = append(all, dep)
	})
	return all
}

func (p *Project) Errors() []error {
	errs := []error{}
	packages.Visit(p.Packages, nil, func(pkg *packages.Package) {
		for _, err := range pkg.Errors {
			errs = append(errs, err)
		}
	})
	return errs
}
