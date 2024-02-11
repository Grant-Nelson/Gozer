package context

import (
	"os"

	"github.com/Snow-Gremlin/Gozer/tools"
)

type contextImp struct {
	tools tools.ToolSet
	args  []string
}

func New() tools.Context {
	return &contextImp{
		tools: newToolSet(),
		args:  os.Args,
	}
}
func (imp *contextImp) Args() []string {
	return imp.args
}

func (imp *contextImp) Tools() tools.ToolSet {
	return imp.tools
}
