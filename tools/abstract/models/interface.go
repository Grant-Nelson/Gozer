package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type interfaceImp struct {
	index      uint64
	typeIndex  uint64
	name       string
	typeParams collections.List[TypeModel]
	signatures collections.List[SignatureModel]
}

func NewInterface(name string) InterfaceModel {
	return &interfaceImp{
		index:      0,
		typeIndex:  0,
		name:       name,
		typeParams: list.New[TypeModel](),
		signatures: list.New[SignatureModel](),
	}
}

func (imp *interfaceImp) Index() uint64             { return imp.index }
func (imp *interfaceImp) setIndex(index uint64)     { imp.index = index }
func (imp *interfaceImp) TypeIndex() uint64         { return imp.typeIndex }
func (imp *interfaceImp) setTypeIndex(index uint64) { imp.typeIndex = index }
func (imp *interfaceImp) Name() string              { return imp.name }

func (imp *interfaceImp) TypeParams() collections.List[TypeModel]      { return imp.typeParams }
func (imp *interfaceImp) Signatures() collections.List[SignatureModel] { return imp.signatures }

func (imp *interfaceImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addString(data, `name`, imp.Name())
	addIndex(data, `index`, imp.Index())
	addIndex(data, `typeIndex`, imp.TypeIndex())
	addTypeIndices(data, `typeParams`, imp.TypeParams())
	addIndices(data, `signatures`, imp.Signatures())
	return json.Marshal(data)
}
