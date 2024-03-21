package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type typeCollectionImp struct {
	interfaces collections.List[InterfaceModel]
	objects    collections.List[ObjectModel]
	signatures collections.List[SignatureModel]
	extraTypes collections.List[ExtraTypeModel]
}

func NewTypeCollection() TypeCollectionModel {
	return &typeCollectionImp{
		interfaces: list.New[InterfaceModel](),
		objects:    list.New[ObjectModel](),
		signatures: list.New[SignatureModel](),
		extraTypes: list.New[ExtraTypeModel](),
	}
}
func (imp *typeCollectionImp) Interfaces() collections.List[InterfaceModel] { return imp.interfaces }
func (imp *typeCollectionImp) Objects() collections.List[ObjectModel]       { return imp.objects }
func (imp *typeCollectionImp) Signatures() collections.List[SignatureModel] { return imp.signatures }
func (imp *typeCollectionImp) ExtraTypes() collections.List[ExtraTypeModel] { return imp.extraTypes }

func (imp *typeCollectionImp) MarshalJSON() ([]byte, error) {
	setIndices(imp.interfaces)
	setIndices(imp.objects)
	setIndices(imp.signatures)
	setIndices(imp.extraTypes)

	var typeIndex uint64
	setTypeIndices(&typeIndex, imp.interfaces)
	setTypeIndices(&typeIndex, imp.objects)
	setTypeIndices(&typeIndex, imp.signatures)
	setTypeIndices(&typeIndex, imp.extraTypes)

	data := map[string]any{}
	addData(data, `interfaces`, imp.Interfaces())
	addData(data, `objects`, imp.Objects())
	addData(data, `signatures`, imp.Signatures())
	addData(data, `extra`, imp.ExtraTypes())
	return json.Marshal(data)
}
