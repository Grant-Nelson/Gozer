package convert

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/readers/reader"
	"github.com/Snow-Gremlin/Gozer/tools"
	"github.com/Snow-Gremlin/Gozer/writers/writer"
)

// New creates a new convert tool.
func New() tools.Tool { return &toolImp{} }

type toolImp struct{}

func (t *toolImp) Run(args []string) error {
	_ = reader.New()
	_ = writer.New()

	// TODO: Implement
	return terror.New(`convert not implemented yet`)
}
