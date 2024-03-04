package cProject

import (
	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/constructs/cPackageSet"
)

func New() constructs.CProject {
	return &projectImp{
		name:     `unnamed`,
		packages: cPackageSet.New(),
	}
}

type projectImp struct {
	name     string
	packages constructs.CPackageSet
}

func (imp *projectImp) Name() string {
	return imp.name
}

func (imp *projectImp) SetName(name string) {
	imp.name = name
}

func (imp *projectImp) Packages() constructs.CPackageSet {
	return imp.packages
}
