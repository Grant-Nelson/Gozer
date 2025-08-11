package fileMod

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/internal/iterator"
)

type DeclSpec struct {
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

func (ds *DeclSpec) Start() token.Position {
	return ds.Position(ds.Node.Pos())
}

func (ds *DeclSpec) End() token.Position {
	return ds.Position(ds.Node.End())
}

func (ds *DeclSpec) Position(pos token.Pos) token.Position {
	return ds.FileMod.FileSet().Position(pos)
}

func newDeclSpecFunc(fm *FileMod, i int, d *ast.FuncDecl) *DeclSpec {
	return &DeclSpec{FileMod: fm, DeclIndex: i, FuncDecl: d, SpecIndex: -1, Node: d}
}

func newDeclSpecType(fm *FileMod, i, j int, d *ast.GenDecl, s *ast.TypeSpec) *DeclSpec {
	return &DeclSpec{FileMod: fm, DeclIndex: i, SpecIndex: j, GenDecl: d, TypeSpec: s, Node: s}
}

func newDeclSpecValue(fm *FileMod, i, j int, d *ast.GenDecl, s *ast.ValueSpec) *DeclSpec {
	return &DeclSpec{FileMod: fm, DeclIndex: i, SpecIndex: j, GenDecl: d, ValueSpec: s, Node: s}
}

func newDeclSpecImport(fm *FileMod, i, j int, d *ast.GenDecl, s *ast.ImportSpec) *DeclSpec {
	return &DeclSpec{FileMod: fm, DeclIndex: i, SpecIndex: j, GenDecl: d, ImportSpec: s, Node: s}
}

func (fm *FileMod) DeclSpecs() iterator.Iterator[*DeclSpec] {
	return func(yield func(v *DeclSpec) bool) {
		for i, decl := range fm.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if !yield(newDeclSpecFunc(fm, i, d)) {
					return
				}
			case *ast.GenDecl:
				for j, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ImportSpec:
						if !yield(newDeclSpecImport(fm, i, j, d, s)) {
							return
						}
					case *ast.TypeSpec:
						if !yield(newDeclSpecType(fm, i, j, d, s)) {
							return
						}
					case *ast.ValueSpec:
						if !yield(newDeclSpecValue(fm, i, j, d, s)) {
							return
						}
					}
				}
			}
		}
	}
}
