package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type extraTypeImp struct {
	id      uint64
	name    string
	extends collections.List[TypeModel]
}

func (proj *projectImp) AddExtraType(name string) ExtraTypeModel {
	imp := &extraTypeImp{
		id:      proj.nextId(),
		name:    name,
		extends: list.New[TypeModel](),
	}
	proj.ExtraTypes().Append(imp)
	return imp
}

func (imp *extraTypeImp) _typeModel()  {}
func (imp *extraTypeImp) Id() uint64   { return imp.id }
func (imp *extraTypeImp) Name() string { return imp.name }

func (imp *extraTypeImp) Extends() collections.List[TypeModel] { return imp.extends }

func (imp *extraTypeImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addVal(data, `id`, imp.Id())
	addVal(data, `name`, imp.Name())
	addIds(data, `extends`, imp.Extends())
	return json.Marshal(data)
}
