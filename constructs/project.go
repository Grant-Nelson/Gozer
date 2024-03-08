package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type IProject interface {
	INamed
	Packages() collections.List[IPackage]
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

func NewProject(name string) IProject {
	return &projectImp{
		namedImp: newName(name),
		packages: list.New[IPackage](),
	}
}
