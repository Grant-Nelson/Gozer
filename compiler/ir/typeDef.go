package ir

import (
	"go/token"
	"go/types"
)

type TypeDef struct {
	// Comment for this type declaration.
	Comment string

	// Directives attached to this type declaration.
	Directives []Directive

	TypeObj *types.TypeName
}

var (
	_ Def    = (*TypeDef)(nil)
	_ Stmt   = (*TypeDef)(nil)
	_ Parent = (*TypeDef)(nil)
)

func (*TypeDef) DefNode()  {}
func (*TypeDef) StmtNode() {}

func (n *TypeDef) Pos() token.Pos       { return n.TypeObj.Pos() }
func (n *TypeDef) Type() types.Type     { return n.TypeObj.Type() }
func (n *TypeDef) Object() types.Object { return n.TypeObj }
func (n *TypeDef) String() string       { return toString(n.TypeObj) }
func (n *TypeDef) ChildCount() int      { return len(n.Directives) }

func (n *TypeDef) Children(yield func(Node) bool) {
	_ = YieldSlice(n.Directives, yield)
}
