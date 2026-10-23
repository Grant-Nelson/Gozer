package ir

import (
	"go/token"
	"go/types"
)

// ConstDef is the declaration/definition for a single constant.
type ConstDef struct {
	// Comment for this constant.
	Comment string

	// Directives attached to this constant.
	Directives []Directive

	// ConstObj is the object for this constant.
	ConstObj *types.Const
}

var (
	_ Def    = (*ConstDef)(nil)
	_ Expr   = (*ConstDef)(nil)
	_ Stmt   = (*ConstDef)(nil)
	_ Parent = (*ConstDef)(nil)
)

func (*ConstDef) DefNode()  {}
func (*ConstDef) ExprNode() {}
func (*ConstDef) StmtNode() {}
func (*ConstDef) RefNode()  {}

func (n *ConstDef) Pos() token.Pos       { return n.ConstObj.Pos() }
func (n *ConstDef) Type() types.Type     { return n.ConstObj.Type() }
func (n *ConstDef) Object() types.Object { return n.ConstObj }

func (n *ConstDef) String() string {
	return toString(n.ConstObj) + ` = ` + toString(n.ConstObj.Val())
}

func (n *ConstDef) ChildCount() int { return len(n.Directives) }
func (n *ConstDef) Children(yield func(Node) bool) {
	_ = YieldSlice(n.Directives, yield)
}
