package ir

type Directive interface {
	Node
	DirectiveNode()
}
