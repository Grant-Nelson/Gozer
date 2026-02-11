package typeChecker

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
)

var (
	defaultSize = &types.StdSizes{WordSize: 4, MaxAlign: 8}
)

type TypeChecker struct {
	goVersion string
	ctx       *types.Context
	sizes     types.Sizes
}

var (
	_ mods.Modifier       = (*TypeChecker)(nil)
	_ mods.PackageDoneExt = (*TypeChecker)(nil)
	_ types.Importer      = (*importer)(nil)
)

func New(goVersion string, ctx *types.Context, sizes types.Sizes) *TypeChecker {
	if sizes == nil {
		sizes = defaultSize
	}
	return &TypeChecker{
		goVersion: goVersion,
		ctx:       ctx,
		sizes:     sizes,
	}
}

func (tc *TypeChecker) ModName() string { return `TypeChecker` }

type importer struct {
	pkg *project.Package
}

func (i *importer) Import(path string) (*types.Package, error) {
	if imp, ok := i.pkg.Imports[path]; ok {
		if imp.Types != nil {
			return imp.Types, nil
		}
		return nil, fmt.Errorf(`the types from package %q was nil`, path)
	}
	return nil, fmt.Errorf(`failed to find import for package %q`, path)
}

func (tc *TypeChecker) PackageDone(pkg *project.Package, errGroup *faults.Group) (con bool, err error) {
	ctx := tc.ctx
	if ctx == nil {
		ctx = types.NewContext()
	}

	cfg := &types.Config{
		Context:   tc.ctx,
		GoVersion: tc.goVersion,
		Error:     func(err error) { errGroup.Add(err) },
		Importer:  &importer{pkg: pkg},
		Sizes:     tc.sizes,
	}

	tPkg := &types.Package{}
	info := &types.Info{
		Types:        map[ast.Expr]types.TypeAndValue{},
		Instances:    map[*ast.Ident]types.Instance{},
		Defs:         map[*ast.Ident]types.Object{},
		Uses:         map[*ast.Ident]types.Object{},
		Implicits:    map[ast.Node]types.Object{},
		Selections:   map[*ast.SelectorExpr]*types.Selection{},
		Scopes:       map[ast.Node]*types.Scope{},
		InitOrder:    []*types.Initializer{},
		FileVersions: map[*ast.File]string{},
	}

	err = types.NewChecker(cfg, pkg.Fset, tPkg, info).Files(pkg.Syntax)
	if err != nil {
		return false, err
	}

	pkg.Types = tPkg
	pkg.TypesInfo = info
	pkg.TypesSizes = tc.sizes
	return true, nil
}
