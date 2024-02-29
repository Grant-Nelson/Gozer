package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/set"
)

type CPackageSet struct {
	collections.Set[*CPackage]
}

func NewPackageSet() *CPackageSet {
	return &CPackageSet{
		Set: set.New[*CPackage](),
	}
}

func (p *CPackageSet) TryGet(path string) (*CPackage, bool) {
	return p.Set.Enumerate().Where(p.getIsPath(path)).First()
}

func (p *CPackageSet) getIsPath(path string) collections.Predicate[*CPackage] {
	return func(ci *CPackage) bool { return ci.Path == path }
}
