package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type signatureImp struct {
	id         uint64
	name       string
	exported   bool
	extends    collections.List[TypeModel]
	typeParams collections.List[TypeModel]
	params     collections.List[TypeModel]
	returns    collections.List[TypeModel]
}

func (proj *projectImp) AddSignature(name string, exported bool) SignatureModel {
	imp := &signatureImp{
		id:         proj.nextId(),
		name:       name,
		exported:   exported,
		extends:    list.New[TypeModel](),
		typeParams: list.New[TypeModel](),
		params:     list.New[TypeModel](),
		returns:    list.New[TypeModel](),
	}
	proj.AllSignatures().Append(imp)
	return imp
}

func (imp *signatureImp) _typeModel()    {}
func (imp *signatureImp) Id() uint64     { return imp.id }
func (imp *signatureImp) Name() string   { return imp.name }
func (imp *signatureImp) Exported() bool { return imp.exported }

func (imp *signatureImp) Extends() collections.List[TypeModel]    { return imp.extends }
func (imp *signatureImp) TypeParams() collections.List[TypeModel] { return imp.typeParams }
func (imp *signatureImp) Params() collections.List[TypeModel]     { return imp.params }
func (imp *signatureImp) Returns() collections.List[TypeModel]    { return imp.returns }

func (imp *signatureImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addVal(data, `id`, imp.Id())
	addVal(data, `name`, imp.Name())
	addVal(data, `exported`, imp.Exported())
	addIds(data, `extends`, imp.Extends())
	addIds(data, `typeParams`, imp.TypeParams())
	addIds(data, `params`, imp.Params())
	addIds(data, `returns`, imp.Returns())
	return json.Marshal(data)
}
