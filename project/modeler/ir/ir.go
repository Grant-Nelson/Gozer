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

type Node interface {
	Pos() token.Pos
}

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

func astPos(n ast.Node) token.Pos {
	if n == nil {
		return token.NoPos
	}
	return n.Pos()
}
