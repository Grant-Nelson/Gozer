package ir

import (
	"fmt"
	"strings"
)

func csvString[E any, S []E](s S) string {
	const sep = `, `
	elems := make([]string, len(s))
	for i, elem := range s {
		elems[i] = toString(elem)
	}
	return strings.Join(elems, sep)
}

func linesString[E any, S []E](s S) string {
	const indent = `  `
	const nl = "\n"
	elems := make([]string, len(s))
	for i, elem := range s {
		eStr := indent + toString(elem)
		elems[i] = strings.ReplaceAll(eStr, nl, nl+indent)
	}
	return strings.Join(elems, nl)
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
