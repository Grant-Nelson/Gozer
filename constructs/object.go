package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type IObject interface {
	IType
	INamed
	IGeneric
	Data() IType
	SetData(data IType)
	Methods() collections.List[IMethod]
	_objectType()
}

type objectImp struct {
	namedImp
	genericImp
	data    IType
	methods collections.List[IMethod]
}

func (imp *objectImp) _type()       {}
func (imp *objectImp) _objectType() {}

func (imp *objectImp) Data() IType {
	return imp.data
}

func (imp *objectImp) SetData(data IType) {
	imp.data = data
}

func (imp *objectImp) Methods() collections.List[IMethod] {
	return imp.methods
}

func NewObject(name string, data IType, methods ...IMethod) IObject {
	return &objectImp{
		namedImp:   newName(name),
		genericImp: newGeneric(),
		data:       data,
		methods:    list.With(methods...),
	}
}
