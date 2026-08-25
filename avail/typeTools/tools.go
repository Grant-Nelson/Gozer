package typeTools

import "go/types"

func Unalias(t types.Type) types.Type {
	return types.Unalias(t)
}

func Is[T types.Type](t types.Type) bool {
	_, ok := t.(T)
	return ok
}

func IsBasic(t types.Type) bool   { return Is[*types.Basic](t) }
func IsSlice(t types.Type) bool   { return Is[*types.Slice](t) }
func IsArray(t types.Type) bool   { return Is[*types.Array](t) }
func IsMap(t types.Type) bool     { return Is[*types.Map](t) }
func IsStruct(t types.Type) bool  { return Is[*types.Struct](t) }
func IsPointer(t types.Type) bool { return Is[*types.Pointer](t) }

func BasicKind(t types.Type) types.BasicKind {
	if b, ok := t.(*types.Basic); ok {
		return b.Kind()
	}
	return types.Invalid
}

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

const noInfo = types.BasicInfo(0)

func BasicInfo(t types.Type) types.BasicInfo {
	if b, ok := t.(*types.Basic); ok {
		return b.Info()
	}
	return noInfo
}

func IsBoolean(t types.Type) bool   { return BasicInfo(t)&types.IsBoolean != noInfo }
func IsInteger(t types.Type) bool   { return BasicInfo(t)&types.IsInteger != noInfo }
func IsUnsigned(t types.Type) bool  { return BasicInfo(t)&types.IsUnsigned != noInfo }
func IsFloat(t types.Type) bool     { return BasicInfo(t)&types.IsFloat != noInfo }
func IsComplex(t types.Type) bool   { return BasicInfo(t)&types.IsComplex != noInfo }
func IsUntyped(t types.Type) bool   { return BasicInfo(t)&types.IsUntyped != noInfo }
func IsOrdered(t types.Type) bool   { return BasicInfo(t)&types.IsOrdered != noInfo }
func IsNumeric(t types.Type) bool   { return BasicInfo(t)&types.IsNumeric != noInfo }
func IsConstType(t types.Type) bool { return BasicInfo(t)&types.IsConstType != noInfo }

func IsPtrOf[T types.Type](t types.Type) bool {
	b, ok := t.(*types.Pointer)
	return ok && Is[T](b.Elem())
}

func IsBasicPtr(t types.Type) bool  { return IsPtrOf[*types.Basic](t) }
func IsSlicePtr(t types.Type) bool  { return IsPtrOf[*types.Slice](t) }
func IsArrayPtr(t types.Type) bool  { return IsPtrOf[*types.Array](t) }
func IsStructPtr(t types.Type) bool { return IsPtrOf[*types.Struct](t) }
