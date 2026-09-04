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

func namedBasicAlias(exp string) func() *types.Basic {
	return func() *types.Basic {
		tv, err := types.Eval(token.NewFileSet(), nil, token.NoPos, exp)
		if err != nil {
			panic(err)
		}
		return tv.Type.(*types.Basic)
	}
}

var runeType = sync.OnceValue(namedBasicAlias(`rune`))

func RuneType() *types.Basic { return runeType() }

var byteType = sync.OnceValue(namedBasicAlias(`byte`))

func ByteType() *types.Basic { return byteType() }
