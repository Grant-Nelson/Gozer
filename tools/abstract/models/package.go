package models

import (
	"encoding/json"

	"golang.org/x/tools/go/packages"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type packageImp struct {
	index      uint64
	pkg        *packages.Package
	interfaces collections.List[InterfaceModel]
	objects    collections.List[ObjectModel]
	methods    collections.List[MethodModel]
	statics    collections.List[TypeModel]
}

func NewPackage(pkg *packages.Package) PackageModel {
	return &packageImp{
		index:      0,
		pkg:        pkg,
		interfaces: list.New[InterfaceModel](),
		objects:    list.New[ObjectModel](),
		methods:    list.New[MethodModel](),
		statics:    list.New[TypeModel](),
	}
}

func (imp *packageImp) Index() uint64             { return imp.index }
func (imp *packageImp) setIndex(index uint64)     { imp.index = index }
func (imp *packageImp) Name() string              { return imp.pkg.Name }
func (imp *packageImp) Path() string              { return imp.pkg.PkgPath }
func (imp *packageImp) Source() *packages.Package { return imp.pkg }

func (imp *packageImp) Interfaces() collections.List[InterfaceModel] { return imp.interfaces }
func (imp *packageImp) Objects() collections.List[ObjectModel]       { return imp.objects }
func (imp *packageImp) Methods() collections.List[MethodModel]       { return imp.methods }
func (imp *packageImp) Statics() collections.List[TypeModel]         { return imp.statics }

func (imp *packageImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addString(data, `name`, imp.Name())
	addString(data, `path`, imp.Path())
	addIndex(data, `index`, imp.Index())
	addIndices(data, `interfaces`, imp.Interfaces())
	addIndices(data, `objects`, imp.Objects())
	addIndices(data, `methods`, imp.Methods())
	addTypeIndices(data, `statics`, imp.Statics())
	return json.Marshal(data)
}
