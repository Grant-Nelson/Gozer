package json

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/writers"
)

// New creates a new source writer for JSON
func New() writers.Writer { return &writerImp{} }

type writerImp struct{}

func (w *writerImp) Write(cfg *writers.Config, proj *constructs.CProject) error {
	// TODO: Implement
	return terror.New(`not implemented`)
}
