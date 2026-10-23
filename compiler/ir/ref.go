package ir

import "go/types"

type Ref interface {
	Node

	RefNode()

	Type() types.Type

	Object() types.Object
}
