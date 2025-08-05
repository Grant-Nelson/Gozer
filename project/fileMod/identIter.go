package fileMod

import (
	"go/ast"
	"iter"
)

type IdentIter struct {
	fm         *FileMod
	DeclIndex  int
	SpecIndex  int
	ValueIndex int
	Node       ast.Node
	Name       string
}

func (fm *FileMod) IdentIter() iter.Seq[IdentIter] {
	return func(yield func(v IdentIter) bool) {
		for i, decl := range fm.Decls {
			fm.identIterDecl(i, decl, yield)
		}
	}
}

func (fm *FileMod) identIterDecl(i int, decl ast.Decl, yield func(v IdentIter) bool) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		if !yield(IdentIter{
			fm:         fm,
			DeclIndex:  i,
			SpecIndex:  -1,
			ValueIndex: -1,
			Node:       d,
			Name:       d.Name.Name,
		}) {
			return
		}
	case *ast.GenDecl:
		for j, spec := range d.Specs {
			fm.identIterSpec(i, j, spec, yield)
		}
	}
}

func (fm *FileMod) identIterSpec(i, j int, spec ast.Spec, yield func(v IdentIter) bool) {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		if !yield(IdentIter{
			fm:         fm,
			DeclIndex:  i,
			SpecIndex:  j,
			ValueIndex: -1,
			Node:       s,
			Name:       s.Name.Name,
		}) {
			return
		}
	case *ast.ValueSpec:
		for k, ident := range s.Names {
			if ident == nil || ident.Name == `_` {
				continue
			}
			if !yield(IdentIter{
				fm:         fm,
				DeclIndex:  i,
				SpecIndex:  j,
				ValueIndex: k,
				Node:       s,
				Name:       ident.Name,
			}) {
				return
			}
		}
	}
}
