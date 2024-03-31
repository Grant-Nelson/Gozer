package models

import (
	"go/token"
	"go/types"

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

	// ProjectModel is the main model containing
	// the information needed for an application.
	ProjectModel interface {
		nextId() uint64
		Source() *reader.Project
		Path() string

		Packages() collections.List[PackageModel]
		AllInterfaces() collections.List[InterfaceModel]
		AllObjects() collections.List[ObjectModel]
		AllSignatures() collections.List[SignatureModel]
		ExtraTypes() collections.List[ExtraTypeModel]
		Methods() collections.List[MethodModel]

		AddPackage(pkg *packages.Package) PackageModel
		AddSignature(name string, exported bool) SignatureModel
		AddExtraType(name string, exported bool) ExtraTypeModel
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

		AddInterface(obj types.Object) InterfaceModel
		AddObject(obj types.Object) ObjectModel
		AddMethod(name string, exported bool, sig SignatureModel) MethodModel
	}

	// TypeModel is a model which defines a type.
	TypeModel interface {
		Identifier
		Name() string
		Exported() bool
		TypeParams() collections.List[TypeModel]
		Extends() collections.List[TypeModel]
		_typeModel()
	}

	InterfaceModel interface {
		TypeModel
		Signatures() collections.List[SignatureModel]
	}

	ObjectModel interface {
		TypeModel
		FieldTypes() collections.List[TypeModel]
	}

	SignatureModel interface {
		TypeModel
		Params() collections.List[TypeModel]
		Returns() collections.List[TypeModel]
	}

	ExtraTypeModel interface {
		TypeModel
		_extraType()
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
