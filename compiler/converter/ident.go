package converter

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/compiler/ir"
)

func (c *converter) FromIdent(e *ast.Ident) ir.Expr {
	if e == nil {
		return nil
	}
	id := &ir.Ident{
		NamePos: e.NamePos,
		Name:    e.Name,
	}
	info := c.Source.TypesInfo
	if tv, ok := info.Types[e]; ok {
		id.TypeAndValue = &tv
	}
	if in, ok := info.Instances[e]; ok {
		id.Instance = &in
	}
	if ds, ok := info.Defs[e]; ok {
		id.Def = ds
	}
	if us, ok := info.Uses[e]; ok {
		id.Use = us
	}
	return id
}
