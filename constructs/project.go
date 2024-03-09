package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type IProject interface {
	INamed
	Packages() collections.List[IPackage]
	AtPath(path string) (IPackage, bool)
	_projectConstruct()
}

type projectImp struct {
	namedImp
	packages collections.List[IPackage]
}

func (imp *projectImp) _projectConstruct() {}

func (imp *projectImp) Packages() collections.List[IPackage] {
	return imp.packages
}

func (imp *projectImp) AtPath(path string) (IPackage, bool) {
	return imp.packages.Enumerate().
		Where(func(p IPackage) bool { return p.Path() == path }).
		First()
}

func NewProject(name string) IProject {
	return &projectImp{
		namedImp: newName(name),
		packages: list.New[IPackage](),
	}
}
