package ir

import (
	"fmt"
	"strings"
)

const indent = `  `

func csvString[E any, S ~[]E](s S) string {
	const sep = `, `
	elems := make([]string, len(s))
	for i, elem := range s {
		elems[i] = toString(elem)
	}
	return strings.Join(elems, sep)
}

func linesString[E any, S ~[]E](s S) string {
	const nl = "\n"
	elems := make([]string, len(s))
	for i, elem := range s {
		eStr := indent + toString(elem)
		elems[i] = strings.ReplaceAll(eStr, nl, nl+indent)
	}
	return strings.Join(elems, nl)
}

func bodyString[E any, S ~[]E](body S) string {
	count := len(body)
	if count <= 0 {
		return ` {}`
	}
	if count == 1 {
		return "\n" + indent + toString(body[0])
	}
	return " {\n" + linesString(body) + "\n}"
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
	case interface{ String() string }:
		return t.String()
	default:
		return fmt.Sprintf(`%v`, t)
	}
}
