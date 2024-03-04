package cMethodSet

import (
	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/set"
)

func New() constructs.CMethodSet {
	return &methodSetImp{
		Set: set.New[constructs.CMethod](),
	}
}

type methodSetImp struct {
	collections.Set[constructs.CMethod]
}

func (imp *methodSetImp) TryGetByName(path string) (constructs.CMethod, bool) {
	return imp.Set.Enumerate().Where(getIsName(path)).First()
}

func getIsName(name string) collections.Predicate[constructs.CMethod] {
	return func(ci constructs.CMethod) bool { return ci.Name() == name }
}
