package golang

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/readers"
)

// New creates a new source reader for Golang
func New() readers.Reader { return &readerImp{} }

type readerImp struct{}

func (r *readerImp) Name() string { return `golang` }

func (r *readerImp) Read(cfg *readers.Config) (cProj constructs.IProject, err error) {
	defer func() {
		if r := recover(); r != nil {
			cProj = nil
			err = terror.RecoveredPanic(r)
		}
	}()

	proj := &goProject{
		context: &build.Default,
		fileSet: token.NewFileSet(),
		config:  &types.Config{},
		info: &types.Info{
			Types:      map[ast.Expr]types.TypeAndValue{},
			Instances:  map[*ast.Ident]types.Instance{},
			Defs:       map[*ast.Ident]types.Object{},
			Uses:       map[*ast.Ident]types.Object{},
			Implicits:  map[ast.Node]types.Object{},
			Selections: map[*ast.SelectorExpr]*types.Selection{},
			Scopes:     map[ast.Node]*types.Scope{},
		},
		packages: make(map[string]*goPackage),
	}
	proj.config.Importer = proj
	proj.mainPkg = createPackage(proj, cfg.MainPackageDir)
	return convert(proj), nil
}

type (
	goProject struct {
		context  *build.Context
		fileSet  *token.FileSet
		config   *types.Config
		info     *types.Info
		mainPkg  *goPackage
		packages map[string]*goPackage
	}

	goPackage struct {
		proj         *goProject
		order        int
		buildPackage *build.Package
		typesPackage *types.Package
		imports      []*goPackage
		files        []*ast.File
	}
)

func createPackage(proj *goProject, path string) *goPackage {
	if p, exists := proj.packages[path]; exists {
		return p
	}

	pkg, err := proj.context.ImportDir(path, build.FindOnly)
	if err != nil {
		panic(err)
	}

	p := &goPackage{
		proj:         proj,
		buildPackage: pkg,
		imports:      make([]*goPackage, len(pkg.Imports)),
	}

	proj.packages[path] = p
	for i, inPath := range pkg.Imports {
		p.imports[i] = createPackage(proj, inPath)
	}

	p.updateOrder()
	p.populateFiles()
	// TODO: add optional Go file augmentation here
	p.populateInfo()
	return p
}

func (p *goPackage) updateOrder() {
	maxImportOrder := -1
	for _, in := range p.imports {
		maxImportOrder = max(maxImportOrder, in.order)
	}
	p.order = maxImportOrder + 1
}

func (p *goPackage) populateFiles() {
	bp := p.buildPackage
	p.files = make([]*ast.File, len(bp.GoFiles))
	for i, fileName := range bp.GoFiles {
		if !filepath.IsAbs(fileName) {
			fileName = filepath.Join(bp.Dir, fileName)
		}

		f, err := parser.ParseFile(p.proj.fileSet, fileName, nil, parser.ParseComments)
		if err != nil {
			panic(terror.New(`error parsing file`, err).
				With(`file name`, fileName))
		}

		p.files[i] = f
	}
}

func (p *goPackage) populateInfo() {
	tp, err := p.proj.config.Check(p.buildPackage.Dir, p.proj.fileSet, p.files, p.proj.info)
	if err != nil {
		panic(terror.New(`type checker error`).
			With(`path`, p.buildPackage.Dir).
			WithError(err))
	}
	p.typesPackage = tp
}

func (proj *goProject) Import(path string) (*types.Package, error) {
	if p, exists := proj.packages[path]; exists {
		return p.typesPackage, nil
	}
	return nil, terror.New(`parser error: checker requested package not preloaded`).
		With(`path`, path)
}
