package models

import "github.com/Snow-Gremlin/goToolbox/collections"

type (
	ProjectModel interface {
		Path() string
		AllPackages() collections.List[PackageModel]
		AllInterfaces() collections.List[InterfaceModel]
		AllObjects() collections.List[ObjectModel]
		AllSignatures() collections.List[SignatureModel]
		AllMethods() collections.List[MethodModel]
		AllFields() collections.List[FieldModel]
	}

	IndexedModel interface {
		Index() uint64
		setIndex(index uint64)
	}

	TypeModel interface {
		TypeIndex() uint64
		setTypeIndex(index uint64)
	}

	NamedModel interface {
		Name() string
	}

	PackageModel interface {
		IndexedModel
		NamedModel
		Parent() PackageModel
		Interfaces() collections.List[InterfaceModel]
		Objects() collections.List[ObjectModel]
		Methods() collections.List[MethodModel]
		StaticData() collections.List[FieldModel]
	}

	InterfaceModel interface {
		IndexedModel
		TypeModel
		NamedModel
		TypeParams() collections.List[TypeModel]
		Signatures() collections.List[SignatureModel]
	}

	ObjectModel interface {
		IndexedModel
		TypeModel
		NamedModel
		TypeParams() collections.List[TypeModel]
		Fields() collections.List[FieldModel]
	}

	SignatureModel interface {
		IndexedModel
		TypeModel
		NamedModel
		TypeParams() collections.List[TypeModel]
		Params() collections.List[TypeModel]
		Returns() collections.List[TypeModel]
	}

	MethodModel interface {
		IndexedModel
		NamedModel
		Signature() SignatureModel
		// TODO: Writes
		// TODO: Reads
	}

	FieldModel interface {
		IndexedModel
		NamedModel
		Type() TypeModel
	}
)
