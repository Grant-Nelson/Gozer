package golang

import (
	"fmt"
	"go/build"
	"go/token"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/readers"
	"github.com/Snow-Gremlin/Gozer/readers/golang/converter"
	"github.com/Snow-Gremlin/Gozer/readers/golang/packageSet"
)

// New creates a new source reader for Golang
func New() readers.Reader { return &readerImp{} }

type readerImp struct{}

func (r *readerImp) Name() string { return `golang` }

func (r *readerImp) Read(cfg *readers.Config) (*constructs.CProject, error) {
	pkgSet := packageSet.New(build.Default)
	if err := pkgSet.Add(cfg.MainPackageDir); err != nil {
		return nil, err
	}

	fSet := token.NewFileSet()
	proj := constructs.NewProject(cfg.MainPackageDir)

	p, err := converter.Convert(pkgSet.MainPackage(), pkgSet, fSet, proj)
	if err != nil {
		return nil, err
	}

	fmt.Println(p)

	// TODO: Implement
	return nil, terror.New(`not implemented`)
}
