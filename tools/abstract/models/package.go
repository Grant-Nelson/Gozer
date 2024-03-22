package models

import (
	"encoding/json"
	"go/token"

	"golang.org/x/tools/go/packages"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type packageImp struct {
	id         uint64
	proj       *projectImp
	pkg        *packages.Package
	interfaces collections.List[InterfaceModel]
	objects    collections.List[ObjectModel]
	methods    collections.List[MethodModel]
	statics    collections.List[TypeModel]
}

func (proj *projectImp) AddPackage(pkg *packages.Package) PackageModel {
	imp := &packageImp{
		id:         proj.nextId(),
		proj:       proj,
		pkg:        pkg,
		interfaces: list.New[InterfaceModel](),
		objects:    list.New[ObjectModel](),
		methods:    list.New[MethodModel](),
		statics:    list.New[TypeModel](),
	}
	proj.Packages().Append(imp)
	return imp
}

func (imp *packageImp) Id() uint64                   { return imp.id }
func (imp *packageImp) Name() string                 { return imp.pkg.Name }
func (imp *packageImp) Path() string                 { return imp.pkg.PkgPath }
func (imp *packageImp) Source() *packages.Package    { return imp.pkg }
func (imp *packageImp) Project() ProjectModel        { return imp.proj }
func (imp *packageImp) PosPath(pos token.Pos) string { return imp.pkg.Fset.Position(pos).String() }

func (imp *packageImp) Interfaces() collections.List[InterfaceModel] { return imp.interfaces }
func (imp *packageImp) Objects() collections.List[ObjectModel]       { return imp.objects }
func (imp *packageImp) Methods() collections.List[MethodModel]       { return imp.methods }
func (imp *packageImp) Statics() collections.List[TypeModel]         { return imp.statics }

func (imp *packageImp) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	addVal(data, `id`, imp.Id())
	addVal(data, `name`, imp.Name())
	addVal(data, `path`, imp.Path())
	addIds(data, `interfaces`, imp.Interfaces())
	addIds(data, `objects`, imp.Objects())
	addIds(data, `methods`, imp.Methods())
	addIds(data, `statics`, imp.Statics())
	return json.Marshal(data)
}
