package cache

import (
	"errors"
	"go/types"
	"io"
	"os"

	"golang.org/x/tools/go/gcexportdata"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/source"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
)

type Config struct {
	Build     []string
	Converter source.Converter
}

type Cache struct {
	build []string
	conv  source.Converter
}

var (
	_ mods.Modifier        = (*Cache)(nil)
	_ mods.PackageStartExt = (*Cache)(nil)
	_ mods.PackageDoneExt  = (*Cache)(nil)
)

func New(cfg *Config) *Cache {
	return &Cache{
		build: cfg.Build,
		conv:  cfg.Converter,
	}
}

func (c *Cache) ModName() string { return `Cache` }

func (c *Cache) PackageStart(pkg *project.Package, errGroup *faults.Group) (bool, error) {
	ok, path, data, err := c.conv(pkg.Ast.Dir, nil)
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
	imports := make(map[string]*types.Package, len(pkg.Ast.Imports))
	for path, pkg := range pkg.Ast.Imports {
		imports[path] = pkg.Types
	}

	tPkg, err := gcexportdata.Read(in, pkg.Ast.Fset, imports, pkg.Ast.PkgPath)
	if err = errGroup.Add(err); err != nil {
		return err
	}

	pkg.Ast.Types = tPkg
	return nil
}

func writePackage(out io.Writer, pkg *project.Package, errGroup *faults.Group) error {
	err := gcexportdata.Write(out, pkg.Ast.Fset, pkg.Ast.Types)
	return errGroup.Add(err)
}
