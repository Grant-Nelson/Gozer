package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type interfaceImp struct {
	id         uint64
	name       string
	typeParams collections.List[TypeModel]
	signatures collections.List[SignatureModel]
}

func (pkg *packageImp) AddInterface(name string) InterfaceModel {
	imp := &interfaceImp{
		id:         pkg.Project().nextId(),
		name:       name,
		typeParams: list.New[TypeModel](),
		signatures: list.New[SignatureModel](),
	}
	pkg.Project().Interfaces().Append(imp)
	pkg.Interfaces().Append(imp)
	return imp
}

func (imp *interfaceImp) _typeModel()  {}
func (imp *interfaceImp) Id() uint64   { return imp.id }
func (imp *interfaceImp) Name() string { return imp.name }

func (imp *interfaceImp) TypeParams() collections.List[TypeModel]      { return imp.typeParams }
func (imp *interfaceImp) Signatures() collections.List[SignatureModel] { return imp.signatures }

func (imp *interfaceImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addVal(data, `id`, imp.Id())
	addVal(data, `name`, imp.Name())
	addIds(data, `typeParams`, imp.TypeParams())
	addIds(data, `signatures`, imp.Signatures())
	return json.Marshal(data)
}
