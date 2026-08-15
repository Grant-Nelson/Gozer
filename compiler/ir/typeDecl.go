package ir

import (
	"go/token"
	"go/types"
)

type TypeDecl struct {
	Name string

	NamePos token.Pos

	TypeAndValue *types.TypeAndValue
}

var (
	_ Stmt = (*TypeDecl)(nil)
	_ Decl = (*TypeDecl)(nil)
)

func (*TypeDecl) StmtNode() {}
func (*TypeDecl) DeclNode() {}

func (n *TypeDecl) Pos() token.Pos   { return n.NamePos }
func (n *TypeDecl) Type() types.Type { return n.TypeAndValue.Type }

func (n *TypeDecl) String() string {
	if n.TypeAndValue.Value != nil {
		return n.Name + ` ` + n.TypeAndValue.Value.String()
	}
	return n.Name + ` ` + n.TypeAndValue.Type.String()
}
