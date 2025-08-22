package file

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/internal/iterator"
)

type DeclSpecIteratorValue struct {
	File       *File
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
	return ds.File.FileSet.Position(pos)
}

func JoinComments(cgs ...*ast.CommentGroup) []*ast.Comment {
	result := []*ast.Comment{}
	for _, cg := range cgs {
		if cg != nil {
			result = append(result, cg.List...)
		}
	}
	return result
}

func (ds *DeclSpecIteratorValue) Comments() []*ast.Comment {
	switch {
	case ds.FuncDecl != nil:
		return JoinComments(ds.FuncDecl.Doc)
	case ds.TypeSpec != nil:
		return JoinComments(ds.GenDecl.Doc, ds.TypeSpec.Doc, ds.TypeSpec.Comment)
	case ds.ValueSpec != nil:
		return JoinComments(ds.GenDecl.Doc, ds.ValueSpec.Doc, ds.ValueSpec.Comment)
	case ds.ImportSpec != nil:
		return JoinComments(ds.GenDecl.Doc, ds.ImportSpec.Doc, ds.ImportSpec.Comment)
	}
	return nil
}

func newDeclSpecIteratorFunc(f *File, i int, d *ast.FuncDecl) *DeclSpecIteratorValue {
	return &DeclSpecIteratorValue{File: f, DeclIndex: i, FuncDecl: d, SpecIndex: -1, Node: d}
}

func newDeclSpecIteratorType(f *File, i, j int, d *ast.GenDecl, s *ast.TypeSpec) *DeclSpecIteratorValue {
	return &DeclSpecIteratorValue{File: f, DeclIndex: i, SpecIndex: j, GenDecl: d, TypeSpec: s, Node: s}
}

func newDeclSpecIteratorValue(f *File, i, j int, d *ast.GenDecl, s *ast.ValueSpec) *DeclSpecIteratorValue {
	return &DeclSpecIteratorValue{File: f, DeclIndex: i, SpecIndex: j, GenDecl: d, ValueSpec: s, Node: s}
}

func newDeclSpecIteratorImport(f *File, i, j int, d *ast.GenDecl, s *ast.ImportSpec) *DeclSpecIteratorValue {
	return &DeclSpecIteratorValue{File: f, DeclIndex: i, SpecIndex: j, GenDecl: d, ImportSpec: s, Node: s}
}

// DeclSpecs iterates through all the declarations and
// specifications in the file.
func (f *File) DeclSpecs() iterator.Iterator[*DeclSpecIteratorValue] {
	return func(yield func(v *DeclSpecIteratorValue) bool) {
		for i, decl := range f.File.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if !yield(newDeclSpecIteratorFunc(f, i, d)) {
					return
				}
			case *ast.GenDecl:
				for j, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ImportSpec:
						if !yield(newDeclSpecIteratorImport(f, i, j, d, s)) {
							return
						}
					case *ast.TypeSpec:
						if !yield(newDeclSpecIteratorType(f, i, j, d, s)) {
							return
						}
					case *ast.ValueSpec:
						if !yield(newDeclSpecIteratorValue(f, i, j, d, s)) {
							return
						}
					}
				}
			}
		}
	}
}
