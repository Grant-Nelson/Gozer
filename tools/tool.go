package tools

type Tool interface {
	Name() string
	Summary() string
	Aliases() []string
	Description() string
	Run(ctx *Context) (int, error)
}
