package typeChecker

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
)

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

type typeCheckerMod struct {
	ctx       *types.Context
	goVersion string
	sizes     types.Sizes
	pkg       *project.Package
	errGroup  *faults.Group
}

var (
	_ mods.ModFactory = (*TypeChecker)(nil)
	_ mods.Modifier   = (*typeCheckerMod)(nil)
	_ types.Importer  = (*typeCheckerMod)(nil)
)

// The default size to use when no size is provided.
//
// This default size is based on what sizes of int (i.e. 32 bit or 64 bit)
// would fit into a JS number. A JS number is IEEE-754 64 bit double float
// with a mantissa of 53 bits, therefore a 32 bit int is used.
var defaultSize types.Sizes = &types.StdSizes{WordSize: 4, MaxAlign: 8}

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

func (tc *TypeChecker) StartPackage(pkg *project.Package, errGroup *faults.Group) (bool, mods.Modifier, error) {
	mod := &typeCheckerMod{
		ctx:       tc.ctx,
		goVersion: tc.goVersion,
		sizes:     tc.sizes,
		pkg:       pkg,
		errGroup:  errGroup,
	}
	return true, mod, nil
}

func (tc *typeCheckerMod) Import(path string) (*types.Package, error) {
	if imp, ok := tc.pkg.Ast.Imports[path]; ok {
		if imp.Types != nil {
			return imp.Types, nil
		}
		return nil, fmt.Errorf(`the types from package %q was nil`, path) // TODO: Change to faults
	}
	return nil, fmt.Errorf(`failed to find import for package %q`, path) // TODO: Change to faults
}

func (tc *typeCheckerMod) newInfo() *types.Info {
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

func (tc *typeCheckerMod) addError(err error) {
	tc.errGroup.Add(err)
}

func (tc *typeCheckerMod) PackageDone() (bool, error) {
	ctx := tc.ctx
	if ctx == nil {
		ctx = types.NewContext()
	}

	cfg := &types.Config{
		Context:   tc.ctx,
		GoVersion: tc.goVersion,
		Error:     tc.addError,
		Importer:  tc,
		Sizes:     tc.sizes,
	}

	pkg := tc.pkg.Ast
	typPkg := &types.Package{}
	info := tc.newInfo()

	err := types.NewChecker(cfg, pkg.Fset, typPkg, info).
		Files(pkg.Syntax)
	if err != nil {
		return false, tc.errGroup.Add(err)
	}

	pkg.Types = typPkg
	pkg.TypesInfo = info
	pkg.TypesSizes = tc.sizes
	return true, nil
}
