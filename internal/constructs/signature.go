package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type ISignature interface {
	IType
	INamed
	IGeneric
	Variadic() bool
	SetVariadic(variadic bool)
	Parameters() collections.List[IParameter]
	Returns() collections.List[IParameter]
	_signatureType()
}

type signatureImp struct {
	namedImp
	genericImp
	variadic   bool
	parameters collections.List[IParameter]
	returns    collections.List[IParameter]
}

func (imp *signatureImp) _type()          {}
func (imp *signatureImp) _signatureType() {}

func (imp *signatureImp) Variadic() bool {
	return imp.variadic
}

func (imp *signatureImp) SetVariadic(variadic bool) {
	imp.variadic = variadic
}

func (imp *signatureImp) Parameters() collections.List[IParameter] {
	return imp.parameters
}

func (imp *signatureImp) Returns() collections.List[IParameter] {
	return imp.returns
}

func NewSignature(name string) ISignature {
	return &signatureImp{
		namedImp:   newName(name),
		genericImp: newGeneric(),
		variadic:   false,
		parameters: list.New[IParameter](),
		returns:    list.New[IParameter](),
	}
}
