package augmenter

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/project/fileMod"
)

type augAdd struct {
	addImports []*ast.ImportSpec
	addDecls   []ast.Decl
	addFields  []*ast.TypeSpec
}

func (a *augAdd) Modify(fm *fileMod.FileMod) error {
	// TODO: Implement
	return nil
}

func (a *augAdd) Finished() error {
	// TODO: Implement
	return nil
}
