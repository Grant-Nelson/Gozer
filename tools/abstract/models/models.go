package models

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"golang.org/x/tools/go/packages"
)

type (
	ProjectModel interface {
		Path() string
		AllPackages() collections.List[PackageModel]
		AllInterfaces() collections.List[InterfaceModel]
		AllObjects() collections.List[ObjectModel]
		AllSignatures() collections.List[SignatureModel]
		AllMethods() collections.List[MethodModel]
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
		Path() string
		Source() *packages.Package
		Interfaces() collections.List[InterfaceModel]
		Objects() collections.List[ObjectModel]
		Methods() collections.List[MethodModel]
		Statics() collections.List[TypeModel]
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
		Implements() collections.List[InterfaceModel]
		Extends() collections.List[TypeModel]
		TypeParams() collections.List[TypeModel]
		Fields() collections.List[TypeModel]
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
		// TODO: LineCount & other metrics
	}
)
