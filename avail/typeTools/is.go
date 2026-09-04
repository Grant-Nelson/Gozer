package typeTools

import "go/types"

func Is[T types.Type](t types.Type) bool {
	_, ok := t.(T)
	return ok
}

func IsArray(t types.Type) bool     { return Is[*types.Array](t) }
func IsBasic(t types.Type) bool     { return Is[*types.Basic](t) }
func IsChan(t types.Type) bool      { return Is[*types.Chan](t) }
func IsInterface(t types.Type) bool { return Is[*types.Interface](t) }
func IsMap(t types.Type) bool       { return Is[*types.Map](t) }
func IsNamed(t types.Type) bool     { return Is[*types.Named](t) }
func IsPointer(t types.Type) bool   { return Is[*types.Pointer](t) }
func IsSignature(t types.Type) bool { return Is[*types.Signature](t) }
func IsSlice(t types.Type) bool     { return Is[*types.Slice](t) }
func IsStruct(t types.Type) bool    { return Is[*types.Struct](t) }
func IsTuple(t types.Type) bool     { return Is[*types.Tuple](t) }
func IsTypeParam(t types.Type) bool { return Is[*types.TypeParam](t) }
func IsUnion(t types.Type) bool     { return Is[*types.Union](t) }

func IsBool(t types.Type) bool           { return BasicKind(t) == types.Bool }
func IsInt(t types.Type) bool            { return BasicKind(t) == types.Int }
func IsInt8(t types.Type) bool           { return BasicKind(t) == types.Int8 }
func IsInt16(t types.Type) bool          { return BasicKind(t) == types.Int16 }
func IsInt32(t types.Type) bool          { return BasicKind(t) == types.Int32 }
func IsInt64(t types.Type) bool          { return BasicKind(t) == types.Int64 }
func IsUint(t types.Type) bool           { return BasicKind(t) == types.Uint }
func IsUint8(t types.Type) bool          { return BasicKind(t) == types.Uint8 }
func IsUint16(t types.Type) bool         { return BasicKind(t) == types.Uint16 }
func IsUint32(t types.Type) bool         { return BasicKind(t) == types.Uint32 }
func IsUint64(t types.Type) bool         { return BasicKind(t) == types.Uint64 }
func IsUintptr(t types.Type) bool        { return BasicKind(t) == types.Uintptr }
func IsFloat32(t types.Type) bool        { return BasicKind(t) == types.Float32 }
func IsFloat64(t types.Type) bool        { return BasicKind(t) == types.Float64 }
func IsComplex64(t types.Type) bool      { return BasicKind(t) == types.Complex64 }
func IsComplex128(t types.Type) bool     { return BasicKind(t) == types.Complex128 }
func IsString(t types.Type) bool         { return BasicKind(t) == types.String }
func IsUnsafePointer(t types.Type) bool  { return BasicKind(t) == types.UnsafePointer }
func IsUntypedBool(t types.Type) bool    { return BasicKind(t) == types.UntypedBool }
func IsUntypedInt(t types.Type) bool     { return BasicKind(t) == types.UntypedInt }
func IsUntypedRune(t types.Type) bool    { return BasicKind(t) == types.UntypedRune }
func IsUntypedFloat(t types.Type) bool   { return BasicKind(t) == types.UntypedFloat }
func IsUntypedComplex(t types.Type) bool { return BasicKind(t) == types.UntypedComplex }
func IsUntypedString(t types.Type) bool  { return BasicKind(t) == types.UntypedString }
func IsUntypedNil(t types.Type) bool     { return BasicKind(t) == types.UntypedNil }

func IsBoolean(t types.Type) bool   { return BasicInfo(t)&types.IsBoolean != 0 }
func IsInteger(t types.Type) bool   { return BasicInfo(t)&types.IsInteger != 0 }
func IsUnsigned(t types.Type) bool  { return BasicInfo(t)&types.IsUnsigned != 0 }
func IsFloat(t types.Type) bool     { return BasicInfo(t)&types.IsFloat != 0 }
func IsComplex(t types.Type) bool   { return BasicInfo(t)&types.IsComplex != 0 }
func IsUntyped(t types.Type) bool   { return BasicInfo(t)&types.IsUntyped != 0 }
func IsOrdered(t types.Type) bool   { return BasicInfo(t)&types.IsOrdered != 0 }
func IsNumeric(t types.Type) bool   { return BasicInfo(t)&types.IsNumeric != 0 }
func IsConstType(t types.Type) bool { return BasicInfo(t)&types.IsConstType != 0 }

func IsPtrOf[T types.Type](t types.Type) bool {
	b, ok := t.(*types.Pointer)
	return ok && Is[T](b.Elem())
}

func IsArrayPtr(t types.Type) bool  { return IsPtrOf[*types.Array](t) }
func IsBasicPtr(t types.Type) bool  { return IsPtrOf[*types.Basic](t) }
func IsSlicePtr(t types.Type) bool  { return IsPtrOf[*types.Slice](t) }
func IsStructPtr(t types.Type) bool { return IsPtrOf[*types.Struct](t) }

func IsSliceOf[T types.Type](t types.Type) bool {
	s, ok := t.(*types.Slice)
	return ok && Is[T](s.Elem())
}
