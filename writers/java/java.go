package java

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/writers"
)

// New creates a new source writer for Java
func New() writers.Writer { return &writerImp{} }

type writerImp struct{}

func (w *writerImp) Name() string { return `java` }

func (w *writerImp) Write(cfg *writers.Config, proj constructs.IProject) error {
	// TODO: Implement
	return terror.New(`not implemented`)
}
