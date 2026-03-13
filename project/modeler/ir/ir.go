package ir

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
)

const (
	directiveGroup      = `gozer`
	directiveAtomicFunc = `atomic`
)

func csvString[E any, S []E](s S) string {
	elems := make([]string, len(s))
	for i, elem := range s {
		elems[i] = toString(elem)
	}
	return strings.Join(elems, `, `)
}

func linesString[E any, S []E](s S) string {
	const indent = `  `
	elems := make([]string, len(s))
	for i, elem := range s {
		eStr := indent + toString(elem)
		elems[i] = strings.ReplaceAll(eStr, "\n", "\n"+indent)
	}
	return strings.Join(elems, "\n")
}

func toString(t any) string {
	switch t := t.(type) {
	case nil:
		return `<nil>`
	case ast.Node:
		return nodeString(t)
	case interface{ String() string }:
		return t.String()
	default:
		return fmt.Sprintf(`%v`, t)
	}
}

func nodeString(n ast.Node) string {
	buf := &bytes.Buffer{}
	fSet := token.NewFileSet()
	if err := format.Node(buf, fSet, n); err != nil {
		panic(err)
	}
	return buf.String()
}
