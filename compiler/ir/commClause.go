package ir

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

// CommClause is a node that represents a case of a select statement.
type CommClause struct {
	Ast  *ast.CommClause // TODO: REMOVE
	Comm Stmt            // send or receive statement; nil means default case
	Body []Stmt          // statement list; or nil
}

var (
	_ Node   = (*CommClause)(nil)
	_ Parent = (*CommClause)(nil)
)

func (n *CommClause) String() string {
	str := `default:`
	if n.Comm != nil {
		str = fmt.Sprintf(`case %v:`, n.Comm)
	}
	if len(n.Body) > 0 {
		str += "\n" + linesString(n.Body)
	}
	return str
}

func (n *CommClause) Pos() token.Pos { return astPos(n.Ast) }

func (n *CommClause) Children(yield func(Node) bool) bool {
	return yield(n.Comm) && YieldSlice(n.Body, yield)
}

func FromCommClause(s *ast.CommClause, c *Converter) *CommClause {
	if s == nil {
		return nil
	}
	return &CommClause{
		Ast:  s,
		Comm: FromStmt(s.Comm, c),
		Body: FromStmtSlice(s.Body, c),
	}
}

func FromCommClauseSlice(s *ast.BlockStmt, c *Converter) []*CommClause {
	if s == nil {
		return []*CommClause{}
	}
	ccs := make([]*CommClause, 0, len(s.List))
	for i, ct := range s.List {
		if ct == nil {
			continue
		}
		cc, ok := ct.(*ast.CommClause)
		if !ok {
			panic(faults.New(`expected comm clause`).
				WithF(`type`, `%T`, ct).
				With(`index`, i))
		}
		ccs = append(ccs, FromCommClause(cc, c))
	}
	return ccs
}
