package ir

import (
	"go/token"
	"go/types"
)

// ConstDecl is the declaration for a single constant.
type ConstDecl struct {
	// ConstObj is the object for this constant.
	ConstObj *types.Const
}

var (
	_ Stmt = (*ConstDecl)(nil)
	_ Decl = (*ConstDecl)(nil)
)

func (*ConstDecl) StmtNode() {}
func (*ConstDecl) DeclNode() {}

func (n *ConstDecl) Pos() token.Pos       { return n.ConstObj.Pos() }
func (n *ConstDecl) Type() types.Type     { return n.ConstObj.Type() }
func (n *ConstDecl) Object() types.Object { return n.ConstObj }
func (n *ConstDecl) String() string       { return `const ` + n.ConstObj.String() }
