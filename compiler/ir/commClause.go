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

var _ Node = (*CommClause)(nil)

func (s *CommClause) String() string {
	str := `default:`
	if s.Comm != nil {
		str = fmt.Sprintf(`case %v:`, s.Comm)
	}
	if len(s.Body) > 0 {
		str += "\n" + linesString(s.Body)
	}
	return str
}

func (c *CommClause) Pos() token.Pos { return astPos(c.Ast) }

func FromCommClause(s *ast.CommClause, c *Converter) *CommClause {
	if s == nil {
		return nil
	}
	return &CommClause{
		Ast:  s,
		Comm: c.FromStmt(s.Comm),
		Body: c.FromStmtSlice(s.Body),
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
