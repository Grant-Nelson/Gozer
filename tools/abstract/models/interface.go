package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type interfaceImp struct {
	id         uint64
	name       string
	exported   bool
	typeParams collections.List[TypeModel]
	extends    collections.List[TypeModel]
	signatures collections.List[SignatureModel]
}

func (pkg *packageImp) AddInterface(name string, exported bool) InterfaceModel {
	imp := &interfaceImp{
		id:         pkg.Project().nextId(),
		name:       name,
		exported:   exported,
		typeParams: list.New[TypeModel](),
		extends:    list.New[TypeModel](),
		signatures: list.New[SignatureModel](),
	}
	pkg.Project().AllInterfaces().Append(imp)
	pkg.Interfaces().Append(imp)
	return imp
}

func (imp *interfaceImp) _typeModel()    {}
func (imp *interfaceImp) Id() uint64     { return imp.id }
func (imp *interfaceImp) Name() string   { return imp.name }
func (imp *interfaceImp) Exported() bool { return imp.exported }

func (imp *interfaceImp) TypeParams() collections.List[TypeModel]      { return imp.typeParams }
func (imp *interfaceImp) Extends() collections.List[TypeModel]         { return imp.extends }
func (imp *interfaceImp) Signatures() collections.List[SignatureModel] { return imp.signatures }

func (imp *interfaceImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addVal(data, `id`, imp.Id())
	addVal(data, `name`, imp.Name())
	addVal(data, `exported`, imp.Exported())
	addIds(data, `typeParams`, imp.TypeParams())
	addIds(data, `extends`, imp.Extends())
	addIds(data, `signatures`, imp.Signatures())
	return json.Marshal(data)
}
