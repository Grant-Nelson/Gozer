package tools

type Context interface {
	Args() []string
	Tools() ToolSet
}
