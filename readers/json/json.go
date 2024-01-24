package json

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/readers"
)

// New creates a new source reader for JSON
func New() readers.Reader { return &readerImp{} }

type readerImp struct{}

func (r *readerImp) Name() string { return `json` }

func (r *readerImp) Read(cfg *readers.Config) (*constructs.CProject, error) {
	// TODO: Implement
	return nil, terror.New(`not implemented`)
}
