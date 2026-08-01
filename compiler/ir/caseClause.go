package ir

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

// CaseClause is a node that represents a case of an expression or type switch statement.
type CaseClause struct {
	Ast  *ast.CaseClause // TODO: REMOVE
	List []ast.Expr      // TODO: REPLACE // list of expressions or types; nil means default case
	Body []Stmt          // statement list; or nil
}

var (
	_ Node   = (*CaseClause)(nil)
	_ Parent = (*CaseClause)(nil)
)

func (n *CaseClause) String() string {
	str := `default:`
	if len(n.List) > 0 {
		str = fmt.Sprintf(`case %v:`, csvString(n.List))
	}
	if len(n.Body) > 0 {
		str += "\n" + linesString(n.Body)
	}
	return str
}

func (n *CaseClause) Pos() token.Pos { return astPos(n.Ast) }

func (n *CaseClause) Children(yield func(Node) bool) bool {
	return YieldSlice(n.List, yield) && YieldSlice(n.Body, yield)
}

func FromCaseClause(s *ast.CaseClause, c *Converter) *CaseClause {
	if s == nil {
		return nil
	}
	return &CaseClause{
		Ast:  s,
		List: s.List,
		Body: FromStmtSlice(s.Body, c),
	}
}

func FromCaseClauseSlice(s *ast.BlockStmt, c *Converter) []*CaseClause {
	if s == nil {
		return []*CaseClause{}
	}
	ccs := make([]*CaseClause, 0, len(s.List))
	for i, ct := range s.List {
		if ct == nil {
			continue
		}
		cc, ok := ct.(*ast.CaseClause)
		if !ok {
			panic(faults.New(`expected case clause`).
				WithF(`type`, `%T`, ct).
				With(`index`, i))
		}
		ccs = append(ccs, FromCaseClause(cc, c))
	}
	return ccs
}
