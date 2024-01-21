package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type CProject struct {
	Name     string
	Packages collections.List[*CPackage]
}

func NewProject(name string) *CProject {
	return &CProject{
		Name:     name,
		Packages: list.New[*CPackage](),
	}
}
