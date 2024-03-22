package models

import (
	"encoding/json"
)

type methodModelImp struct {
	id   uint64
	name string
	sig  SignatureModel
}

func (pkg *packageImp) AddMethod(name string, sig SignatureModel) MethodModel {
	imp := &methodModelImp{
		id:   pkg.Project().nextId(),
		name: name,
		sig:  sig,
	}
	pkg.Project().Methods().Append(imp)
	pkg.Methods().Append(imp)
	return imp
}

func (imp *methodModelImp) Id() uint64                { return imp.id }
func (imp *methodModelImp) Name() string              { return imp.name }
func (imp *methodModelImp) Signature() SignatureModel { return imp.sig }

func (imp *methodModelImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addVal(data, `id`, imp.Id())
	addVal(data, `name`, imp.Name())
	addVal(data, `sig`, imp.Signature().Id())
	return json.Marshal(data)
}
