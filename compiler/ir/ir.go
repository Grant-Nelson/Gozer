package ir

import (
	"fmt"
	"go/constant"
	"regexp"
	"strings"
	"sync"

	"github.com/Grant-Nelson/Gozer/avail/wad"
)

// Is determines if the given value is of the given type.
func Is[T any](v any) bool {
	_, ok := v.(T)
	return ok
}

func init() {
	wad.Register[*BasicLit]()
	//wad.Register[*BinaryExpr]()
	//wad.Register[*Block]()
	//wad.Register[*BlockRef]()
	//wad.Register[*BranchStmt]()
	//wad.Register[*BuiltinRef]()
	//wad.Register[*CallExpr]()
	//wad.Register[*CaseClause]()
	//wad.Register[*CommClause]()
	//wad.Register[*ConstDef]()
	//wad.Register[*ConstRef]()
	//wad.Register[*DeferStmt]()
	//wad.Register[*Directive]()
	//wad.Register[*ForStmt]()
	//wad.Register[*Func]()
	//wad.Register[*FuncDef]()
	//wad.Register[*FuncLit]()
	//wad.Register[*FuncRef]()
	//wad.Register[*GoStmt]()
	//wad.Register[*GotoBlockStmt]()
	//wad.Register[*IfStmt]()
	//wad.Register[*ImportDef]()
	//wad.Register[*ImportRef]()
	//wad.Register[*IndexExpr]()
	//wad.Register[*IndexListExpr]()
	//wad.Register[*InjectStmt]()

	// TODO: Continue adding
}

const (
	indentStr = `  `
	nlStr     = "\n"
)

func csvString[E any, S ~[]E](s S) string {
	const sep = `, `
	elems := make([]string, len(s))
	for i, elem := range s {
		elems[i] = toString(elem)
	}
	return strings.Join(elems, sep)
}

func linesString[E any, S ~[]E](s S) string {
	elems := make([]string, len(s))
	for i, elem := range s {
		elems[i] = indentStr + indentInner(toString(elem))
	}
	return strings.Join(elems, nlStr)
}

func indentInner(s string) string {
	return strings.ReplaceAll(s, nlStr, nlStr+indentStr)
}

func bodyString[E any, S ~[]E](body S) string {
	if len(body) <= 0 {
		return ` {}`
	}
	return ` {` + nlStr + linesString(body) + nlStr + `}`
}

func groupString[E any, S ~[]E](name string, group S) string {
	if len(group) <= 0 {
		return ``
	}
	return indentStr + name + `:` + nlStr + indentStr +
		indentInner(linesString(group)) + nlStr
}

var wordReg = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[a-zA-Z0-9_$.]+$`)
})

func paren(a Expr) string {
	v := `_`
	if a != nil {
		v = toString(a)
	}
	if wordReg().MatchString(v) {
		return v
	}
	return `(` + v + `)`
}

func emptyZeroOrString[T comparable](t T) string {
	var zero T
	if t == zero {
		return ``
	}
	return toString(t)
}

func toString(t any) string {
	switch t := t.(type) {
	case nil:
		return `<nil>`
	case string:
		return t
	case constant.Value:
		return t.ExactString()
	case interface{ String() string }:
		return t.String()
	default:
		return fmt.Sprintf(`%v`, t)
	}
}
