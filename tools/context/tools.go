package context

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Snow-Gremlin/Gozer/tools"
)

type toolSetImp struct {
	tools      []tools.Tool
	toolByName map[string]tools.Tool
}

func newToolSet() tools.ToolSet {
	return &toolSetImp{
		tools:      []tools.Tool{},
		toolByName: map[string]tools.Tool{},
	}
}

func (imp *toolSetImp) Count() int {
	return len(imp.tools)
}

func (imp *toolSetImp) At(index int) tools.Tool {
	return imp.tools[index]
}

func (imp *toolSetImp) setName(name string, t tools.Tool) {
	if t2, exists := imp.toolByName[name]; exists {
		panic(fmt.Errorf(`the name or alias, %s, for %T already exists for %T`, name, t, t2))
	}
	imp.toolByName[name] = t
}

func (imp *toolSetImp) Add(t tools.Tool) {
	imp.setName(t.Name(), t)
	for _, alias := range t.Aliases() {
		imp.setName(alias, t)
	}
	imp.tools = append(imp.tools, t)
}

func (imp *toolSetImp) Get(name string) tools.Tool {
	return imp.toolByName[name]
}

func (imp *toolSetImp) Names() []string {
	parts := make([]string, len(imp.tools))
	for i, tool := range imp.tools {
		names := append([]string{tool.Name()}, tool.Aliases()...)
		parts[i] = strings.Join(names, `, `)
	}
	sort.Strings(parts)
	return parts
}
