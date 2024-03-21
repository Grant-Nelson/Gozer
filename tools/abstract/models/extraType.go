package models

import "encoding/json"

type extraTypeImp struct {
	index     uint64
	typeIndex uint64
	name      string
}

func NewExtraType(name string) ExtraTypeModel {
	return &extraTypeImp{
		index:     0,
		typeIndex: 0,
		name:      name,
	}
}

func (imp *extraTypeImp) Index() uint64             { return imp.index }
func (imp *extraTypeImp) setIndex(index uint64)     { imp.index = index }
func (imp *extraTypeImp) TypeIndex() uint64         { return imp.typeIndex }
func (imp *extraTypeImp) setTypeIndex(index uint64) { imp.typeIndex = index }
func (imp *extraTypeImp) Name() string              { return imp.name }

func (imp *extraTypeImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addString(data, `name`, imp.Name())
	return json.Marshal(data)
}
