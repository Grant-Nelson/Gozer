package golang

import (
	"go/build"
	"go/token"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/constructs/cPackage"
	"github.com/Snow-Gremlin/Gozer/constructs/cProject"
	"github.com/Snow-Gremlin/Gozer/readers"
)

// New creates a new source reader for Golang
func New() readers.Reader { return &readerImp{} }

type readerImp struct{}

func (r *readerImp) Name() string { return `golang` }

func (r *readerImp) Read(cfg *readers.Config) (proj constructs.CProject, err error) {
	defer func() {
		if r := recover(); r != nil {
			proj = nil
			err = terror.RecoveredPanic(r)
		}
	}()

	context := build.Default
	fSet := token.NewFileSet()
	proj = cProject.New()
	proj.SetName(cfg.MainPackageDir)
	readPackage(cfg.MainPackageDir, context, fSet, proj)
	return proj, nil
}

func readPackage(path string, context build.Context, fSet *token.FileSet, proj constructs.CProject) constructs.CPackage {
	if p, exists := proj.Packages().TryGetByPath(path); exists {
		return p
	}

	p := cPackage.New()
	p.SetPath(path)
	proj.Packages().Add(p)

	pkg, err := context.ImportDir(path, build.FindOnly)
	if err != nil {
		panic(err)
	}

	p.SetName(pkg.Name)
	for _, inPath := range pkg.Imports {
		p.Imports().Add(readPackage(inPath, context, fSet, proj))
	}

	newConverter(p, pkg, fSet, proj).
		convertPackage()
	return p
}
