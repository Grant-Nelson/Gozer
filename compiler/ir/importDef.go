package ir

import (
	"go/token"
	"go/types"
)

// ImportDef is the declaration for a single import.
type ImportDef struct {
	PkgObj *types.PkgName
}

var (
	_ Def  = (*ImportDef)(nil)
	_ Stmt = (*ImportDef)(nil)
)

func (*ImportDef) DefNode()  {}
func (*ImportDef) StmtNode() {}

func (n *ImportDef) Pos() token.Pos       { return n.PkgObj.Pos() }
func (n *ImportDef) Type() types.Type     { return n.PkgObj.Type() }
func (n *ImportDef) Object() types.Object { return n.PkgObj }
func (n *ImportDef) String() string       { return `import ` + toString(n.PkgObj) }
