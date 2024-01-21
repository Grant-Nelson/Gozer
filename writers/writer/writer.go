package writer

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/writers"
	"github.com/Snow-Gremlin/Gozer/writers/golang"
	"github.com/Snow-Gremlin/Gozer/writers/java"
	"github.com/Snow-Gremlin/Gozer/writers/json"
)

// New creates a new source writer that will select the correct writer by name.
func New() writers.Writer {
	return &writerImp{
		ws: map[string]writers.Writer{
			`golang`: golang.New(),
			`java`:   java.New(),
			`json`:   json.New(),
		},
	}
}

type writerImp struct {
	ws map[string]writers.Writer
}

func (w *writerImp) Write(cfg *writers.Config, proj *constructs.CProject) error {
	if w2, ok := w.ws[cfg.WriterName]; ok {
		return w2.Write(cfg, proj)
	}
	return terror.New(`failed to find a writer by the requested name`).
		With(`writer name`, cfg.WriterName)
}
