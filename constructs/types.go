package constructs

import (
	"strconv"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

type (
	IType interface {
		String() string
		_type()
	}

	IBoolType interface {
		IType
		_boolType()
	}

	IIntType interface {
		IType
		Signed() bool
		Size() int
		_intType()
	}

	IFloatType interface {
		IType
		Size() int
		_floatType()
	}

	IStringType interface {
		IType
		_stringType()
	}

	IPointerType interface {
		IType
		Inner() IType
		Unsafe() bool
		_pointerType()
	}

	IInterface interface {
		IType
		IGeneric
		Extends() collections.List[IInterface]
		Methods() collections.List[ISignature]
		_interfaceType()
	}

	IStruct interface {
		IType
		IGeneric
		Fields() collections.List[IParameter]
		_structType()
	}

	ISignature interface {
		IType
		IGeneric
		Variadic() bool
		Parameters() collections.List[IParameter]
		Returns() collections.List[IParameter]
		_signatureType()
	}

	IObject interface {
		IType
		IGeneric
		Implements() collections.List[IInterface]
		Methods() collections.List[IMethod]
		Data() IType
		_classType()
	}

	ITypeParameter interface {
		INamed
		Type() IInterface
	}

	IGeneric interface {
		TypeParameters() collections.List[ITypeParameter]
		IsGeneric() bool
	}

	IParameter interface {
		INamed
		Type() IType
	}
)

type intTypeImp struct {
	name   string
	signed bool
	size   int
}

func (imp *intTypeImp) Signed() bool   { return imp.signed }
func (imp *intTypeImp) Size() int      { return imp.size }
func (imp *intTypeImp) _type()         {}
func (imp *intTypeImp) _intType()      {}
func (imp *intTypeImp) String() string { return imp.name }

var (
	intInst    = &intTypeImp{name: `int`, signed: true, size: 64}
	int8Inst   = &intTypeImp{name: `int8`, signed: true, size: 8}
	int16Inst  = &intTypeImp{name: `int16`, signed: true, size: 16}
	int32Inst  = &intTypeImp{name: `int32`, signed: true, size: 32}
	int64Inst  = &intTypeImp{name: `int64`, signed: true, size: 64}
	uintInst   = &intTypeImp{name: `uint`, signed: false, size: 64}
	uint8Inst  = &intTypeImp{name: `uint8`, signed: false, size: 8}
	uint16Inst = &intTypeImp{name: `uint16`, signed: false, size: 16}
	uint32Inst = &intTypeImp{name: `uint32`, signed: false, size: 32}
	uint64Inst = &intTypeImp{name: `uint64`, signed: false, size: 64}
)

func Int() IIntType    { return intInst }
func Int8() IIntType   { return int8Inst }
func Int16() IIntType  { return int16Inst }
func Int32() IIntType  { return int32Inst }
func Int64() IIntType  { return int64Inst }
func Uint() IIntType   { return uintInst }
func Uint8() IIntType  { return uint8Inst }
func Uint16() IIntType { return uint16Inst }
func Uint32() IIntType { return uint32Inst }
func Uint64() IIntType { return uint64Inst }

type floatTypeImp struct {
	size int
}

func (imp *floatTypeImp) Size() int      { return imp.size }
func (imp *floatTypeImp) _type()         {}
func (imp *floatTypeImp) _floatType()    {}
func (imp *floatTypeImp) String() string { return `float` + strconv.Itoa(imp.size) }

var (
	float32Inst = &floatTypeImp{size: 32}
	float64Inst = &floatTypeImp{size: 64}
)

func Float32() IFloatType { return float32Inst }
func Float64() IFloatType { return float64Inst }

type boolTypeImp struct{}

func (imp *boolTypeImp) _type()         {}
func (imp *boolTypeImp) _boolType()     {}
func (imp *boolTypeImp) String() string { return `bool` }

var boolInst = &boolTypeImp{}

func Bool() IBoolType { return boolInst }

type stringTypeImp struct{}

func (imp *stringTypeImp) _type()         {}
func (imp *stringTypeImp) _stringType()   {}
func (imp *stringTypeImp) String() string { return `string` }

var stringInst = &stringTypeImp{}

func String() IStringType { return stringInst }

/*
	IPointerType interface {
		IType
		Inner() IType
		Unsafe() bool
		_pointerType()
	}

	IInterface interface {
		IType
		IGeneric
		Extends() collections.List[IInterface]
		Methods() collections.List[ISignature]
		_interfaceType()
	}

	IStruct interface {
		IType
		IGeneric
		Fields() collections.List[IParameter]
		_structType()
	}

	ISignature interface {
		IType
		IGeneric
		Variadic() bool
		Parameters() collections.List[IParameter]
		Returns() collections.List[IParameter]
		_signatureType()
	}

	IClass interface {
		IType
		IGeneric
		Implements() collections.List[IInterface]
		Methods() collections.List[IMethod]
		Data() IType
		_classType()
	}
*/

/*
	ITypeParameter interface {
		INamed
		Type() IInterface
	}

	IGeneric interface {
		TypeParameters() collections.List[ITypeParameter]
		IsGeneric() bool
	}

	IParameter interface {
		INamed
		Type() IType
	}
*/
