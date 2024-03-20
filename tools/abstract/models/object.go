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
	implements collections.List[InterfaceModel]
	extends    collections.List[TypeModel]
	typeParams collections.List[TypeModel]
	fields     collections.List[TypeModel]
}

func NewObject(name string) ObjectModel {
	return &objectImp{
		index:      0,
		typeIndex:  0,
		name:       name,
		implements: list.New[InterfaceModel](),
		extends:    list.New[TypeModel](),
		typeParams: list.New[TypeModel](),
		fields:     list.New[TypeModel](),
	}
}

func (imp *objectImp) Index() uint64             { return imp.index }
func (imp *objectImp) setIndex(index uint64)     { imp.index = index }
func (imp *objectImp) TypeIndex() uint64         { return imp.typeIndex }
func (imp *objectImp) setTypeIndex(index uint64) { imp.typeIndex = index }
func (imp *objectImp) Name() string              { return imp.name }

func (imp *objectImp) Implements() collections.List[InterfaceModel] { return imp.implements }
func (imp *objectImp) Extends() collections.List[TypeModel]         { return imp.extends }
func (imp *objectImp) TypeParams() collections.List[TypeModel]      { return imp.typeParams }
func (imp *objectImp) Fields() collections.List[TypeModel]          { return imp.fields }

func (imp *objectImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addString(data, `name`, imp.Name())
	addIndex(data, `index`, imp.Index())
	addIndex(data, `typeIndex`, imp.TypeIndex())
	addIndices(data, `implements`, imp.Implements())
	addTypeIndices(data, `extends`, imp.Extends())
	addTypeIndices(data, `typeParams`, imp.TypeParams())
	addTypeIndices(data, `fields`, imp.Fields())
	return json.Marshal(data)
}
