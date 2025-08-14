package astMod

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/internal/iterator"
)

type DeclSpecIteratorValue struct {
	FileMod    *FileMod
	DeclIndex  int
	SpecIndex  int
	FuncDecl   *ast.FuncDecl
	GenDecl    *ast.GenDecl
	TypeSpec   *ast.TypeSpec
	ValueSpec  *ast.ValueSpec
	ImportSpec *ast.ImportSpec
	Node       ast.Node
}

func (ds *DeclSpecIteratorValue) Start() token.Position {
	return ds.Position(ds.Node.Pos())
}

func (ds *DeclSpecIteratorValue) End() token.Position {
	return ds.Position(ds.Node.End())
}

func (ds *DeclSpecIteratorValue) Position(pos token.Pos) token.Position {
	return ds.FileMod.Package().FileSet().Position(pos)
}

func newDeclSpecIteratorFunc(fm *FileMod, i int, d *ast.FuncDecl) *DeclSpecIteratorValue {
	return &DeclSpecIteratorValue{FileMod: fm, DeclIndex: i, FuncDecl: d, SpecIndex: -1, Node: d}
}

func newDeclSpecIteratorType(fm *FileMod, i, j int, d *ast.GenDecl, s *ast.TypeSpec) *DeclSpecIteratorValue {
	return &DeclSpecIteratorValue{FileMod: fm, DeclIndex: i, SpecIndex: j, GenDecl: d, TypeSpec: s, Node: s}
}

func newDeclSpecIteratorValue(fm *FileMod, i, j int, d *ast.GenDecl, s *ast.ValueSpec) *DeclSpecIteratorValue {
	return &DeclSpecIteratorValue{FileMod: fm, DeclIndex: i, SpecIndex: j, GenDecl: d, ValueSpec: s, Node: s}
}

func newDeclSpecIteratorImport(fm *FileMod, i, j int, d *ast.GenDecl, s *ast.ImportSpec) *DeclSpecIteratorValue {
	return &DeclSpecIteratorValue{FileMod: fm, DeclIndex: i, SpecIndex: j, GenDecl: d, ImportSpec: s, Node: s}
}

func (fm *FileMod) DeclSpecs() iterator.Iterator[*DeclSpecIteratorValue] {
	return func(yield func(v *DeclSpecIteratorValue) bool) {
		for i, decl := range fm.file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if !yield(newDeclSpecIteratorFunc(fm, i, d)) {
					return
				}
			case *ast.GenDecl:
				for j, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ImportSpec:
						if !yield(newDeclSpecIteratorImport(fm, i, j, d, s)) {
							return
						}
					case *ast.TypeSpec:
						if !yield(newDeclSpecIteratorType(fm, i, j, d, s)) {
							return
						}
					case *ast.ValueSpec:
						if !yield(newDeclSpecIteratorValue(fm, i, j, d, s)) {
							return
						}
					}
				}
			}
		}
	}
}
