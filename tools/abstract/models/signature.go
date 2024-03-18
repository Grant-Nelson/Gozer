package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type signatureImp struct {
	index      uint64
	typeIndex  uint64
	name       string
	typeParams collections.List[TypeModel]
	params     collections.List[TypeModel]
	returns    collections.List[TypeModel]
}

func NewSignature(name string) SignatureModel {
	return &signatureImp{
		index:      0,
		typeIndex:  0,
		name:       name,
		typeParams: list.New[TypeModel](),
		params:     list.New[TypeModel](),
		returns:    list.New[TypeModel](),
	}
}

func (imp *signatureImp) Index() uint64             { return imp.index }
func (imp *signatureImp) setIndex(index uint64)     { imp.index = index }
func (imp *signatureImp) TypeIndex() uint64         { return imp.typeIndex }
func (imp *signatureImp) setTypeIndex(index uint64) { imp.typeIndex = index }
func (imp *signatureImp) Name() string              { return imp.name }

func (imp *signatureImp) TypeParams() collections.List[TypeModel] { return imp.typeParams }
func (imp *signatureImp) Params() collections.List[TypeModel]     { return imp.params }
func (imp *signatureImp) Returns() collections.List[TypeModel]    { return imp.returns }

func (imp *signatureImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addName(data, imp.name)
	addIndex(data, `index`, imp.index)
	addIndex(data, `typeIndex`, imp.typeIndex)
	addTypeIndices(data, `typeParams`, imp.typeParams)
	addTypeIndices(data, `params`, imp.params)
	addTypeIndices(data, `returns`, imp.returns)
	return json.Marshal(data)
}
