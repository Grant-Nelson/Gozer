package ir

import "go/types"

type Def interface {
	Stmt

	DefNode()

	Type() types.Type

	Object() types.Object
}
