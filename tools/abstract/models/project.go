package models

import (
	"encoding/json"

	"github.com/Snow-Gremlin/Gozer/reader"
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type projectImp struct {
	lastId     uint64
	proj       *reader.Project
	packages   collections.List[PackageModel]
	interfaces collections.List[InterfaceModel]
	objects    collections.List[ObjectModel]
	signatures collections.List[SignatureModel]
	extraTypes collections.List[ExtraTypeModel]
	methods    collections.List[MethodModel]
}

func NewProject(proj *reader.Project) ProjectModel {
	return &projectImp{
		lastId:     0,
		proj:       proj,
		packages:   list.New[PackageModel](),
		interfaces: list.New[InterfaceModel](),
		objects:    list.New[ObjectModel](),
		signatures: list.New[SignatureModel](),
		extraTypes: list.New[ExtraTypeModel](),
		methods:    list.New[MethodModel](),
	}
}

func (imp *projectImp) nextId() uint64 {
	imp.lastId++
	return imp.lastId
}

func (imp *projectImp) Source() *reader.Project { return imp.proj }
func (imp *projectImp) Path() string            { return imp.proj.Packages[0].PkgPath }

func (imp *projectImp) Packages() collections.List[PackageModel]        { return imp.packages }
func (imp *projectImp) AllInterfaces() collections.List[InterfaceModel] { return imp.interfaces }
func (imp *projectImp) AllObjects() collections.List[ObjectModel]       { return imp.objects }
func (imp *projectImp) AllSignatures() collections.List[SignatureModel] { return imp.signatures }
func (imp *projectImp) ExtraTypes() collections.List[ExtraTypeModel]    { return imp.extraTypes }
func (imp *projectImp) Methods() collections.List[MethodModel]          { return imp.methods }

func (imp *projectImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addVal(data, `path`, imp.Path())
	addList(data, `packages`, imp.Packages())
	addList(data, `interfaces`, imp.AllInterfaces())
	addList(data, `objects`, imp.AllObjects())
	addList(data, `signatures`, imp.AllSignatures())
	addList(data, `extra`, imp.ExtraTypes())
	addList(data, `methods`, imp.Methods())
	return json.Marshal(data)
}
