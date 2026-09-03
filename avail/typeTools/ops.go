package typeTools

import (
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/typeTools/typeOp"
)

type OpTypes struct {
	Ops      typeOp.Op
	Key      types.Type
	Elem     types.Type
	Complex  types.Type
	RealImag types.Type
	Slice    types.Type
	Deref    types.Type
	Range1   types.Type
	Range2   types.Type
}

func IntersectOps(a, b OpTypes) OpTypes {
	// TODO: Implement
	return OpTypes{}
}

func Ops(t types.Type) OpTypes {
	switch t := t.Underlying().(type) {
	case *types.Array:
		return arrayTypeOps(t)
	case *types.Basic:
		return basicTypeOps(t)
	case *types.Chan:
		return chanTypeOps(t)
	case *types.Interface:
		return interfaceTypeOps(t)
	case *types.Map:
		return mapTypeOps(t)
	case *types.Pointer:
		return pointerTypeOps(t)
	case *types.Signature:
		return signatureTypeOps(t)
	case *types.Slice:
		return sliceTypeOps(t)
	case *types.Struct:
		return structTypeOps(t)
	case *types.Union:
		return unionTypeOps(t)
	}
	return OpTypes{}
}

func arrayTypeOps(t *types.Array) OpTypes {
	ops := OpTypes{
		Ops: typeOp.Clear | typeOp.GetIndex | typeOp.IsNil | typeOp.Len |
			typeOp.Make | typeOp.Make3 | typeOp.Range | typeOp.Range2 | typeOp.Ref |
			typeOp.RefIndex | typeOp.SetIndex | typeOp.Slice | typeOp.Slice3,
		Elem: t.Elem(),
	}
	if IsUint8(t.Elem()) {
		ops.Ops |= typeOp.ByteSlice
	}
	return ops
}

func basicTypeOps(t *types.Basic) OpTypes {
	switch t.Kind() {
	case types.Bool, types.UntypedBool:
		return booleanTypeOps()
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
		types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
		types.UntypedInt, types.UntypedRune, types.Uintptr:
		return integerTypeOps(t)
	case types.UntypedFloat, types.Float32, types.Float64:
		return floatTypeOps(t)
	case types.UntypedComplex, types.Complex64, types.Complex128:
		return complexTypeOps(t)
	case types.UntypedString, types.String:
		return stringTypeOps()
	case types.UnsafePointer:
		return unsafePointerTypeOps()
	case types.UntypedNil:
		return untypedNilTypeOps()
	}
	return OpTypes{}
}

func booleanTypeOps() OpTypes {
	return OpTypes{Ops: typeOp.Comparable | typeOp.Ref}
}

func integerTypeOps(t *types.Basic) OpTypes {
	return OpTypes{
		Ops: typeOp.Add | typeOp.Arith | typeOp.Bitwise | typeOp.Comparable |
			typeOp.Mod | typeOp.Orderable | typeOp.Range | typeOp.Ref,
		Range1: t,
	}
}

func floatTypeOps(t *types.Basic) OpTypes {
	ops := OpTypes{
		Ops: typeOp.Add | typeOp.Arith | typeOp.Comparable | typeOp.Complex |
			typeOp.Orderable | typeOp.Ref,
		Complex: types.Typ[types.UntypedComplex],
	}
	switch t.Kind() {
	case types.Float32:
		ops.Complex = types.Typ[types.Complex64]
	case types.Float64:
		ops.Complex = types.Typ[types.Complex128]
	default:
		ops.Complex = types.Typ[types.UntypedComplex]
	}
	return ops
}

func complexTypeOps(t *types.Basic) OpTypes {
	ops := OpTypes{
		Ops: typeOp.Add | typeOp.Arith | typeOp.Comparable | typeOp.RealImag |
			typeOp.Ref,
	}
	switch t.Kind() {
	case types.Complex64:
		ops.RealImag = types.Typ[types.Float32]
	case types.Complex128:
		ops.RealImag = types.Typ[types.Float64]
	default:
		ops.RealImag = types.Typ[types.UntypedFloat]
	}
	return ops
}

func stringTypeOps() OpTypes {
	return OpTypes{
		Ops: typeOp.Add | typeOp.ByteSlice | typeOp.GetIndex | typeOp.Len |
			typeOp.Comparable | typeOp.Orderable | typeOp.Range | typeOp.Range2 |
			typeOp.Ref | typeOp.Slice,
		Key:    types.Typ[types.UntypedInt],
		Elem:   types.Typ[types.Byte],
		Slice:  types.Typ[types.String],
		Range1: types.Typ[types.Int],
		Range2: types.Typ[types.Rune],
	}
}

func unsafePointerTypeOps() OpTypes {
	return OpTypes{
		Ops: typeOp.Comparable | typeOp.IsNil | typeOp.Orderable | typeOp.Ref,
	}
}

func untypedNilTypeOps() OpTypes {
	return OpTypes{Ops: typeOp.IsNil | typeOp.Ref}
}

func chanTypeOps(t *types.Chan) OpTypes {
	ops := OpTypes{
		Ops: typeOp.Len | typeOp.Cap | typeOp.IsNil | typeOp.Comparable,
	}
	switch t.Dir() {
	case types.SendOnly:
		ops.Ops |= typeOp.Send
	case types.RecvOnly:
		ops.Ops |= typeOp.Range | typeOp.Recv
		ops.Range1 = t.Elem()
	default:
		ops.Ops |= typeOp.Range | typeOp.Recv | typeOp.Send
		ops.Range1 = t.Elem()
	}
	return ops
}

func interfaceTypeOps(t *types.Interface) OpTypes {
	if t.IsMethodSet() {
		if t.IsComparable() {
			return OpTypes{Ops: typeOp.IsNil | typeOp.Comparable}
		}
		return OpTypes{Ops: typeOp.IsNil}
	}

	var ops OpTypes
	for i := range t.NumEmbeddeds() {
		u := Ops(t.EmbeddedType(i))
		if i == 0 {
			ops = u
		} else {
			ops = IntersectOps(ops, u)
		}
	}
	return ops
}

func mapTypeOps(t *types.Map) OpTypes {
	return OpTypes{
		Ops: typeOp.Cap | typeOp.Clear | typeOp.GetIndex | typeOp.IsNil | typeOp.Len |
			typeOp.Make | typeOp.Make3 | typeOp.Range | typeOp.Range2 | typeOp.Ref |
			typeOp.SetIndex,
		Key:    t.Key(),
		Elem:   t.Elem(),
		Range1: t.Key(),
		Range2: t.Key(),
	}
}

func pointerTypeOps(t *types.Pointer) OpTypes {
	var ops OpTypes
	if at, ok := t.Elem().Underlying().(*types.Array); ok {
		ops = arrayTypeOps(at)
	}
	ops.Ops |= typeOp.Comparable | typeOp.Deref | typeOp.IsNil |
		typeOp.Orderable | typeOp.Ref
	ops.Deref = t.Elem()
	return ops
}

func signatureTypeOps(t *types.Signature) OpTypes {
	// Check for iter.Seq and iter.Seq2 (receivers are optional).
	if (t.Params() != nil && t.Params().Len() == 1) &&
		(t.Results() == nil || t.Results().Len() == 0) {
		par := t.Params().At(0).Type()
		if fn, ok := par.(*types.Signature); ok &&
			fn.Params() != nil &&
			fn.Results() != nil && t.Results().Len() == 1 && IsBool(t.Results().At(0).Type()) {
			// at this point we know the signature looks like `func(func(<params>)bool)`
			switch t.Params().Len() {
			case 1:
				return OpTypes{
					Ops:    typeOp.IsNil | typeOp.Range | typeOp.Ref,
					Range1: t.Params().At(0).Type(),
				}
			case 2:
				return OpTypes{
					Ops:    typeOp.IsNil | typeOp.Range2 | typeOp.Ref,
					Range1: t.Params().At(0).Type(),
					Range2: t.Params().At(1).Type(),
				}
			}
		}
	}
	return OpTypes{Ops: typeOp.IsNil | typeOp.Ref}
}

func sliceTypeOps(t *types.Slice) OpTypes {
	const sliceOp = typeOp.Cap | typeOp.Clear | typeOp.GetIndex | typeOp.IsNil | typeOp.Len |
		typeOp.Make | typeOp.Make3 | typeOp.Range | typeOp.Range2 | typeOp.Ref | typeOp.RefIndex |
		typeOp.SetIndex | typeOp.Slice | typeOp.Slice3
	if IsUint8(t.Elem()) {
		return sliceOp | typeOp.ByteSlice
	}
	return sliceOp
}

func structTypeOps(t *types.Struct) OpTypes {
	if types.Comparable(t) {
		return typeOp.Comparable | typeOp.IsNil
	}
	return typeOp.IsNil
}

func unionTypeOps(t *types.Union) OpTypes {
	union := typeOp.None
	for i := range t.Len() {
		u := Ops(t.Term(i).Type())
		if i == 0 {
			union = u
		} else {
			// TODO: Need to check that the types match for ops too
			union &= u
		}
	}
	return union
}
