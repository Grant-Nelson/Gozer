package ir

import (
	"go/token"
	"go/types"
)

type TypeDecl struct {
	// Comment for this type declaration.
	Comment string
	TypeObj *types.TypeName
}

var (
	_ Stmt = (*TypeDecl)(nil)
	_ Decl = (*TypeDecl)(nil)
)

func (*TypeDecl) StmtNode() {}
func (*TypeDecl) DeclNode() {}

func (n *TypeDecl) Pos() token.Pos       { return n.TypeObj.Pos() }
func (n *TypeDecl) Type() types.Type     { return n.TypeObj.Type() }
func (n *TypeDecl) Object() types.Object { return n.TypeObj }
func (n *TypeDecl) String() string       { return toString(n.TypeObj) }
