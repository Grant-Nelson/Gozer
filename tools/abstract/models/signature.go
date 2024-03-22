package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type signatureImp struct {
	id         uint64
	typeParams collections.List[TypeModel]
	params     collections.List[TypeModel]
	returns    collections.List[TypeModel]
}

func (proj *projectImp) AddSignature() SignatureModel {
	imp := &signatureImp{
		id:         proj.nextId(),
		typeParams: list.New[TypeModel](),
		params:     list.New[TypeModel](),
		returns:    list.New[TypeModel](),
	}
	proj.Signatures().Append(imp)
	return imp
}

func (imp *signatureImp) _typeModel() {}
func (imp *signatureImp) Id() uint64  { return imp.id }

func (imp *signatureImp) TypeParams() collections.List[TypeModel] { return imp.typeParams }
func (imp *signatureImp) Params() collections.List[TypeModel]     { return imp.params }
func (imp *signatureImp) Returns() collections.List[TypeModel]    { return imp.returns }

func (imp *signatureImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addVal(data, `id`, imp.Id())
	addIds(data, `typeParams`, imp.TypeParams())
	addIds(data, `params`, imp.Params())
	addIds(data, `returns`, imp.Returns())
	return json.Marshal(data)
}
