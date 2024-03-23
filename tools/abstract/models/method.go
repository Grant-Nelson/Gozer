package models

import (
	"encoding/json"
)

type methodModelImp struct {
	id       uint64
	name     string
	exported bool
	sig      SignatureModel
}

func (pkg *packageImp) AddMethod(name string, exported bool, sig SignatureModel) MethodModel {
	imp := &methodModelImp{
		id:       pkg.Project().nextId(),
		name:     name,
		exported: exported,
		sig:      sig,
	}
	pkg.Project().Methods().Append(imp)
	pkg.Methods().Append(imp)
	return imp
}

func (imp *methodModelImp) Id() uint64     { return imp.id }
func (imp *methodModelImp) Name() string   { return imp.name }
func (imp *methodModelImp) Exported() bool { return imp.exported }

func (imp *methodModelImp) Signature() SignatureModel { return imp.sig }

func (imp *methodModelImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addVal(data, `id`, imp.Id())
	addVal(data, `name`, imp.Name())
	addVal(data, `exported`, imp.Exported())
	addVal(data, `sig`, imp.Signature().Id())
	return json.Marshal(data)
}
