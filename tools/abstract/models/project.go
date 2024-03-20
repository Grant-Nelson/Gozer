package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type projectImp struct {
	path          string
	allPackages   collections.List[PackageModel]
	allInterfaces collections.List[InterfaceModel]
	allObjects    collections.List[ObjectModel]
	allSignatures collections.List[SignatureModel]
	allMethods    collections.List[MethodModel]
}

func NewProject(path string) ProjectModel {
	return &projectImp{
		path:          path,
		allPackages:   list.New[PackageModel](),
		allInterfaces: list.New[InterfaceModel](),
		allObjects:    list.New[ObjectModel](),
		allSignatures: list.New[SignatureModel](),
		allMethods:    list.New[MethodModel](),
	}
}

func (imp *projectImp) Path() string { return imp.path }

func (imp *projectImp) AllPackages() collections.List[PackageModel]     { return imp.allPackages }
func (imp *projectImp) AllInterfaces() collections.List[InterfaceModel] { return imp.allInterfaces }
func (imp *projectImp) AllObjects() collections.List[ObjectModel]       { return imp.allObjects }
func (imp *projectImp) AllSignatures() collections.List[SignatureModel] { return imp.allSignatures }
func (imp *projectImp) AllMethods() collections.List[MethodModel]       { return imp.allMethods }

func (imp *projectImp) MarshalJSON() ([]byte, error) {
	setIndices(imp.allPackages)
	setIndices(imp.allInterfaces)
	setIndices(imp.allObjects)
	setIndices(imp.allSignatures)
	setIndices(imp.allMethods)

	var typeIndex uint64
	setTypeIndices(&typeIndex, imp.allInterfaces)
	setTypeIndices(&typeIndex, imp.allObjects)
	setTypeIndices(&typeIndex, imp.allSignatures)

	data := map[string]any{`path`: imp.path}
	addData(data, `packages`, imp.AllPackages())
	addData(data, `interfaces`, imp.AllInterfaces())
	addData(data, `objects`, imp.AllObjects())
	addData(data, `signatures`, imp.AllSignatures())
	addData(data, `methods`, imp.AllMethods())
	return json.Marshal(data)
}
