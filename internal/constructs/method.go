package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type IMethod interface {
	INamed
	Signature() ISignature
	SetSignature(signature ISignature)
	Body() collections.List[IStatement]
	_methodConstruct()
}

type methodImp struct {
	namedImp
	signature ISignature
	body      collections.List[IStatement]
}

func (imp *methodImp) _methodConstruct() {}

func (imp *methodImp) Signature() ISignature {
	return imp.signature
}

func (imp *methodImp) SetSignature(signature ISignature) {
	imp.signature = signature
}

func (imp *methodImp) Body() collections.List[IStatement] {
	return imp.body
}

func NewMethod(name string, signature ISignature) IMethod {
	return &methodImp{
		namedImp:  newName(name),
		signature: signature,
		body:      list.New[IStatement](),
	}
}
