package typeTools

import (
	"go/types"
)

func Unalias(t types.Type) types.Type {
	return types.Unalias(t)
}

func BasicKind(t types.Type) types.BasicKind {
	if b, ok := t.(*types.Basic); ok {
		return b.Kind()
	}
	return types.Invalid
}

func BasicInfo(t types.Type) types.BasicInfo {
	if b, ok := t.(*types.Basic); ok {
		return b.Info()
	}
	return 0
}

func Deref(t types.Type) types.Type {
	if b, ok := t.(*types.Pointer); ok {
		return b.Elem()
	}
	return nil
}

func Key(t types.Type) types.Type {
	switch t := t.(type) {
	case *types.Map:
		return t.Key()
	case *types.Slice:
		return types.Typ[types.UntypedInt]
	case *types.Array:
		return types.Typ[types.UntypedInt]
	case *types.Basic:
		if t.Kind() == types.String {
			return types.Typ[types.UntypedInt]
		}
	case *types.Pointer:
		if _, ok := t.Elem().(*types.Array); ok {
			return types.Typ[types.UntypedInt]
		}
	}
	return nil
}

func Elem(t types.Type) types.Type {
	switch t := t.(type) {
	case *types.Map:
		return t.Elem()
	case *types.Slice:
		return t.Elem()
	case *types.Array:
		return t.Elem()
	case *types.Basic:
		if t.Kind() == types.String {
			return types.Typ[types.Byte]
		}
	case *types.Pointer:
		if s, ok := t.Elem().(*types.Array); ok {
			return s.Elem()
		}
	}
	return nil
}

/*
func Ops(t types.Type) typeOp.Op {
	switch t := t.(type) {
	case *types.Basic:

	case *types.Slice:
	case *types.Array:
	case *types.Map:
	case *types.Struct:
	case *types.Pointer:
	}
	return typeOp.None
}
*/
