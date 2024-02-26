package packageSet

import (
	"go/build"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/dictionary"
)

type PackageSet interface {
	Context() build.Context

	Add(mainPackageDirPath string) error

	Packages() collections.ReadonlyDictionary[string, *build.Package]
}

func New(context build.Context) PackageSet {
	return &packageSetImp{
		context:  context,
		packages: dictionary.New[string, *build.Package](),
	}
}
