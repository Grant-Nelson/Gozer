package irc

import (
	"fmt"
	"go/token"
)

type (
	// CaseClause is a node for a case of an expression or type switch statement.
	CaseClause struct {
		Case  token.Pos // position of "case" or "default" keyword
		List  []Expr    // list of expressions or types; nil means default case
		Colon token.Pos // position of ":"
		Body  []Stmt    // statement list; or empty
	}

	// CommClause is a node for a select statement.
	CommClause struct {
		Case  token.Pos // position of "case" or "default" keyword
		Comm  Stmt      // send or receive statement; nil means default case
		Colon token.Pos // position of ":"
		Body  []Stmt    // statement list; or empty
	}
)

var (
	_ Node = (*CaseClause)(nil)
	_ Node = (*CommClause)(nil)
)

//====[String]==================================================================

func (s *CaseClause) String() string {
	str := `default:`
	if len(s.List) > 0 {
		str = fmt.Sprintf(`case %v:`, csvString(s.List))
	}
	if len(s.Body) > 0 {
		str += "\n" + linesString(s.Body, `  `)
	}
	return str
}
func (s *CommClause) String() string {
	str := `default:`
	if s.Comm != nil {
		str = fmt.Sprintf(`case %v:`, s.Comm)
	}
	if len(s.Body) > 0 {
		str += "\n" + linesString(s.Body, `  `)
	}
	return str
}

//====[Pos]=====================================================================

func (s *CaseClause) Pos() token.Pos { return s.Case }
func (s *CommClause) Pos() token.Pos { return s.Case }

//====[End]=====================================================================

func (s *CaseClause) End() token.Pos {
	if p := endOfSlice(s.Body); p.IsValid() {
		return p
	}
	return s.Colon + 1 // len(":")
}
func (s *CommClause) End() token.Pos {
	if p := endOfSlice(s.Body); p.IsValid() {
		return p
	}
	return s.Colon + 1 // len(":")
}
