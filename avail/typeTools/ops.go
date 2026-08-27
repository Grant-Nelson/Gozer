package typeTools

import (
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/typeTools/typeOp"
)

func Ops(t types.Type) typeOp.Op {
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
		return mapTypeOps()
	case *types.Pointer:
		return pointerTypeOps(t)
	case *types.Signature:
		return signatureTypeOps()
	case *types.Slice:
		return sliceTypeOps(t)
	case *types.Struct:
		return structTypeOps(t)
	case *types.Union:
		return unionTypeOps(t)
	}
	return typeOp.None
}

func arrayTypeOps(t *types.Array) typeOp.Op {
	const arrayOp = typeOp.Clear | typeOp.GetIndex | typeOp.IsNil | typeOp.Len |
		typeOp.Make | typeOp.Make3 | typeOp.Ref | typeOp.RefIndex |
		typeOp.SetIndex | typeOp.Slice | typeOp.Slice3
	if IsUint8(t.Elem()) {
		return arrayOp | typeOp.ByteSlice
	}
	return arrayOp
}

func basicTypeOps(t *types.Basic) typeOp.Op {
	switch t.Kind() {
	case types.Bool, types.UntypedBool:
		return typeOp.Comparable | typeOp.Ref

	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
		types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
		types.UntypedInt, types.UntypedRune, types.Uintptr:
		return typeOp.Add | typeOp.Arith | typeOp.Bitwise | typeOp.Comparable |
			typeOp.Mod | typeOp.Orderable | typeOp.Ref

	case types.UntypedFloat, types.Float32, types.Float64:
		return typeOp.Add | typeOp.Arith | typeOp.Comparable | typeOp.Complex |
			typeOp.Orderable | typeOp.Ref

	case types.UntypedComplex, types.Complex64, types.Complex128:
		return typeOp.Add | typeOp.Arith | typeOp.Comparable | typeOp.RealImag |
			typeOp.Ref

	case types.UntypedString, types.String:
		return typeOp.Add | typeOp.ByteSlice | typeOp.GetIndex | typeOp.Len |
			typeOp.Comparable | typeOp.Orderable | typeOp.Ref | typeOp.Slice

	case types.UnsafePointer:
		return typeOp.Comparable | typeOp.IsNil | typeOp.Orderable | typeOp.Ref

	case types.UntypedNil:
		return typeOp.IsNil | typeOp.Ref
	}
	return typeOp.None
}

func chanTypeOps(t *types.Chan) typeOp.Op {
	const chanOp = typeOp.Len | typeOp.Cap | typeOp.IsNil | typeOp.Comparable
	switch t.Dir() {
	case types.SendOnly:
		return chanOp | typeOp.Send
	case types.RecvOnly:
		return chanOp | typeOp.Recv
	default:
		return chanOp | typeOp.Recv | typeOp.Send
	}
}

func interfaceTypeOps(t *types.Interface) typeOp.Op {
	if t.IsMethodSet() {
		if t.IsComparable() {
			return typeOp.IsNil | typeOp.Comparable
		}
		return typeOp.IsNil
	}

	union := typeOp.None
	for i := range t.NumEmbeddeds() {
		u := Ops(t.EmbeddedType(i))
		if i == 0 {
			union = u
		} else {
			union &= u
		}
	}
	return union
}

func mapTypeOps() typeOp.Op {
	return typeOp.Cap | typeOp.Clear | typeOp.GetIndex | typeOp.IsNil | typeOp.Len |
		typeOp.Make | typeOp.Make3 | typeOp.Ref | typeOp.SetIndex
}

func pointerTypeOps(t *types.Pointer) typeOp.Op {
	const ptrOp = typeOp.Comparable | typeOp.Deref | typeOp.IsNil |
		typeOp.Orderable | typeOp.Ref
	if at, ok := t.Elem().Underlying().(*types.Array); ok {
		return arrayTypeOps(at) | ptrOp
	}
	return ptrOp
}

func signatureTypeOps() typeOp.Op {
	return typeOp.IsNil | typeOp.Ref
}

func sliceTypeOps(t *types.Slice) typeOp.Op {
	const sliceOp = typeOp.Cap | typeOp.Clear | typeOp.GetIndex | typeOp.IsNil | typeOp.Len |
		typeOp.Make | typeOp.Make3 | typeOp.Ref | typeOp.RefIndex |
		typeOp.SetIndex | typeOp.Slice | typeOp.Slice3
	if IsUint8(t.Elem()) {
		return sliceOp | typeOp.ByteSlice
	}
	return sliceOp
}

func structTypeOps(t *types.Struct) typeOp.Op {
	if types.Comparable(t) {
		return typeOp.Comparable | typeOp.IsNil
	}
	return typeOp.IsNil
}

func unionTypeOps(t *types.Union) typeOp.Op {
	union := typeOp.None
	for i := range t.Len() {
		u := Ops(t.Term(i).Type())
		if i == 0 {
			union = u
		} else {
			union &= u
		}
	}
	return union
}
