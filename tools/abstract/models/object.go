package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type objectImp struct {
	index      uint64
	typeIndex  uint64
	name       string
	typeParams collections.List[TypeModel]
	fields     collections.List[FieldModel]
}

func NewObject(name string) ObjectModel {
	return &objectImp{
		index:      0,
		typeIndex:  0,
		name:       name,
		typeParams: list.New[TypeModel](),
		fields:     list.New[FieldModel](),
	}
}

func (imp *objectImp) Index() uint64             { return imp.index }
func (imp *objectImp) setIndex(index uint64)     { imp.index = index }
func (imp *objectImp) TypeIndex() uint64         { return imp.typeIndex }
func (imp *objectImp) setTypeIndex(index uint64) { imp.typeIndex = index }
func (imp *objectImp) Name() string              { return imp.name }

func (imp *objectImp) TypeParams() collections.List[TypeModel] { return imp.typeParams }
func (imp *objectImp) Fields() collections.List[FieldModel]    { return imp.fields }

func (imp *objectImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addName(data, imp.name)
	addIndex(data, `index`, imp.index)
	addIndex(data, `typeIndex`, imp.typeIndex)
	addTypeIndices(data, `typeParams`, imp.typeParams)
	addIndices(data, `fields`, imp.fields)
	return json.Marshal(data)
}
