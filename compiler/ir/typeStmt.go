package ir

import (
	"go/token"
	"go/types"
)

type TypeStmt struct {
	TypePos      token.Pos
	TypeAndValue *types.TypeAndValue
}

var _ Stmt = (*TypeStmt)(nil)

func (n *TypeStmt) Pos() token.Pos { return n.TypePos }

func (n *TypeStmt) Type() types.Type { return n.TypeAndValue.Type }

func (n *TypeStmt) StmtNode() {}

func (n *TypeStmt) String() string {
	if n.TypeAndValue.Value != nil {
		return n.TypeAndValue.Value.String()
	}
	return n.TypeAndValue.Type.String()
}
