package reader

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/readers"
	"github.com/Snow-Gremlin/Gozer/readers/golang"
	"github.com/Snow-Gremlin/Gozer/readers/json"
)

// New creates a new source reader that will select the correct reader by name.
func New() readers.Reader {
	return &readerImp{
		rs: map[string]readers.Reader{
			`golang`: golang.New(),
			`json`:   json.New(),
		},
	}
}

type readerImp struct {
	rs map[string]readers.Reader
}

func (r *readerImp) Read(cfg *readers.Config) (*constructs.CProject, error) {
	if r2, ok := r.rs[cfg.ReaderName]; ok {
		return r2.Read(cfg)
	}
	return nil, terror.New(`failed to find a reader by the requested name`).
		With(`reader name`, cfg.ReaderName)
}
