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
	"github.com/Grant-Nelson/Gozer/project/enums/buildState"
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

type cacheMod struct {
	pkg      *project.Package
	errGroup *faults.Group
}

var (
	_ mods.ModFactory = (*Cache)(nil)
	_ mods.Modifier   = (*cacheMod)(nil)
)

func New(cfg *Config) *Cache {
	return &Cache{
		build: cfg.Build,
		conv:  cfg.Converter,
	}
}

func (c *Cache) StartPackage(pkg *project.Package, errGroup *faults.Group) (bool, mods.Modifier, error) {
	cm := &cacheMod{
		pkg:      pkg,
		errGroup: errGroup,
	}

	ok, path, data, err := c.conv(pkg.Ast.Dir, nil)
	if !ok || err != nil {
		return true, cm, errGroup.Add(err)
	}

	// TODO: See TODOs in README.md

	reader, err := source.ToReader(path, data)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Cache miss
			return true, cm, nil
		}
		return true, cm, errGroup.Add(err)
	}

	if err = readPackage(reader, pkg); err != nil {
		return true, cm, errGroup.Add(err)
	}

	// Cache hit so skip rest of loading.
	pkg.State = buildState.Finished
	return false, nil, nil
}

func (c *cacheMod) PackageDone() (con bool, err error) {
	f, err := os.CreateTemp(``, `gozerTypePkg*.a`)
	if err != nil {
		return true, c.errGroup.Add(err)
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = c.errGroup.Add(closeErr)
		}
	}()

	if err = writePackage(f, c.pkg); err != nil {
		return true, c.errGroup.Add(err)
	}

	c.pkg.TempTypeFile = f.Name()
	return true, nil
}

func readPackage(in io.Reader, pkg *project.Package) error {
	imports := make(map[string]*types.Package, len(pkg.Ast.Imports))
	for path, pkg := range pkg.Ast.Imports {
		imports[path] = pkg.Types
	}

	tPkg, err := gcexportdata.Read(in, pkg.Ast.Fset, imports, pkg.Ast.PkgPath)
	if err != nil {
		return err
	}

	pkg.Ast.Types = tPkg
	return nil
}

func writePackage(out io.Writer, pkg *project.Package) error {
	return gcexportdata.Write(out, pkg.Ast.Fset, pkg.Ast.Types)
}
