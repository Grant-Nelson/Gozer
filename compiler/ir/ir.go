package ir

import (
	"fmt"
	"go/constant"
	"strings"
)

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
	count := len(body)
	if count <= 0 {
		return ` {}`
	}
	if count == 1 {
		return nlStr + indentStr + indentInner(toString(body[0]))
	}
	return ` {` + nlStr + linesString(body) + nlStr + `}`
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
	case constant.Value:
		return t.ExactString()
	case interface{ String() string }:
		return t.String()
	default:
		return fmt.Sprintf(`%v`, t)
	}
}
