package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type extraTypeImp struct {
	id         uint64
	name       string
	exported   bool
	extends    collections.List[TypeModel]
	typeParams collections.List[TypeModel]
}

func (proj *projectImp) AddExtraType(name string, exported bool) ExtraTypeModel {
	imp := &extraTypeImp{
		id:         proj.nextId(),
		name:       name,
		exported:   exported,
		extends:    list.New[TypeModel](),
		typeParams: list.New[TypeModel](),
	}
	proj.ExtraTypes().Append(imp)
	return imp
}

func (imp *extraTypeImp) _typeModel()    {}
func (imp *extraTypeImp) _extraType()    {}
func (imp *extraTypeImp) Id() uint64     { return imp.id }
func (imp *extraTypeImp) Name() string   { return imp.name }
func (imp *extraTypeImp) Exported() bool { return imp.exported }

func (imp *extraTypeImp) Extends() collections.List[TypeModel]    { return imp.extends }
func (imp *extraTypeImp) TypeParams() collections.List[TypeModel] { return imp.typeParams }

func (imp *extraTypeImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addVal(data, `id`, imp.Id())
	addVal(data, `name`, imp.Name())
	addVal(data, `exported`, imp.Exported())
	addIds(data, `extends`, imp.Extends())
	addIds(data, `typeParams`, imp.TypeParams())
	return json.Marshal(data)
}
