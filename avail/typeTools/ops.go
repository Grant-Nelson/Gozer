package typeTools

import (
	"go/types"
	"slices"
	"strings"

	"golang.org/x/exp/typeparams"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/typeTools/typeOp"
)

type OpTypes struct {
	Ops typeOp.Op

	// Key is the type for the parameter of the ops
	// GetIndex, GetIndex2, RefIndex, and SetIndex.
	Key types.Type

	// Elem is the type for the parameter of the SetIndex op,
	// the return types of the ops GetIndex and GetIndex2, and
	// the element type in the pointer returned by the RefIndex op.
	Elem types.Type

	// Complex is the returned type from the Complex op.
	Complex types.Type

	// RealImag is the returned type from the RealImag op.
	RealImag types.Type

	// Slice is the returned type from the ops Slice and Slice3.
	Slice types.Type

	// Deref is the returned type from the Deref op.
	Deref types.Type

	// Range1 is the returned type from a Range op and
	// the first returned type from a Range2 op.
	Range1 types.Type

	// Range2 is the second returned type from a Range2 op.
	Range2 types.Type
}

func (ops OpTypes) String() string {
	parts := []string{}
	add := func(name string, t types.Type) {
		if t != nil {
			parts = append(parts, name+`:`+t.String())
		}
	}
	add(`Key`, ops.Key)
	add(`Elem`, ops.Elem)
	add(`Complex`, ops.Complex)
	add(`RealImag`, ops.RealImag)
	add(`Slice`, ops.Slice)
	add(`Deref`, ops.Deref)
	add(`Range1`, ops.Range1)
	add(`Range2`, ops.Range2)
	tail := ``
	if len(parts) > 0 {
		tail = `{ ` + strings.Join(parts, `, `) + ` }`
	}
	return ops.Ops.String() + tail
}

func Ops(t types.Type) OpTypes {
	return getOps(t, t)
}

func getOps(orig, t types.Type) OpTypes {
	switch t2 := t.Underlying().(type) {
	case *types.Array:
		return arrayTypeOps(t2)
	case *types.Basic:
		return basicTypeOps(orig, t2)
	case *types.Chan:
		return chanTypeOps(t2)
	case *types.Interface:
		return interfaceTypeOps(orig, t2)
	case *types.Map:
		return mapTypeOps(t2)
	case *types.Pointer:
		return pointerTypeOps(t2)
	case *types.Signature:
		return signatureTypeOps(t2)
	case *types.Slice:
		return sliceTypeOps(orig, t2)
	case *types.Struct:
		return structTypeOps(t2)
	case *types.Union:
		return unionTypeOps(orig, t2)
	}
	return OpTypes{}
}

func arrayTypeOps(t2 *types.Array) OpTypes {
	ops := OpTypes{
		Ops: typeOp.Clear | typeOp.GetIndex | typeOp.IsNil | typeOp.Len |
			typeOp.Make | typeOp.Make3 | typeOp.Range | typeOp.Range2 | typeOp.Ref |
			typeOp.RefIndex | typeOp.SetIndex | typeOp.Slice | typeOp.Slice3,
		Key:    types.Typ[types.UntypedInt],
		Elem:   t2.Elem(),
		Slice:  types.NewSlice(t2.Elem()),
		Range1: types.Typ[types.Int],
		Range2: t2.Elem(),
	}
	if IsUint8(t2.Elem()) {
		ops.Ops |= typeOp.ByteSlice
	}
	return ops
}

func basicTypeOps(orig types.Type, t2 *types.Basic) OpTypes {
	switch t2.Kind() {
	case types.Bool, types.UntypedBool:
		return booleanTypeOps()
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
		types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
		types.UntypedInt, types.UntypedRune, types.Uintptr:
		return integerTypeOps(t2)
	case types.UntypedFloat, types.Float32, types.Float64:
		return floatTypeOps(t2)
	case types.UntypedComplex, types.Complex64, types.Complex128:
		return complexTypeOps(t2)
	case types.UntypedString, types.String:
		return stringTypeOps(orig)
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

func integerTypeOps(t2 *types.Basic) OpTypes {
	return OpTypes{
		Ops: typeOp.Add | typeOp.Arith | typeOp.Bitwise | typeOp.Comparable |
			typeOp.Mod | typeOp.Orderable | typeOp.Range | typeOp.Ref,
		Range1: t2,
	}
}

func floatTypeOps(t2 *types.Basic) OpTypes {
	ops := OpTypes{
		Ops: typeOp.Add | typeOp.Arith | typeOp.Comparable | typeOp.Complex |
			typeOp.Orderable | typeOp.Ref,
		Complex: types.Typ[types.UntypedComplex],
	}
	switch t2.Kind() {
	case types.Float32:
		ops.Complex = types.Typ[types.Complex64]
	case types.Float64:
		ops.Complex = types.Typ[types.Complex128]
	default:
		ops.Complex = types.Typ[types.UntypedComplex]
	}
	return ops
}

func complexTypeOps(t2 *types.Basic) OpTypes {
	ops := OpTypes{
		Ops: typeOp.Add | typeOp.Arith | typeOp.Comparable | typeOp.RealImag |
			typeOp.Ref,
	}
	switch t2.Kind() {
	case types.Complex64:
		ops.RealImag = types.Typ[types.Float32]
	case types.Complex128:
		ops.RealImag = types.Typ[types.Float64]
	default:
		ops.RealImag = types.Typ[types.UntypedFloat]
	}
	return ops
}

func stringTypeOps(orig types.Type) OpTypes {
	return OpTypes{
		Ops: typeOp.Add | typeOp.ByteSlice | typeOp.GetIndex | typeOp.Len |
			typeOp.Comparable | typeOp.Orderable | typeOp.Range | typeOp.Range2 |
			typeOp.Ref | typeOp.Slice,
		Key:    types.Typ[types.UntypedInt],
		Elem:   types.Typ[types.Byte],
		Slice:  orig,
		Range1: types.Typ[types.Int],
		Range2: RuneType(),
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

func chanTypeOps(t2 *types.Chan) OpTypes {
	ops := OpTypes{
		Ops: typeOp.Len | typeOp.Cap | typeOp.IsNil | typeOp.Comparable,
	}
	switch t2.Dir() {
	case types.SendOnly:
		ops.Ops |= typeOp.Send
	case types.RecvOnly:
		ops.Ops |= typeOp.Range | typeOp.Recv
		ops.Range1 = t2.Elem()
	default:
		ops.Ops |= typeOp.Range | typeOp.Recv | typeOp.Send
		ops.Range1 = t2.Elem()
	}
	return ops
}

func interfaceTypeOps(orig types.Type, t2 *types.Interface) OpTypes {
	if t2.IsMethodSet() {
		if t2.IsComparable() {
			return OpTypes{Ops: typeOp.IsNil | typeOp.Comparable}
		}
		return OpTypes{Ops: typeOp.IsNil}
	}

	terms, err := typeparams.NormalTerms(t2)
	if err != nil {
		panic(faults.New(`failed to determine op type of interface`, err).
			With(`interface`, t2))
	}
	return unionTermsOps(orig, terms)
}

func mapTypeOps(t2 *types.Map) OpTypes {
	return OpTypes{
		Ops: typeOp.Cap | typeOp.Clear | typeOp.GetIndex | typeOp.IsNil | typeOp.Len |
			typeOp.Make | typeOp.Make3 | typeOp.Range | typeOp.Range2 | typeOp.Ref |
			typeOp.SetIndex,
		Key:    t2.Key(),
		Elem:   t2.Elem(),
		Range1: t2.Key(),
		Range2: t2.Key(),
	}
}

func pointerTypeOps(t2 *types.Pointer) OpTypes {
	var ops OpTypes
	if at, ok := t2.Elem().Underlying().(*types.Array); ok {
		ops = arrayTypeOps(at)
	}
	ops.Ops |= typeOp.Comparable | typeOp.Deref | typeOp.IsNil |
		typeOp.Orderable | typeOp.Ref
	ops.Deref = t2.Elem()
	return ops
}

func signatureTypeOps(t2 *types.Signature) OpTypes {
	// Check for iter.Seq and iter.Seq2 (receivers are optional).
	if t2.Params().Len() == 1 && t2.Results().Len() == 0 {
		par := t2.Params().At(0).Type()
		if fn, ok := par.(*types.Signature); ok &&
			fn.Results().Len() == 1 && IsBool(fn.Results().At(0).Type()) {
			// At this point we know the signature looks like `func(func(<params>)bool)`.
			switch fn.Params().Len() {
			case 1:
				return OpTypes{
					Ops:    typeOp.IsNil | typeOp.Range | typeOp.Ref,
					Range1: fn.Params().At(0).Type(),
				}
			case 2:
				return OpTypes{
					Ops:    typeOp.IsNil | typeOp.Range2 | typeOp.Ref,
					Range1: fn.Params().At(0).Type(),
					Range2: fn.Params().At(1).Type(),
				}
			}
		}
	}
	return OpTypes{Ops: typeOp.IsNil | typeOp.Ref}
}

func sliceTypeOps(orig types.Type, t2 *types.Slice) OpTypes {
	ops := OpTypes{
		Ops: typeOp.Cap | typeOp.Clear | typeOp.GetIndex | typeOp.IsNil | typeOp.Len |
			typeOp.Make | typeOp.Make3 | typeOp.Range | typeOp.Range2 | typeOp.Ref | typeOp.RefIndex |
			typeOp.SetIndex | typeOp.Slice | typeOp.Slice3,
		Key:    types.Typ[types.UntypedInt],
		Elem:   t2.Elem(),
		Slice:  orig,
		Range1: types.Typ[types.Int],
		Range2: t2.Elem(),
	}
	if IsUint8(t2.Elem()) {
		ops.Ops |= typeOp.ByteSlice
	}
	return ops
}

func structTypeOps(t2 *types.Struct) OpTypes {
	if types.Comparable(t2) {
		return OpTypes{Ops: typeOp.Comparable | typeOp.IsNil}
	}
	return OpTypes{Ops: typeOp.IsNil}
}

func unionTypeOps(orig types.Type, t2 *types.Union) OpTypes {
	return unionTermsOps(orig, slices.Collect(t2.Terms()))
}

func termOps(t *types.Term) OpTypes {
	if t.Tilde() {
		return getOps(types.NewUnion([]*types.Term{t}), t.Type())
	}
	return getOps(t.Type(), t.Type())
}

func unionTermsOps(orig types.Type, t2 []*types.Term) OpTypes {
	switch len(t2) {
	case 0:
		return OpTypes{}
	case 1:
		getOps(orig, t2[0].Type())
	}

	terms := make([]OpTypes, len(t2))
	ops := typeOp.None
	for i, t := range t2 {
		u := termOps(t)
		terms[i] = u
		if i == 0 {
			ops = u.Ops
		} else {
			ops &= u.Ops
		}
	}
	if ops == typeOp.None {
		return OpTypes{}
	}

	adj := func(adjOps typeOp.Op, getType func(OpTypes) types.Type) types.Type {
		if ops.Any(adjOps) {
			if result := innerUnionType(terms, getType); result != nil {
				return result
			}
			ops &^= adjOps
		}
		return nil
	}

	return OpTypes{
		Ops:      ops,
		Key:      adj(typeOp.GetIndex|typeOp.GetIndex2|typeOp.RefIndex|typeOp.SetIndex, func(op OpTypes) types.Type { return op.Key }),
		Elem:     adj(typeOp.GetIndex|typeOp.GetIndex2|typeOp.RefIndex|typeOp.SetIndex, func(op OpTypes) types.Type { return op.Elem }),
		Complex:  adj(typeOp.Complex, func(op OpTypes) types.Type { return op.Complex }),
		RealImag: adj(typeOp.RealImag, func(op OpTypes) types.Type { return op.RealImag }),
		Slice:    adj(typeOp.Slice|typeOp.Slice3, func(op OpTypes) types.Type { return op.Slice }),
		Deref:    adj(typeOp.Deref, func(op OpTypes) types.Type { return op.Deref }),
		Range1:   adj(typeOp.Range|typeOp.Range2, func(op OpTypes) types.Type { return op.Range1 }),
		Range2:   adj(typeOp.Range2, func(op OpTypes) types.Type { return op.Range2 }),
	}
}

func innerUnionType(terms []OpTypes, getType func(OpTypes) types.Type) types.Type {
	ts := make([]*types.Term, 0, len(terms))
	for _, term := range terms {
		t := getType(term)
		if u, ok := t.(*types.Union); ok {
			ts = slices.AppendSeq(ts, u.Terms())
		} else {
			ts = append(ts, types.NewTerm(false, t))
		}
	}

	it := types.NewUnion(ts)
	t2, err := typeparams.NormalTerms(it)
	if err != nil {
		panic(faults.New(`failed to determine inner op type`, err).
			With(`union`, it))
	}

	switch len(t2) {
	case 0:
		return anyType()
	case 1:
		if !t2[0].Tilde() {
			return t2[0].Type()
		}
	}
	return types.NewUnion(t2)
}
