package models

import "encoding/json"

type fieldImp struct {
	index uint64
	name  string
	fType TypeModel
}

func NewField(name string, fType TypeModel) FieldModel {
	return &fieldImp{
		index: 0,
		name:  name,
		fType: fType,
	}
}

func (imp *fieldImp) Index() uint64         { return imp.index }
func (imp *fieldImp) setIndex(index uint64) { imp.index = index }
func (imp *fieldImp) Name() string          { return imp.name }
func (imp *fieldImp) Type() TypeModel       { return imp.fType }

func (imp *fieldImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addName(data, imp.name)
	addIndex(data, `index`, imp.index)
	addTypeIndex(data, `type`, imp.fType)
	return json.Marshal(data)
}
