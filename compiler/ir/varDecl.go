package ir

import (
	"go/token"
	"go/types"
)

// VarDecl is the declaration for a single variable.
type VarDecl struct {

	// Name is the name for the variable.
	Name string

	// NamePos is the position for the name of the variable.
	NamePos token.Pos

	// Type is the type of this variable.
	VarType types.Type

	// VarObj is the object for this variable.
	VarObj *types.Var

	// Value is the optional initial assigned of the variable or nil.
	Value Expr
}

var (
	_ Stmt   = (*VarDecl)(nil)
	_ Decl   = (*VarDecl)(nil)
	_ Parent = (*VarDecl)(nil)
)

func (*VarDecl) StmtNode() {}
func (*VarDecl) DeclNode() {}

func (n *VarDecl) Pos() token.Pos       { return n.NamePos }
func (n *VarDecl) Type() types.Type     { return n.VarType }
func (n *VarDecl) Object() types.Object { return n.VarObj }

func (n *VarDecl) String() string {
	result := `var ` + n.Name + ` ` + n.VarType.String()
	if n.Value != nil {
		result += ` = ` + n.Value.String()
	}
	return result
}

func (n *VarDecl) Children(yield func(Node) bool) { _ = yield(n.Value) }
