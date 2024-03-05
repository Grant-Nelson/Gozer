package cObjectSet

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/set"

	"github.com/Snow-Gremlin/Gozer/constructs"
)

func New() constructs.CObjectSet {
	return &objectSetImp{
		Set: set.New[constructs.CObject](),
	}
}

type objectSetImp struct {
	collections.Set[constructs.CObject]
}

func (imp *objectSetImp) TryGetByName(name string) (constructs.CObject, bool) {
	return imp.Set.Enumerate().Where(getIsName(name)).First()
}

func getIsName(name string) collections.Predicate[constructs.CObject] {
	return func(ci constructs.CObject) bool { return ci.Name() == name }
}
