package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type IGeneric interface {
	TypeParameters() collections.List[ITypeParameter]
	IsGeneric() bool
}

type genericImp struct {
	typeParams collections.List[ITypeParameter]
}

func (imp *genericImp) TypeParameters() collections.List[ITypeParameter] {
	return imp.typeParams
}

func (imp *genericImp) IsGeneric() bool {
	return !imp.typeParams.Empty()
}

func newGeneric() genericImp {
	return genericImp{
		typeParams: list.New[ITypeParameter](),
	}
}
