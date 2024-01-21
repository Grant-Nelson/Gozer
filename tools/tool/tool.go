package tool

import (
	"strings"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/Snow-Gremlin/Gozer/tools"
	"github.com/Snow-Gremlin/Gozer/tools/convert"
	"github.com/Snow-Gremlin/Gozer/tools/help"
)

// New creates a new main tool.
func New() tools.Tool {
	return &toolImp{
		ts: map[string]tools.Tool{
			`help`:    help.New(),
			`convert`: convert.New(),
		},
	}
}

type toolImp struct {
	ts map[string]tools.Tool
}

func (t *toolImp) Run(args []string) error {
	if len(args) <= 0 {
		return terror.New(`arguments require a tool name`).
			With(`tool names`, strings.Join(utils.SortedKeys(t.ts), `, `))
	}
	if t2, ok := t.ts[args[0]]; ok {
		return t2.Run(args[1:])
	}
	return terror.New(`unknown tool name`).
		With(`given name`, args[0]).
		With(`expected names`, strings.Join(utils.SortedKeys(t.ts), `, `))
}
