package cache

import (
	"errors"
	"go/types"
	"io"
	"os"

	"golang.org/x/tools/go/gcexportdata"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/source"
)

// TODO: Need to add a manifest file to manage several different builds for the
// same package built for different build flags and versions. The build flags
// should be filtered to flags that actually affect the package.

// TODO: Also note that the stored cache needs to be put into a temporary location
// until the resulting transpiled code has been written, then the stored cache
// and transpiled code can be moved to the location caches are read from.

type Cache struct {
	conv source.Converter
}

var (
	_ mods.Modifier        = (*Cache)(nil)
	_ mods.PackageStartExt = (*Cache)(nil)
	_ mods.PackageDoneExt  = (*Cache)(nil)
)

func New(conv source.Converter) *Cache {
	return &Cache{conv: conv}
}

func (c *Cache) ModName() string { return `Cache` }

func (c *Cache) PackageStart(pkg *project.Package, errGroup *faults.Group) (bool, error) {
	ok, path, data, err := c.conv(pkg.Dir, nil)
	if !ok || err != nil {
		return true, errGroup.Add(err)
	}

	reader, err := source.ToReader(path, data)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Cache miss
			return true, nil
		}
		return true, errGroup.Add(err)
	}

	if err = readPackage(reader, pkg, errGroup); err != nil {
		return true, errGroup.Add(err)
	}

	// Cache hit so skip rest of loading.
	return false, nil
}

func (c *Cache) PackageDone(pkg *project.Package, errGroup *faults.Group) (con bool, err error) {
	f, err := os.CreateTemp(``, `gozerTypePkg*.a`)
	if err != nil {
		return true, errGroup.Add(err)
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = errGroup.Add(closeErr)
		}
	}()

	if err = writePackage(f, pkg, errGroup); err != nil {
		return true, errGroup.Add(err)
	}

	pkg.TempTypeFile = f.Name()
	return true, nil
}

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
