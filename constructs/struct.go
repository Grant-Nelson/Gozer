package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type IStruct interface {
	IType
	IGeneric
	Fields() collections.List[IParameter]
	_structType()
}

type structImp struct {
	genericImp
	fields collections.List[IParameter]
}

func (imp *structImp) _type()       {}
func (imp *structImp) _structType() {}

func (imp *structImp) Fields() collections.List[IParameter] {
	return imp.fields
}

func (imp *structImp) String() string {
	return `struct`
}

func NewStruct(fields ...IParameter) IStruct {
	return &structImp{
		genericImp: newGeneric(),
		fields:     list.With(fields...),
	}
}
