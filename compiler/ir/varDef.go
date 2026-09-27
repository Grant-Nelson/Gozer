package ir

import (
	"go/token"
	"go/types"
)

// VarDef is the declaration/definition for a single variable.
type VarDef struct {
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
	_ Expr   = (*VarDef)(nil)
	_ Def    = (*VarDef)(nil)
	_ Expr   = (*VarDef)(nil)
	_ Stmt   = (*VarDef)(nil)
	_ Parent = (*VarDef)(nil)
)

func (*VarDef) ExprNode() {}
func (*VarDef) DefNode()  {}
func (*VarDef) StmtNode() {}

func (n *VarDef) Pos() token.Pos       { return n.VarObj.Pos() }
func (n *VarDef) Type() types.Type     { return n.VarObj.Type() }
func (n *VarDef) Object() types.Object { return n.VarObj }

func (n *VarDef) String() string {
	if n.Value == nil {
		return `decl ` + toString(n.VarObj)
	}
	return `decl ` + toString(n.VarObj) + ` = ` + toString(n.Value)
}

func (n *VarDef) ChildCount() int { return 1 + len(n.Directives) }

func (n *VarDef) Children(yield func(Node) bool) {
	_ = YieldSlice(n.Directives, yield) &&
		YieldNode(n.Value, yield)
}
