package packageSet

import (
	"go/build"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/dictionary"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/set"
)

type packageSet struct {
	context     build.Context
	mainPackage *build.Package
	packages    collections.Dictionary[string, *build.Package]
}

func newPackageSet() *packageSet {
	return &packageSet{
		context:     build.Default,
		mainPackage: nil,
		packages:    dictionary.New[string, *build.Package](),
	}
}

func (ps *packageSet) loadAllPackages(mainPackageDirPath string) error {
	pending := set.New[string]()
	path, hadPending := mainPackageDirPath, true
	for hadPending {
		pkg, err := ps.addPackage(path)
		if err != nil {
			return err
		}

		ps.addPendingImport(pending, pkg.Imports)
		ps.addPendingImport(pending, pkg.TestImports)
		ps.addPendingImport(pending, pkg.XTestImports)
		path, hadPending = takePendingImport(pending)
	}
	return nil
}

func takePendingImport(pending collections.Set[string]) (string, bool) {
	path, hadPending := pending.Enumerate().First()
	if !hadPending {
		return ``, false
	}
	pending.Remove(path)
	return path, true
}

func (ps *packageSet) addPackage(path string) (*build.Package, error) {
	pkg, err := ps.context.ImportDir(path, build.FindOnly)
	if err != nil {
		return nil, err
	}

	if ps.mainPackage == nil {
		ps.mainPackage = pkg
	}
	ps.packages.Add(path, pkg)

	return pkg, nil
}

func (ps *packageSet) addPendingImport(pending collections.Set[string], imports []string) {
	pending.AddFrom(enumerator.Enumerate(imports...).WhereNot(ps.packages.Contains))
}

func (ps *packageSet) Context() build.Context {
	return ps.context
}

func (ps *packageSet) MainPackage() *build.Package {
	return ps.mainPackage
}

func (ps *packageSet) Packages() collections.ReadonlyDictionary[string, *build.Package] {
	return ps.packages
}
