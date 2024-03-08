package constructs

import "github.com/Snow-Gremlin/goToolbox/collections"

type (
	INamed interface {
		String() string
		Exported() bool
	}

	IProject interface {
		INamed
		Packages() collections.List[IPackages]
		AllPackages() collections.Enumerator[IPackages]
	}

	IPackages interface {
		INamed
		SubPackages() collections.List[IPackages]
		Interfaces() collections.List[IInterface]
		Objects() collections.List[IObject]
		Functions() collections.List[IMethod]
		Variables() collections.List[IVariable]
	}

	IVariable interface {
		INamed
		Constant() bool
		Type() IType
		Assignment() IExpression
	}

	IMethod interface {
		INamed
		Signature() ISignature
		Body() collections.List[IStatement]
	}

	IExpression interface{}

	IStatement interface{}
)
