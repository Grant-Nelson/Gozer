package cache

import (
	"errors"
	"go/types"
	"io"
	"os"

	"golang.org/x/tools/go/gcexportdata"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/avail/source"
	"github.com/Grant-Nelson/Gozer/compiler/loader/mods"
	"github.com/Grant-Nelson/Gozer/compiler/project"
	"github.com/Grant-Nelson/Gozer/compiler/project/enums/buildState"
)

type Config struct {
	Build        []string
	Converter    source.Converter
	Logger       *logger.Logger
	ErrGroup     *faults.ErrGroup
	DisableRead  bool
	DisableWrite bool
}

type Cache struct {
	build        []string
	conv         source.Converter
	logger       *logger.Logger
	errGroup     *faults.ErrGroup
	disableRead  bool
	disableWrite bool
}

type cacheMod struct {
	pkg      *project.Package
	logger   *logger.Logger
	errGroup *faults.ErrGroup
}

var (
	_ mods.ModFactory = (*Cache)(nil)
	_ mods.Modifier   = (*cacheMod)(nil)
)

func New(cfg *Config) *Cache {
	return &Cache{
		build:        cfg.Build,
		conv:         cfg.Converter,
		logger:       cfg.Logger,
		errGroup:     cfg.ErrGroup,
		disableRead:  cfg.DisableRead,
		disableWrite: cfg.DisableWrite,
	}
}

func (c *Cache) StartPackage(pkg *project.Package) (con bool, mod mods.Modifier, err error) {
	defer c.errGroup.Recover(&err)

	if !c.disableWrite {
		mod = &cacheMod{
			pkg:      pkg,
			logger:   c.logger,
			errGroup: c.errGroup,
		}
	}

	if c.disableRead {
		return true, mod, nil
	}

	ok, path, data, err := c.conv(pkg.Ast.Dir, nil)
	if !ok || err != nil {
		return true, mod, c.errGroup.Add(err)
	}

	// TODO: See TODOs in README.md

	reader, err := source.ToReader(path, data)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Cache miss
			c.logger.Printf("Cache miss for %s\n", pkg.PkgPath())
			return true, mod, nil
		}
		return true, mod, c.errGroup.Add(err)
	}

	if err = readPackage(reader, pkg); err != nil {
		return true, mod, c.errGroup.Add(err)
	}

	// Cache hit so skip rest of loading.
	c.logger.Printf("Cache hit for %s\n", pkg.PkgPath())
	pkg.State = buildState.Finished
	return false, nil, nil
}

func (c *cacheMod) PackageDone() (con bool, err error) {
	defer c.errGroup.Recover(&err)

	f, err := os.CreateTemp(``, `gozerTypePkg*.a`)
	if err != nil {
		return true, c.errGroup.Add(err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = c.errGroup.Add(closeErr)
		}
	}()

	if err = writePackage(f, c.pkg); err != nil {
		return true, c.errGroup.Add(err)
	}

	c.logger.Printf("Cache written for %s\n", c.pkg.PkgPath())
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
