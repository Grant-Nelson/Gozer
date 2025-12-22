package artifacts

import (
	"github.com/Grant-Nelson/Gozer/avail/iterator"
)

type IdentIteratorValue struct {
	*DeclSpecIteratorValue
	ValueIndex int
	Ident      string
}

func newIdentIteratorValue(ds *DeclSpecIteratorValue, k int, ident string) *IdentIteratorValue {
	return &IdentIteratorValue{DeclSpecIteratorValue: ds, ValueIndex: k, Ident: ident}
}

func (f *File) Idents() iterator.Iterator[*IdentIteratorValue] {
	return func(yield func(v *IdentIteratorValue) bool) {
		for ds := range f.DeclSpecs() {
			switch {
			case ds.ImportSpec != nil:
				if !yield(newIdentIteratorValue(ds, -1, ds.ImportSpec.Path.Value)) {
					return
				}
			case ds.FuncDecl != nil:
				if !yield(newIdentIteratorValue(ds, -1, ds.FuncDecl.Name.Name)) {
					return
				}
			case ds.TypeSpec != nil:
				if !yield(newIdentIteratorValue(ds, -1, ds.TypeSpec.Name.Name)) {
					return
				}
			case ds.ValueSpec != nil:
				for k, ident := range ds.ValueSpec.Names {
					if ident == nil || ident.Name == `_` {
						continue
					}
					ds.Node = ident
					if !yield(newIdentIteratorValue(ds, k, ident.Name)) {
						return
					}
				}
			}
		}
	}
}
