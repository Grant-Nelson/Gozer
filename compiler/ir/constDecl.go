package ir

import (
	"go/token"
	"go/types"
)

// ConstDecl is the declaration for a single constant.
type ConstDecl struct {

	// Name is the name for the value.
	Name string

	// NamePos is the position for the name of the value.
	NamePos token.Pos

	// TypeAndValue is the type and constant value.
	TypeAndValue types.TypeAndValue

	// ConstObj is the object for this constant.
	ConstObj *types.Const

	// Value is the optional initial assigned of the value or nil.
	Value Expr
}

var (
	_ Stmt   = (*ConstDecl)(nil)
	_ Decl   = (*ConstDecl)(nil)
	_ Parent = (*ConstDecl)(nil)
)

func (*ConstDecl) StmtNode() {}
func (*ConstDecl) DeclNode() {}

func (n *ConstDecl) Pos() token.Pos       { return n.NamePos }
func (n *ConstDecl) Type() types.Type     { return n.TypeAndValue.Type }
func (n *ConstDecl) Object() types.Object { return n.ConstObj }

func (n *ConstDecl) String() string {
	result := `const ` + n.Name + ` ` + n.TypeAndValue.Type.String()
	if n.Value != nil {
		result += ` = ` + n.Value.String()
	}
	return result
}

func (n *ConstDecl) Children(yield func(Node) bool) { _ = yield(n.Value) }
