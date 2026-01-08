package artifacts

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/iterator"
)

type DeclSpecIteratorValue struct {
	FileSet    *token.FileSet
	File       *ast.File
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
	return ds.FileSet.Position(pos)
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

func newDeclSpecIteratorFunc(fSet *token.FileSet, f *ast.File, i int, d *ast.FuncDecl) *DeclSpecIteratorValue {
	return &DeclSpecIteratorValue{FileSet: fSet, File: f, DeclIndex: i, FuncDecl: d, SpecIndex: -1, Node: d}
}

func newDeclSpecIteratorType(fSet *token.FileSet, f *ast.File, i, j int, d *ast.GenDecl, s *ast.TypeSpec) *DeclSpecIteratorValue {
	return &DeclSpecIteratorValue{FileSet: fSet, File: f, DeclIndex: i, SpecIndex: j, GenDecl: d, TypeSpec: s, Node: s}
}

func newDeclSpecIteratorValue(fSet *token.FileSet, f *ast.File, i, j int, d *ast.GenDecl, s *ast.ValueSpec) *DeclSpecIteratorValue {
	return &DeclSpecIteratorValue{FileSet: fSet, File: f, DeclIndex: i, SpecIndex: j, GenDecl: d, ValueSpec: s, Node: s}
}

func newDeclSpecIteratorImport(fSet *token.FileSet, f *ast.File, i, j int, d *ast.GenDecl, s *ast.ImportSpec) *DeclSpecIteratorValue {
	return &DeclSpecIteratorValue{FileSet: fSet, File: f, DeclIndex: i, SpecIndex: j, GenDecl: d, ImportSpec: s, Node: s}
}

// DeclSpecs iterates through all the declarations and
// specifications in the file.
func DeclSpecs(fSet *token.FileSet, f *ast.File) iterator.Iterator[*DeclSpecIteratorValue] {
	return func(yield func(v *DeclSpecIteratorValue) bool) {
		for i, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if !yield(newDeclSpecIteratorFunc(fSet, f, i, d)) {
					return
				}
			case *ast.GenDecl:
				for j, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ImportSpec:
						if !yield(newDeclSpecIteratorImport(fSet, f, i, j, d, s)) {
							return
						}
					case *ast.TypeSpec:
						if !yield(newDeclSpecIteratorType(fSet, f, i, j, d, s)) {
							return
						}
					case *ast.ValueSpec:
						if !yield(newDeclSpecIteratorValue(fSet, f, i, j, d, s)) {
							return
						}
					}
				}
			}
		}
	}
}
