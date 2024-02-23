package golang

import (
	"fmt"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/readers"
	"github.com/Snow-Gremlin/Gozer/readers/golang/packageSet"
)

// New creates a new source reader for Golang
func New() readers.Reader { return &readerImp{} }

type readerImp struct{}

func (r *readerImp) Name() string { return `golang` }

func (r *readerImp) Read(cfg *readers.Config) (*constructs.CProject, error) {
	ps, err := packageSet.New(cfg.MainPackageDir)
	if err != nil {
		return nil, err
	}

	fmt.Println(ps)

	// TODO: Implement
	return nil, terror.New(`not implemented`)
}
