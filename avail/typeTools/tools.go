package typeTools

import (
	"go/token"
	"go/types"
	"sync"
)

func BasicKind(t types.Type) types.BasicKind {
	if b, ok := t.Underlying().(*types.Basic); ok {
		return b.Kind()
	}
	return types.Invalid
}

func BasicInfo(t types.Type) types.BasicInfo {
	if b, ok := t.Underlying().(*types.Basic); ok {
		return b.Info()
	}
	return 0
}

func Deref(t types.Type) types.Type {
	if b, ok := t.Underlying().(*types.Pointer); ok {
		return b.Elem()
	}
	return nil
}

func namedType[T types.Type](exp string) func() T {
	return func() T {
		tv, err := types.Eval(token.NewFileSet(), nil, token.NoPos, exp)
		if err != nil {
			panic(err)
		}
		return tv.Type.(T)
	}
}

var (
	runeType      = sync.OnceValue(namedType[*types.Basic](`rune`))
	byteType      = sync.OnceValue(namedType[*types.Basic](`byte`))
	anyType       = sync.OnceValue(namedType[*types.Alias](`any`))
	byteSliceType = sync.OnceValue(namedType[*types.Slice](`[]byte`))
)

func RuneType() *types.Basic { return runeType() }

func ByteType() *types.Basic { return byteType() }

func AnyType() *types.Alias { return anyType() }

func ByteSliceType() *types.Slice { return byteSliceType() }
