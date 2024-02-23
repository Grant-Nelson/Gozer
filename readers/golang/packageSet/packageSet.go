package packageSet

import (
	"go/build"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

type PackageSet interface {
	Context() build.Context

	MainPackage() *build.Package

	Packages() collections.ReadonlyDictionary[string, *build.Package]
}

func New(mainPackageDirPath string) (PackageSet, error) {
	ps := newPackageSet()
	if err := ps.loadAllPackages(mainPackageDirPath); err != nil {
		return nil, err
	}
	return ps, nil
}
