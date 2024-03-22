package models

import (
	"go/token"

	"github.com/Snow-Gremlin/Gozer/reader"
	"github.com/Snow-Gremlin/goToolbox/collections"
	"golang.org/x/tools/go/packages"
)

type (
	// Identifier is part of a model which
	// identifies the model with a unique value.
	Identifier interface {
		// Id is the unique value for the model.
		// The zero Id indicates not set.
		Id() uint64
	}

	// TypeModel is a model which defines a type.
	TypeModel interface {
		Identifier
		_typeModel()
	}

	// ProjectModel is the main model containing
	// the information needed for an application.
	ProjectModel interface {
		nextId() uint64
		Source() *reader.Project
		Path() string

		Packages() collections.List[PackageModel]
		Interfaces() collections.List[InterfaceModel]
		Objects() collections.List[ObjectModel]
		Signatures() collections.List[SignatureModel]
		ExtraTypes() collections.List[ExtraTypeModel]
		Methods() collections.List[MethodModel]

		AddPackage(pkg *packages.Package) PackageModel
		AddSignature() SignatureModel
		AddExtraType(name string) ExtraTypeModel
	}

	PackageModel interface {
		Identifier
		Name() string
		Path() string
		Project() ProjectModel
		Source() *packages.Package
		PosPath(pos token.Pos) string

		Interfaces() collections.List[InterfaceModel]
		Objects() collections.List[ObjectModel]
		Methods() collections.List[MethodModel]
		Statics() collections.List[TypeModel]

		AddInterface(name string) InterfaceModel
		AddObject(name string) ObjectModel
		AddMethod(name string, sig SignatureModel) MethodModel
	}

	InterfaceModel interface {
		TypeModel
		Name() string
		TypeParams() collections.List[TypeModel]
		Signatures() collections.List[SignatureModel]
	}

	ObjectModel interface {
		TypeModel
		Name() string
		Extends() collections.List[TypeModel]
		TypeParams() collections.List[TypeModel]
		FieldTypes() collections.List[TypeModel]
	}

	SignatureModel interface {
		TypeModel
		TypeParams() collections.List[TypeModel]
		Params() collections.List[TypeModel]
		Returns() collections.List[TypeModel]
	}

	ExtraTypeModel interface {
		TypeModel
		Name() string
		Extends() collections.List[TypeModel]
	}

	MethodModel interface {
		Identifier
		Name() string
		Signature() SignatureModel
		// TODO: Writes
		// TODO: Reads
		// TODO: LineCount & other metrics
	}
)
