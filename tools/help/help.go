package help

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/tools"
)

// New creates a new help tool.
func New() tools.Tool { return &toolImp{} }

type toolImp struct{}

func (t *toolImp) Run(args []string) error {
	// TODO: Implement
	return terror.New(`help not implemented yet`)
}
