package cPackageSet

import (
	"github.com/Snow-Gremlin/Gozer/constructs"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/set"
)

func New() constructs.CPackageSet {
	return &packageSetImp{
		Set: set.New[constructs.CPackage](),
	}
}

type packageSetImp struct {
	collections.Set[constructs.CPackage]
}

func (imp *packageSetImp) TryGetByPath(path string) (constructs.CPackage, bool) {
	return imp.Set.Enumerate().Where(getIsPath(path)).First()
}

func getIsPath(path string) collections.Predicate[constructs.CPackage] {
	return func(ci constructs.CPackage) bool { return ci.Path() == path }
}
