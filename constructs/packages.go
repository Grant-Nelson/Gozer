package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
)

type IPackage interface {
	INamed
	SubPackages() collections.List[IPackage]
	Interfaces() collections.List[IInterface]
	Objects() collections.List[IObject]
	Functions() collections.List[IMethod]
	Values() collections.List[IValue]
	_packageConstruct()
}

type packageImp struct {
	namedImp
	subPackages collections.List[IPackage]
	interfaces  collections.List[IInterface]
	objects     collections.List[IObject]
	functions   collections.List[IMethod]
	values      collections.List[IValue]
}

func (imp *packageImp) _packageConstruct() {}

func (imp *packageImp) SubPackages() collections.List[IPackage] {
	return imp.subPackages
}

func (imp *packageImp) Interfaces() collections.List[IInterface] {
	return imp.interfaces
}

func (imp *packageImp) Objects() collections.List[IObject] {
	return imp.objects
}

func (imp *packageImp) Functions() collections.List[IMethod] {
	return imp.functions
}

func (imp *packageImp) Values() collections.List[IValue] {
	return imp.values
}

func NewPackage(name string) IPackage {
	return &packageImp{
		namedImp:    newName(name),
		subPackages: list.New[IPackage](),
		interfaces:  list.New[IInterface](),
		objects:     list.New[IObject](),
		functions:   list.New[IMethod](),
		values:      list.New[IValue](),
	}
}
