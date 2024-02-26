package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/set"
)

type CPackage struct {
	Name    string
	Path    string
	Imports collections.Set[*CImport]
}

func NewPackage(name, path string) *CPackage {
	return &CPackage{
		Name:    name,
		Path:    path,
		Imports: set.New[*CImport](),
	}
}

func getIsPath(path string) collections.Predicate[*CImport] {
	return func(ci *CImport) bool { return ci.Path == path }
}

func (p *CPackage) ImportForPath(path string) *CImport {
	i, _ := p.Imports.Enumerate().Where(getIsPath(path)).First()
	return i
}
