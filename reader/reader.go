package reader

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/constructs"
)

type Config struct {
	Name            string
	MainPackagePath string
	Context         *build.Context
	AugmentFiles    func(path string, fileSet *token.FileSet, files []*ast.File) []*ast.File
}

func Read(config *Config) (_ constructs.IProject, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = terror.RecoveredPanic(r)
		}
	}()

	proj := newProject(config)
	proj.parsePackage(config.MainPackagePath)
	return proj.project, nil
}

type goProject struct {
	*Config
	packages map[string]*types.Package
	project  constructs.IProject
}

func newProject(config *Config) *goProject {
	if config.Context == nil {
		config.Context = &build.Default
	}
	return &goProject{
		Config:   config,
		project:  constructs.NewProject(config.Name),
		packages: make(map[string]*types.Package),
	}
}

func (proj *goProject) Import(path string) (_ *types.Package, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = terror.RecoveredPanic(r)
		}
	}()

	return proj.parsePackage(path), nil
}

func (proj *goProject) parsePackage(path string) *types.Package {
	if p, exists := proj.packages[path]; exists {
		return p
	}

	paths := proj.findPackageFiles(path)
	fileSet, files := parseFiles(paths)
	tPack, info := proj.getInfo(path, fileSet, files)
	convert(proj.project, fileSet, tPack, info, files)
	return tPack
}

func (proj *goProject) findPackageFiles(path string) []string {
	buildPackage, err := proj.Context.ImportDir(path, build.FindOnly)
	if err != nil {
		panic(terror.New(`error reading import directory`, err).
			With(`path`, path))
	}

	paths := buildPackage.GoFiles
	return normalizePaths(buildPackage.Dir, paths)
}

func normalizePaths(dir string, paths []string) []string {
	for i, path := range paths {
		if !filepath.IsAbs(path) {
			paths[i] = filepath.Join(dir, path)
		}
	}
	return paths
}

func parseFiles(paths []string) (*token.FileSet, []*ast.File) {
	fileSet := token.NewFileSet()
	files := make([]*ast.File, len(paths))
	for i, path := range paths {
		f, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			panic(terror.New(`error parsing file`, err).
				With(`path`, path))
		}
		files[i] = f
	}
	return fileSet, files
}

func (proj *goProject) getInfo(path string, fileSet *token.FileSet, files []*ast.File) (*types.Package, *types.Info) {
	info := &types.Info{
		Types:      map[ast.Expr]types.TypeAndValue{},
		Instances:  map[*ast.Ident]types.Instance{},
		Defs:       map[*ast.Ident]types.Object{},
		Uses:       map[*ast.Ident]types.Object{},
		Implicits:  map[ast.Node]types.Object{},
		Selections: map[*ast.SelectorExpr]*types.Selection{},
		Scopes:     map[ast.Node]*types.Scope{},
	}

	config := &types.Config{
		Importer: proj,
	}

	pkg, err := config.Check(path, fileSet, files, info)
	if err != nil {
		panic(terror.New(`type checker error`).
			With(`path`, path).
			WithError(err))
	}

	proj.packages[path] = pkg
	return pkg, info
}
