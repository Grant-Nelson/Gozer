package augmenter

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/project/fileMod"
)

type augDel struct {
	delImport []*ast.ImportSpec
	delFunc   []*ast.FuncDecl
	delVar    []*ast.ValueSpec
	delType   []*ast.TypeSpec
	delFields []*ast.TypeSpec
}

func (a *augDel) Modify(fm *fileMod.FileMod) error {
	// TODO: Implement
	return nil
}

func (a *augDel) Finished() error {
	// TODO: Implement
	return nil
}
