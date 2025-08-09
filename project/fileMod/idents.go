package fileMod

import (
	"github.com/Grant-Nelson/Gozer/internal/iterator"
)

type Ident struct {
	*DeclSpec
	ValueIndex int
	Name       string
}

func newIdent(ds *DeclSpec, k int, name string) *Ident {
	return &Ident{DeclSpec: ds, ValueIndex: k, Name: name}
}

func (fm *FileMod) Idents() iterator.Iterator[*Ident] {
	return func(yield func(v *Ident) bool) {
		for ds := range fm.DeclSpecs() {
			switch {
			case ds.FuncDecl != nil:
				if !yield(newIdent(ds, -1, ds.FuncDecl.Name.Name)) {
					return
				}
			case ds.TypeSpec != nil:
				if !yield(newIdent(ds, -1, ds.TypeSpec.Name.Name)) {
					return
				}
			case ds.ValueSpec != nil:
				for k, ident := range ds.ValueSpec.Names {
					if ident == nil || ident.Name == `_` {
						continue
					}
					ds.Node = ident
					if !yield(newIdent(ds, k, ident.Name)) {
						return
					}
				}
			}
		}
	}
}
