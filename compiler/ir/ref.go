package ir

type Ref interface {
	Node

	Decl() Decl
}
