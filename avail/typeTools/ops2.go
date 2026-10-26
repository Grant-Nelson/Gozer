package typeTools

import (
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/typeTools/typeOp"
)

func arrayTypeOps2(t2 *types.Array) OpTypes {
	return OpTypes{
		Ops: typeOp.Cap | typeOp.Comparable | typeOp.GetIndex | typeOp.Len |
			typeOp.Range0 | typeOp.Range1 | typeOp.Range2 | typeOp.Ref |
			typeOp.RefIndex | typeOp.SetIndex | typeOp.Slice | typeOp.Slice3,
		Key:    types.Typ[types.UntypedInt],
		Elem:   t2.Elem(),
		Slice:  types.NewSlice(t2.Elem()),
		Range1: types.Typ[types.Int],
		Range2: t2.Elem(),
	}
}

func integerTypeOps2(orig types.Type) OpTypes {
	return OpTypes{
		Ops: typeOp.Add | typeOp.Arith | typeOp.Bitwise | typeOp.Comparable |
			typeOp.Mod | typeOp.Orderable | typeOp.Range0 | typeOp.Range1 | typeOp.Ref,
		Range1: orig,
	}
}

func unionTermsOps2(orig types.Type, t2 []*types.Term) OpTypes {
	if len(t2) == 0 {
		return OpTypes{}
	}

	// Range, make, and (other than for []byte|string) slicing are only
	// permitted when every term has the same underlying type, so in that
	// case the union behaves exactly like that single type.
	if under := sharedUnderlying(t2); under != nil {
		return getOps(orig, under)
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

	// Go never allows these on a type parameter whose terms have different
	// underlying types, and never allows real, imag, or complex on any
	// type parameter.
	ops &^= typeOp.Range0 | typeOp.Range1 | typeOp.Range2 | typeOp.Make |
		typeOp.Complex | typeOp.Float
	if ops == typeOp.None {
		return OpTypes{}
	}

	pick := func(adjOps typeOp.Op, getType func(OpTypes) types.Type) types.Type {
		if !ops.Any(adjOps) {
			return nil
		}
		if t := commonType(terms, getType); t != nil {
			return t
		}
		ops &^= adjOps
		return nil
	}

	// Key must be picked before Elem since a mismatched key (e.g. a map mixed
	// with a slice, `untyped int` vs `int`) removes the index ops.
	indexOps := typeOp.GetIndex | typeOp.RefIndex | typeOp.SetIndex
	ot := OpTypes{}
	ot.Key = pick(indexOps, func(op OpTypes) types.Type { return op.Key })
	ot.Elem = pick(indexOps, func(op OpTypes) types.Type { return op.Elem })
	ot.Slice = pick(typeOp.Slice|typeOp.Slice3, func(op OpTypes) types.Type { return op.Slice })
	ot.Ops = ops
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
