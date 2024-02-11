package tools

type ToolSet interface {
	Count() int
	At(index int) Tool
	Add(t Tool)
	Get(name string) Tool
	Names() []string
}
