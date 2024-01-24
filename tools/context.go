package tools

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Context struct {
	tools      []Tool
	toolByName map[string]Tool
	args       []string
}

func NewContext() *Context {
	return &Context{
		tools:      []Tool{},
		toolByName: map[string]Tool{},
		args:       os.Args,
	}
}
func (c *Context) Args() []string {
	return c.args
}

func (c *Context) Tools() []Tool {
	return c.tools
}

func (c *Context) setName(name string, t Tool) {
	if t2, exists := c.toolByName[name]; exists {
		panic(fmt.Errorf(`the name or alias, %s, for %T already exists for %T`, name, t, t2))
	}
	c.toolByName[name] = t
}

func (c *Context) AddTool(t Tool) {
	c.tools = append(c.tools, t)
	c.setName(t.Name(), t)
	for _, alias := range t.Aliases() {
		c.setName(alias, t)
	}
}

func (c *Context) GetTool(name string) Tool {
	return c.toolByName[name]
}

func (c *Context) ToolNames() []string {
	parts := make([]string, len(c.tools))
	for i, tool := range c.tools {
		names := append([]string{tool.Name()}, tool.Aliases()...)
		parts[i] = strings.Join(names, `, `)
	}
	sort.Strings(parts)
	return parts
}
