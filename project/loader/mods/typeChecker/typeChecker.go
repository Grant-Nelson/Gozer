package typeChecker

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
)

// The default size to use when no size is provided.
//
// This default size is based on what sizes of int (i.e. 32 bit or 64 bit)
// would fit into a JS number. A JS number is IEEE-754 64 bit double float
// with a mantissa of 53 bits, therefore a 32 bit int is used.
var defaultSize types.Sizes = &types.StdSizes{WordSize: 4, MaxAlign: 8}

type Config struct {

	// Context is the context used for resolving global identifiers. If nil, the
	// type checker will initialize this field with a newly created context.
	Context *types.Context

	// GoVersion describes the accepted Go language version. The string must
	// start with a prefix of the form "go%d.%d" (e.g. "go1.20", "go1.21rc1", or
	// "go1.21.0") or it must be empty; an empty string disables Go language
	// version checks. If the format is invalid, invoking the type checker will
	// result in an error.
	GoVersion string

	// Sizes is the byte size and byte alignment to use when generations types.
	// This will affect int and uint sizes and how the fields align in struct.
	// If nil, then a default 32 bit sizes and 8 byte max alignment.
	Sizes types.Sizes
}

type TypeChecker struct {
	ctx       *types.Context
	goVersion string
	sizes     types.Sizes
}

var (
	_ mods.Modifier       = (*TypeChecker)(nil)
	_ mods.PackageDoneExt = (*TypeChecker)(nil)
	_ types.Importer      = (*importer)(nil)
)

func New(cfg *Config) *TypeChecker {
	mod := &TypeChecker{
		sizes: defaultSize,
	}
	if cfg != nil {
		mod.ctx = cfg.Context
		mod.goVersion = cfg.GoVersion
		if cfg.Sizes != nil {
			mod.sizes = cfg.Sizes
		}
	}
	return mod
}

func (tc *TypeChecker) ModName() string { return `TypeChecker` }

type importer struct {
	pkg *project.Package
}

func (i *importer) Import(path string) (*types.Package, error) {
	if imp, ok := i.pkg.Ast.Imports[path]; ok {
		if imp.Types != nil {
			return imp.Types, nil
		}
		return nil, fmt.Errorf(`the types from package %q was nil`, path)
	}
	return nil, fmt.Errorf(`failed to find import for package %q`, path)
}

func (tc *TypeChecker) newInfo() *types.Info {
	return &types.Info{
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
	info := tc.newInfo()
	err = types.NewChecker(cfg, pkg.Ast.Fset, tPkg, info).Files(pkg.Ast.Syntax)
	if err != nil {
		return false, errGroup.Add(err)
	}

	pkg.Ast.Types = tPkg
	pkg.Ast.TypesInfo = info
	pkg.Ast.TypesSizes = tc.sizes
	return true, nil
}
