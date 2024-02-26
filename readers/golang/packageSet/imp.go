package packageSet

import (
	"go/build"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/set"
)

type packageSetImp struct {
	context  build.Context
	packages collections.Dictionary[string, *build.Package]
}

func (ps *packageSetImp) Context() build.Context {
	return ps.context
}

func (ps *packageSetImp) Add(mainPackageDirPath string) error {
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

func (ps *packageSetImp) Packages() collections.ReadonlyDictionary[string, *build.Package] {
	return ps.packages
}

func takePendingImport(pending collections.Set[string]) (string, bool) {
	path, hadPending := pending.Enumerate().First()
	if !hadPending {
		return ``, false
	}
	pending.Remove(path)
	return path, true
}

func (ps *packageSetImp) addPackage(path string) (*build.Package, error) {
	pkg, err := ps.context.ImportDir(path, build.FindOnly)
	if err != nil {
		return nil, err
	}
	ps.packages.Add(path, pkg)
	return pkg, nil
}

func (ps *packageSetImp) addPendingImport(pending collections.Set[string], imports []string) {
	pending.AddFrom(enumerator.Enumerate(imports...).WhereNot(ps.packages.Contains))
}
