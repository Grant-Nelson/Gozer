package cache

import (
	"errors"
	"go/types"
	"io"

	"golang.org/x/tools/go/gcexportdata"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/parser"
)

// TODO: Need to add a manifest file to manage several different builds for the
// same package built for different build flags and versions. The build flags
// should be filtered to flags that actually affect the package.

type Cache struct {
	conv parser.SourceConverter
}

var (
	_ mods.Modifier        = (*Cache)(nil)
	_ mods.PackageStartExt = (*Cache)(nil)
	_ mods.PackageDoneExt  = (*Cache)(nil)
)

func New(conv parser.SourceConverter) *Cache {
	return &Cache{conv: conv}
}

func (c *Cache) ModName() string { return `Cache` }

func (c *Cache) PackageStart(pkg *project.Package, errGroup *faults.Group) (con bool, err error) {

	// TODO: Implement

	return true, errors.New(`unimplemented`)
}

func (c *Cache) PackageDone(pkg *project.Package, errGroup *faults.Group) (con bool, err error) {

	// TODO: Implement

	return true, errors.New(`unimplemented`)
}

var _ = readPackage  // TODO: Remove, keeps readPackage from being "unused"
var _ = writePackage // TODO: Remove, keeps writePackage from being "unused"

func readPackage(in io.Reader, pkg *project.Package, errGroup *faults.Group) error {
	imports := make(map[string]*types.Package, len(pkg.Imports))
	for path, pkg := range pkg.Imports {
		imports[path] = pkg.Types
	}

	tPkg, err := gcexportdata.Read(in, pkg.Fset, imports, pkg.PkgPath)
	if err = errGroup.Add(err); err != nil {
		return err
	}

	pkg.Types = tPkg
	return nil
}

func writePackage(out io.Writer, pkg *project.Package, errGroup *faults.Group) error {
	err := gcexportdata.Write(out, pkg.Fset, pkg.Types)
	return errGroup.Add(err)
}
