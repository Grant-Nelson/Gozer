package reader

import (
	"golang.org/x/tools/go/packages"
)

// Project is the collection of packages that were parsed.
type Project struct {
	Packages []*packages.Package
}

// PreOrder returns all the packages with the imports after the package
// importing them. This will output via depth first traversal.
// The packages will only be outputted once, duplicates will be skipped.
// The packages are outputted in a consistent order.
func (p *Project) PreOrder() []*packages.Package {
	all := []*packages.Package{}
	packages.Visit(p.Packages, func(dep *packages.Package) bool {
		all = append(all, dep)
		return true
	}, nil)
	return all
}

// PostOrder returns all the packages with the imports before the package
// importing them. This will output via depth first traversal.
// The packages will only be outputted once, duplicates will be skipped.
// The packages are outputted in a consistent order.
func (p *Project) PostOrder() []*packages.Package {
	all := []*packages.Package{}
	packages.Visit(p.Packages, nil, func(dep *packages.Package) {
		all = append(all, dep)
	})
	return all
}

// Errors returns nil if there were no errors, a package error if here was
// only one, or a project error if there were multiple errors.
func (p *Project) Errors() error {
	pe := &ProjectErrors{}
	packages.Visit(p.Packages, nil, func(pkg *packages.Package) {
		pe.Errs = append(pe.Errs, pkg.Errors...)
	})

	switch pe.Count() {
	case 0:
		return nil
	case 1:
		return pe.Errs[0]
	default:
		return pe
	}
}
