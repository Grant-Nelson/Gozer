package convert

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/tools"
)

// New creates a new convert tool.
func New() tools.Tool { return &toolImp{} }

type toolImp struct{}

func (t *toolImp) Name() string { return `convert` }

func (t *toolImp) Aliases() []string { return []string{`conv`} }

func (t *toolImp) Summary() string {
	return `` // TODO: Implement
}

func (t *toolImp) Description() string {
	return `` // TODO: Implement
}

func (t *toolImp) Run(ctx tools.Context) (int, error) {

	// TODO: Implement
	return 0, terror.New(`convert not implemented yet`)
}
