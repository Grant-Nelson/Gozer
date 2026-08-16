package ir

import (
	"go/token"
	"go/types"
)

// ImportDecl is the declaration for a single import.
type ImportDecl struct {
	PkgObj *types.PkgName
}

var (
	_ Stmt = (*ImportDecl)(nil)
	_ Decl = (*ImportDecl)(nil)
)

func (*ImportDecl) StmtNode() {}
func (*ImportDecl) DeclNode() {}

func (n *ImportDecl) Pos() token.Pos       { return n.PkgObj.Pos() }
func (n *ImportDecl) Type() types.Type     { return n.PkgObj.Type() }
func (n *ImportDecl) Object() types.Object { return n.PkgObj }
func (n *ImportDecl) String() string       { return `const ` + n.PkgObj.String() }
