package reader

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
)

type Config struct {
	Name            string
	MainPackagePath string
	Context         *build.Context
	AugmentFiles    func(args *AugmentFilesArgs) []*ast.File
	ConvertPackage  func(args *ConvertPackageArgs)
}

type AugmentFilesArgs struct {
	Path    string
	FileSet *token.FileSet
	Files   []*ast.File
}

type ConvertPackageArgs struct {
	FileSet *token.FileSet
	Package *types.Package
	Info    *types.Info
	Files   []*ast.File
}

func Read(config *Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = terror.RecoveredPanic(r)
		}
	}()

	proj := newReader(config)
	proj.parsePackage(config.MainPackagePath)
	return nil
}

type reader struct {
	config   *Config
	packages map[string]*types.Package
}

func newReader(config *Config) *reader {
	if config.Context == nil {
		config.Context = &build.Default
	}
	return &reader{
		config:   config,
		packages: make(map[string]*types.Package),
	}
}

func (r *reader) Import(path string) (_ *types.Package, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = terror.RecoveredPanic(r)
		}
	}()

	return r.parsePackage(path), nil
}

func (r *reader) parsePackage(path string) *types.Package {
	if p, exists := r.packages[path]; exists {
		return p
	}

	paths := r.findPackageFiles(path)
	fileSet, files := parseFiles(paths)
	files = r.tryAugmentFiles(path, fileSet, files)
	pkg, info := r.getInfo(path, fileSet, files)
	r.tryConvertPackage(fileSet, pkg, info, files)
	return pkg
}

func (r *reader) tryAugmentFiles(path string, fileSet *token.FileSet, files []*ast.File) []*ast.File {
	if r.config.AugmentFiles != nil {
		files = r.config.AugmentFiles(&AugmentFilesArgs{
			Path:    path,
			FileSet: fileSet,
			Files:   files,
		})
	}
	return files
}

func (r *reader) tryConvertPackage(fileSet *token.FileSet, pkg *types.Package, info *types.Info, files []*ast.File) {
	if r.config.ConvertPackage != nil {
		r.config.ConvertPackage(&ConvertPackageArgs{
			FileSet: fileSet,
			Package: pkg,
			Info:    info,
			Files:   files,
		})
	}
}

func (r *reader) findPackageFiles(path string) []string {
	buildPackage, err := r.config.Context.ImportDir(path, build.FindOnly)
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

func (r *reader) getInfo(path string, fileSet *token.FileSet, files []*ast.File) (*types.Package, *types.Info) {
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
		Importer: r,
	}

	pkg, err := config.Check(path, fileSet, files, info)
	if err != nil {
		panic(terror.New(`type checker error`).
			With(`path`, path).
			WithError(err))
	}

	r.packages[path] = pkg
	return pkg, info
}
