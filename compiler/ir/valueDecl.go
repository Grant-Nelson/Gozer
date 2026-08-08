package ir

import (
	"go/token"
)

// ValueDecl is the declaration for a single variable or constant.
type ValueDecl struct {

	// Constant indicates the value is constant or variable.
	Constant bool

	// Name is the name for the value.
	Name *Ident

	// Value is the optional initial assigned of the value.
	Value Expr
}

var (
	_ Stmt   = (*ValueDecl)(nil)
	_ Parent = (*ValueDecl)(nil)
)

func (n *ValueDecl) Pos() token.Pos { return n.Name.Pos() }

func (*ValueDecl) StmtNode() {}

func (n *ValueDecl) String() string {
	result := `var `
	if n.Constant {
		result = `const `
	}
	result += n.Name.String() + ` ` + n.Name.Type().String()
	if n.Value != nil {
		result += ` = ` + n.Value.String()
	}
	return result
}

func (n *ValueDecl) Children(yield func(Node) bool) {
	_ = yield(n.Name) && yield(n.Value)
}
