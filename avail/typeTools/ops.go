package typeTools

import (
	"go/types"
	"strings"

	"golang.org/x/exp/typeparams"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/typeTools/typeOp"
)

type OpTypes struct {
	Ops typeOp.Op

	// Type is the type this OpTypes was built for.
	Type types.Type

	// TODO: Determine if we should add an indication if the type will always be boxed with the type
	// or if the underlying type is consistent (add the underlying as a type)
	// or if a type dictionary would been to be passed along side the type.

	// Key is the type for the parameter of the ops
	// GetIndex, GetIndex2, RefIndex, and SetIndex.
	Key types.Type

	// Elem is the type for the parameter of the SetIndex op,
	// the return types of the ops GetIndex and GetIndex2,
	// the value that can be sent or received from a channel, and
	// the element type in the pointer returned by the RefIndex op.
	Elem types.Type

	// Complex is the returned type from the Float op.
	Complex types.Type

	// Float is the returned type from the Complex op.
	Float types.Type

	// Slice is the returned type from the ops Slice and Slice3.
	Slice types.Type

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
			str := t.String()
			if strings.Contains(str, ` | `) {
				str = `{` + str + `}`
			}
			parts = append(parts, name+`:`+str)
		}
	}
	add(`Key`, ops.Key)
	add(`Elem`, ops.Elem)
	add(`Complex`, ops.Complex)
	add(`Float`, ops.Float)
	add(`Slice`, ops.Slice)
	add(`Range1`, ops.Range1)
	add(`Range2`, ops.Range2)
	tail := ``
	if len(parts) > 0 {
		tail = `{ ` + strings.Join(parts, `, `) + ` }`
	}
	return ops.Type.String() + `: ` + ops.Ops.String() + tail
}

func Ops(t types.Type) OpTypes {
	return getOps(t, t)
}

func blankOps(orig types.Type) OpTypes {
	return OpTypes{
		Ops:  typeOp.None,
		Type: orig,
	}
}

func getOps(orig, t types.Type) OpTypes {
	switch t2 := t.Underlying().(type) {
	case *types.Array:
		return arrayTypeOps(orig, t2)
	case *types.Basic:
		return basicTypeOps(orig, t2)
	case *types.Chan:
		return chanTypeOps(orig, t2)
	case *types.Interface:
		return interfaceTypeOps(orig, t2)
	case *types.Map:
		return mapTypeOps(orig, t2)
	case *types.Pointer:
		return pointerTypeOps(orig, t2)
	case *types.Signature:
		return signatureTypeOps(orig, t2)
	case *types.Slice:
		return sliceTypeOps(orig, t2)
	case *types.Struct:
		return structTypeOps(orig, t2)
	case *types.Union:
		return unionTypeOps(orig, t2)
	}
	return blankOps(orig)
}

func arrayTypeOps(orig types.Type, t2 *types.Array) OpTypes {
	return OpTypes{
		Ops: typeOp.Cap | typeOp.Comparable | typeOp.GetIndex | typeOp.Len |
			typeOp.Range0 | typeOp.Range1 | typeOp.Range2 | typeOp.Ref |
			typeOp.RefIndex | typeOp.SetIndex | typeOp.Slice | typeOp.Slice3,
		Type:   orig,
		Key:    types.Typ[types.UntypedInt],
		Elem:   t2.Elem(),
		Slice:  types.NewSlice(t2.Elem()),
		Range1: types.Typ[types.Int],
		Range2: t2.Elem(),
	}
}

func basicTypeOps(orig types.Type, t2 *types.Basic) OpTypes {
	switch t2.Kind() {
	case types.Bool, types.UntypedBool:
		return booleanTypeOps(orig)
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
		types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
		types.UntypedInt, types.UntypedRune, types.Uintptr:
		return integerTypeOps(orig)
	case types.UntypedFloat, types.Float32, types.Float64:
		return floatTypeOps(orig, t2)
	case types.UntypedComplex, types.Complex64, types.Complex128:
		return complexTypeOps(orig, t2)
	case types.UntypedString, types.String:
		return stringTypeOps(orig)
	case types.UnsafePointer:
		return unsafePointerTypeOps(orig)
	case types.UntypedNil:
		return untypedNilTypeOps(orig)
	}
	return blankOps(orig)
}

func booleanTypeOps(orig types.Type) OpTypes {
	return OpTypes{
		Ops:  typeOp.Comparable | typeOp.Ref,
		Type: orig,
	}
}

func integerTypeOps(orig types.Type) OpTypes {
	return OpTypes{
		Ops: typeOp.Add | typeOp.Arith | typeOp.Bitwise | typeOp.Comparable |
			typeOp.Mod | typeOp.Orderable | typeOp.Range0 | typeOp.Range1 | typeOp.Ref,
		Type:   orig,
		Range1: orig,
	}
}

func floatTypeOps(orig types.Type, t2 *types.Basic) OpTypes {
	ops := OpTypes{
		Ops: typeOp.Add | typeOp.Arith | typeOp.Comparable | typeOp.Float |
			typeOp.Orderable | typeOp.Ref,
		Type: orig,
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

func complexTypeOps(orig types.Type, t2 *types.Basic) OpTypes {
	ops := OpTypes{
		Ops: typeOp.Add | typeOp.Arith | typeOp.Comparable |
			typeOp.Complex | typeOp.Ref,
		Type: orig,
	}
	switch t2.Kind() {
	case types.Complex64:
		ops.Float = types.Typ[types.Float32]
	case types.Complex128:
		ops.Float = types.Typ[types.Float64]
	default:
		ops.Float = types.Typ[types.UntypedFloat]
	}
	return ops
}

func stringTypeOps(orig types.Type) OpTypes {
	return OpTypes{
		Ops: typeOp.Add | typeOp.GetIndex | typeOp.Len | typeOp.Comparable |
			typeOp.Orderable | typeOp.Range0 | typeOp.Range1 | typeOp.Range2 |
			typeOp.Ref | typeOp.Slice,
		Type:   orig,
		Key:    types.Typ[types.UntypedInt],
		Elem:   types.Typ[types.Byte],
		Slice:  orig,
		Range1: types.Typ[types.Int],
		Range2: RuneType(),
	}
}

func unsafePointerTypeOps(orig types.Type) OpTypes {
	return OpTypes{
		Ops:  typeOp.Comparable | typeOp.IsNil | typeOp.Ref,
		Type: orig,
	}
}

func untypedNilTypeOps(orig types.Type) OpTypes {
	return OpTypes{
		Ops:  typeOp.IsNil | typeOp.Ref,
		Type: orig,
	}
}

func chanTypeOps(orig types.Type, t2 *types.Chan) OpTypes {
	ops := OpTypes{
		Ops:  typeOp.Len | typeOp.Cap | typeOp.Make1 | typeOp.Make2 | typeOp.Ref | typeOp.IsNil | typeOp.Comparable,
		Type: orig,
		Elem: t2.Elem(),
	}
	switch t2.Dir() {
	case types.SendOnly:
		ops.Ops |= typeOp.Send
	case types.RecvOnly:
		ops.Ops |= typeOp.Range0 | typeOp.Range1 | typeOp.Recv
		ops.Range1 = t2.Elem()
	default:
		ops.Ops |= typeOp.Range0 | typeOp.Range1 | typeOp.Recv | typeOp.Send
		ops.Range1 = t2.Elem()
	}
	return ops
}

func mapTypeOps(orig types.Type, t2 *types.Map) OpTypes {
	return OpTypes{
		Ops: typeOp.GetIndex | typeOp.IsNil | typeOp.Len | typeOp.Make1 | typeOp.Make2 |
			typeOp.Range0 | typeOp.Range1 | typeOp.Range2 | typeOp.Ref | typeOp.SetIndex,
		Type:   orig,
		Key:    t2.Key(),
		Elem:   t2.Elem(),
		Range1: t2.Key(),
		Range2: t2.Elem(),
	}
}

func pointerTypeOps(orig types.Type, t2 *types.Pointer) OpTypes {
	var ops OpTypes
	if at, ok := t2.Elem().Underlying().(*types.Array); ok {
		ops = arrayTypeOps(orig, at)
	}
	ops.Type = orig
	ops.Ops |= typeOp.Comparable | typeOp.IsNil | typeOp.Ref
	return ops
}

func signatureTypeOps(orig types.Type, t2 *types.Signature) OpTypes {
	// Check for iter.Seq and iter.Seq2 (receivers are optional).
	if t2.Params().Len() == 1 && t2.Results().Len() == 0 {
		par := t2.Params().At(0).Type()
		if fn, ok := par.(*types.Signature); ok &&
			fn.Results().Len() == 1 && IsBool(fn.Results().At(0).Type()) {
			// At this point we know the signature looks like `func(func(<params>)bool)`.
			switch fn.Params().Len() {
			case 0: // Found an `func(func() bool)`
				return OpTypes{
					Ops:  typeOp.IsNil | typeOp.Range0 | typeOp.Ref,
					Type: orig,
				}
			case 1: // Found an `iter.Seq[K]`
				return OpTypes{
					Ops:    typeOp.IsNil | typeOp.Range0 | typeOp.Range1 | typeOp.Ref,
					Type:   orig,
					Range1: fn.Params().At(0).Type(),
				}
			case 2: // Found an `iter.Seq[K, V]`
				return OpTypes{
					Ops:    typeOp.IsNil | typeOp.Range0 | typeOp.Range1 | typeOp.Range2 | typeOp.Ref,
					Type:   orig,
					Range1: fn.Params().At(0).Type(),
					Range2: fn.Params().At(1).Type(),
				}
			}
		}
	}
	return OpTypes{
		Ops:  typeOp.IsNil | typeOp.Ref,
		Type: orig,
	}
}

func sliceTypeOps(orig types.Type, t2 *types.Slice) OpTypes {
	return OpTypes{
		Ops: typeOp.Cap | typeOp.GetIndex | typeOp.IsNil | typeOp.Len |
			typeOp.Make2 | typeOp.Make3 | typeOp.Range0 | typeOp.Range1 | typeOp.Range2 | typeOp.Ref |
			typeOp.RefIndex | typeOp.SetIndex | typeOp.Slice | typeOp.Slice3,
		Type:   orig,
		Key:    types.Typ[types.UntypedInt],
		Elem:   t2.Elem(),
		Slice:  orig,
		Range1: types.Typ[types.Int],
		Range2: t2.Elem(),
	}
}

func structTypeOps(orig types.Type, t2 *types.Struct) OpTypes {
	ops := OpTypes{
		Ops:  typeOp.Ref,
		Type: orig,
	}
	if types.Comparable(t2) {
		ops.Ops |= typeOp.Comparable
	}
	return ops
}

func interfaceTypeOps(orig types.Type, t2 *types.Interface) OpTypes {
	if t2.IsMethodSet() {
		if _, ok := orig.(*types.TypeParam); !ok {
			return OpTypes{
				Ops:  typeOp.IsNil | typeOp.Comparable | typeOp.Ref,
				Type: orig,
			}
		}
	}
	return unionTermsOps(orig, t2)
}

func unionTypeOps(orig types.Type, t2 *types.Union) OpTypes {
	return unionTermsOps(orig, t2)
}

func termOps(t *types.Term) OpTypes {
	if t.Tilde() {
		return getOps(types.NewUnion([]*types.Term{t}), t.Type())
	}
	return getOps(t.Type(), t.Type())
}

func unionTermsOps(orig types.Type, t2 types.Type) OpTypes {
	terms, err := typeparams.NormalTerms(t2)
	if err != nil {
		panic(faults.New(`failed to determine op type from terms`, err).
			With(`union`, t2).
			WithTypeOf(`type`, t2))
	}

	switch len(terms) {
	case 0:
		ops := OpTypes{
			Ops:  typeOp.Ref,
			Type: orig,
		}
		if it, ok := t2.(*types.Interface); ok && it.IsComparable() {
			ops.Ops |= typeOp.Comparable
		}
		return ops
	case 1:
		return getOps(orig, terms[0].Type())
	}

	// Range, make, and (other than for []byte|string) slicing are only
	// permitted when every term has the same underlying type, so in that
	// case the union behaves exactly like that single type.
	if under := sharedUnderlying(terms); under != nil {
		return getOps(orig, under)
	}

	opTerms := make([]OpTypes, len(terms))
	ops := typeOp.None
	for i, t := range terms {
		u := termOps(t)
		opTerms[i] = u
		if i == 0 {
			ops = u.Ops
		} else {
			ops &= u.Ops
		}
	}

	// Go never allows these on a type parameter whose terms have different underlying types.
	ops &^= typeOp.Range0 | typeOp.Range1 | typeOp.Range2 |
		typeOp.Make1 | typeOp.Make2 | typeOp.Make3 |
		// real, imag, and complex currently do not work in Go type parameters.
		typeOp.Complex | typeOp.Float
	if ops == typeOp.None {
		return blankOps(orig)
	}

	pick := func(adjOps typeOp.Op, getType func(OpTypes) types.Type) types.Type {
		if !ops.Any(adjOps) {
			return nil
		}
		if t := commonType(opTerms, getType); t != nil {
			return t
		}
		ops &^= adjOps
		return nil
	}

	// Key must be picked before Elem since a mismatched key
	// (e.g. a map mixed with a slice, `untyped int` vs `int`)
	// removes the index, range1, and range2 ops,
	// so there is no need to set the Elem.
	ot := OpTypes{}
	ot.Key = pick(typeOp.IndexOps, func(op OpTypes) types.Type { return op.Key })
	ot.Elem = pick(typeOp.IndexOps, func(op OpTypes) types.Type { return op.Elem })
	ot.Slice = pick(typeOp.SliceOps, func(op OpTypes) types.Type { return op.Slice })

	// set the result ops to the ops value after adj has modified it.
	ot.Ops = ops
	ot.Type = orig
	return ot
}

func sharedUnderlying(t2 []*types.Term) types.Type {
	under := t2[0].Type().Underlying()
	for _, t := range t2[1:] {
		if !types.Identical(under, t.Type().Underlying()) {
			return nil
		}
	}
	return under
}

func commonType(opTerms []OpTypes, getType func(OpTypes) types.Type) types.Type {
	first := getType(opTerms[0])
	for _, term := range opTerms[1:] {
		if t := getType(term); t == nil || !types.Identical(first, t) {
			return nil
		}
	}
	return first
}
