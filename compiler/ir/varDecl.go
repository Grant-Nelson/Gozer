package ir

import (
	"go/token"
	"go/types"
)

// VarDecl is the declaration for a single variable.
type VarDecl struct {
	// Comment for this variable.
	Comment string

	// Directives attached to this variable.
	Directives []Directive

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

func (n *VarDecl) Pos() token.Pos       { return n.VarObj.Pos() }
func (n *VarDecl) Type() types.Type     { return n.VarObj.Type() }
func (n *VarDecl) Object() types.Object { return n.VarObj }

func (n *VarDecl) String() string {
	if n.Value == nil {
		return toString(n.VarObj)
	}
	return toString(n.VarObj) + ` = ` + toString(n.Value)
}

func (n *VarDecl) ChildCount() int { return 1 + len(n.Directives) }

func (n *VarDecl) Children(yield func(Node) bool) {
	_ = YieldSlice(n.Directives, yield) &&
		YieldNode(n.Value, yield)
}
