package ir

import "go/types"

type Decl interface {
	Node

	DeclNode()

	Type() types.Type

	Object() types.Object
}
