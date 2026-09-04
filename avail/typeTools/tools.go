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

var runeType = sync.OnceValue(func() *types.Basic {
	tv, err := types.Eval(token.NewFileSet(), nil, token.NoPos, `rune`)
	if err != nil {
		panic(err)
	}
	return tv.Type.(*types.Basic)
})

func RuneType() *types.Basic {
	return runeType()
}
