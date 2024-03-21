package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type projectImp struct {
	path     string
	packages collections.List[PackageModel]
	types    TypeCollectionModel
	methods  collections.List[MethodModel]
}

func NewProject(path string) ProjectModel {
	return &projectImp{
		path:     path,
		packages: list.New[PackageModel](),
		types:    NewTypeCollection(),
		methods:  list.New[MethodModel](),
	}
}

func (imp *projectImp) Path() string { return imp.path }

func (imp *projectImp) Packages() collections.List[PackageModel] { return imp.packages }
func (imp *projectImp) Types() TypeCollectionModel               { return imp.types }
func (imp *projectImp) Methods() collections.List[MethodModel]   { return imp.methods }

func (imp *projectImp) MarshalJSON() ([]byte, error) {
	setIndices(imp.packages)
	setIndices(imp.methods)

	data := map[string]any{}
	addString(data, `path`, imp.path)
	addData(data, `packages`, imp.Packages())
	data[`types`] = imp.Types()
	addData(data, `methods`, imp.Methods())
	return json.Marshal(data)
}
