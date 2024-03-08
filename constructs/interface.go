package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type IInterface interface {
	IType
	INamed
	IGeneric
	Methods() collections.List[ISignature]
	_interfaceType()
}

type interfaceImp struct {
	namedImp
	genericImp
	methods collections.List[ISignature]
}

func (imp *interfaceImp) _type()          {}
func (imp *interfaceImp) _interfaceType() {}

func (imp *interfaceImp) Methods() collections.List[ISignature] {
	return imp.methods
}

func NewInterface(name string, methods ...ISignature) IInterface {
	return &interfaceImp{
		namedImp:   newName(name),
		genericImp: newGeneric(),
		methods:    list.With(methods...),
	}
}
