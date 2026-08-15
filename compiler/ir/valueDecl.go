package ir

import (
	"go/token"
	"go/types"
)

// ValueDecl is the declaration for a single variable or constant.
type ValueDecl struct {

	// Constant indicates the value is constant or variable.
	Constant bool

	// Name is the name for the value.
	Name string

	// NamePos is the position for the name of the value.
	NamePos token.Pos

	// TypeAndValue is the type and constant value.
	TypeAndValue types.TypeAndValue

	// Value is the optional initial assigned of the value or nil.
	Value Expr
}

var (
	_ Stmt   = (*ValueDecl)(nil)
	_ Decl   = (*ValueDecl)(nil)
	_ Parent = (*ValueDecl)(nil)
)

func (n *ValueDecl) Pos() token.Pos   { return n.NamePos }
func (n *ValueDecl) Type() types.Type { return n.TypeAndValue.Type }

func (*ValueDecl) StmtNode() {}
func (*ValueDecl) DeclNode() {}

func (n *ValueDecl) String() string {
	result := `var `
	if n.Constant {
		result = `const `
	}
	result += n.Name + ` ` + n.TypeAndValue.Type.String()
	if n.Value != nil {
		result += ` = ` + n.Value.String()
	}
	return result
}

func (n *ValueDecl) Children(yield func(Node) bool) { _ = yield(n.Value) }
