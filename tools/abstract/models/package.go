package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type packageImp struct {
	index      uint64
	name       string
	parent     PackageModel
	interfaces collections.List[InterfaceModel]
	objects    collections.List[ObjectModel]
	methods    collections.List[MethodModel]
	statics    collections.List[TypeModel]
}

func NewPackage(parent PackageModel, name string) PackageModel {
	return &packageImp{
		index:      0,
		name:       name,
		parent:     parent,
		interfaces: list.New[InterfaceModel](),
		objects:    list.New[ObjectModel](),
		methods:    list.New[MethodModel](),
		statics:    list.New[TypeModel](),
	}
}

func (imp *packageImp) Index() uint64         { return imp.index }
func (imp *packageImp) setIndex(index uint64) { imp.index = index }
func (imp *packageImp) Name() string          { return imp.name }
func (imp *packageImp) Parent() PackageModel  { return imp.parent }

func (imp *packageImp) Interfaces() collections.List[InterfaceModel] { return imp.interfaces }
func (imp *packageImp) Objects() collections.List[ObjectModel]       { return imp.objects }
func (imp *packageImp) Methods() collections.List[MethodModel]       { return imp.methods }
func (imp *packageImp) Statics() collections.List[TypeModel]         { return imp.statics }

func (imp *packageImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addName(data, imp.name)
	addIndex(data, `index`, imp.index)
	addIndexed(data, `parent`, imp.parent)
	addIndices(data, `interfaces`, imp.interfaces)
	addIndices(data, `objects`, imp.objects)
	addIndices(data, `methods`, imp.methods)
	addTypeIndices(data, `statics`, imp.statics)
	return json.Marshal(data)
}
