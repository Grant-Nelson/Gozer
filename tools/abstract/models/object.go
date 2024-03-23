package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type objectImp struct {
	id         uint64
	name       string
	exported   bool
	extends    collections.List[TypeModel]
	typeParams collections.List[TypeModel]
	fields     collections.List[TypeModel]
}

func (pkg *packageImp) AddObject(name string, exported bool) ObjectModel {
	imp := &objectImp{
		id:         pkg.Project().nextId(),
		name:       name,
		exported:   exported,
		extends:    list.New[TypeModel](),
		typeParams: list.New[TypeModel](),
		fields:     list.New[TypeModel](),
	}
	pkg.Project().AllObjects().Append(imp)
	pkg.Objects().Append(imp)
	return imp
}

func (imp *objectImp) _typeModel()    {}
func (imp *objectImp) Id() uint64     { return imp.id }
func (imp *objectImp) Name() string   { return imp.name }
func (imp *objectImp) Exported() bool { return imp.exported }

func (imp *objectImp) Extends() collections.List[TypeModel]    { return imp.extends }
func (imp *objectImp) TypeParams() collections.List[TypeModel] { return imp.typeParams }
func (imp *objectImp) FieldTypes() collections.List[TypeModel] { return imp.fields }

func (imp *objectImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addVal(data, `id`, imp.Id())
	addVal(data, `name`, imp.Name())
	addVal(data, `exported`, imp.Exported())
	addIds(data, `extends`, imp.Extends())
	addIds(data, `typeParams`, imp.TypeParams())
	addIds(data, `fieldTypes`, imp.FieldTypes())
	return json.Marshal(data)
}
