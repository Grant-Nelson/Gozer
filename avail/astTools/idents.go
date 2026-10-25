package astTools

import (
	"go/ast"
	"go/token"

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

// Idents will return all identifier values at the package-level, including method names
// (names for functions with receivers) and `init` functions, for the given file.
// Note that the *IdentIteratorValue in the iterator may be reused
// and mutated between iteration steps.
func Idents(fSet *token.FileSet, f *ast.File) iterator.Iterator[*IdentIteratorValue] {
	return func(yield func(v *IdentIteratorValue) bool) {
		for ds := range DeclSpecs(fSet, f) {
			switch {
			case ds.ImportSpec != nil:
				if !yield(newIdentIteratorValue(ds, -1, ds.ImportSpec.Name.Name)) {
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
