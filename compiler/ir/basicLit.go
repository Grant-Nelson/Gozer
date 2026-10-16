package ir

import (
	"go/token"
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/wad"
)

// BasicLit is a node that represents a literal of basic type.
type BasicLit struct {

	// ValuePos is the position of the literal.
	ValuePos token.Pos

	// TypeAndValue is the type and value of the literal.
	TypeAndValue *types.TypeAndValue
}

var (
	_ Expr         = (*BasicLit)(nil)
	_ Stmt         = (*BasicLit)(nil)
	_ wad.Waddable = (*BasicLit)(nil)
)

func (*BasicLit) ExprNode() {}
func (*BasicLit) StmtNode() {}

func (n *BasicLit) Pos() token.Pos   { return n.ValuePos }
func (n *BasicLit) Type() types.Type { return n.TypeAndValue.Type }
func (n *BasicLit) String() string   { return toString(n.TypeAndValue.Value) }

func (n *BasicLit) Wad(w *wad.Writer) {
	w.Pos(n.ValuePos)
	// TODO: Finish
}

func (n *BasicLit) UnWad(w *wad.Reader) {
	n.ValuePos = w.Pos()
	// TODO: Finish
}
