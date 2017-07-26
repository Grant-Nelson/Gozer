package statements

import (
	"fmt"

	"github.com/grant-nelson/Gozer/constructs/expressions"
)

var _ Statement = (*ForStat)(nil)

// ForStat is an for-statement.
type ForStat struct {

	// Init is the statement to call when the For loop starts.
	Init Statement

	// Cond is the expession for the conditional of this for.
	Cond expressions.Expression

	// Post is the statement to call when the For loop is about to re-evaluate.
	// Usually used for incrementing some state.
	Post Statement

	// Body is the statement to call for each loop.
	Body Statement
}

// For creates a new for-statement.
func For(init Statement, cond expressions.Expression, post Statement, body Statement) *ForStat {
	return &ForStat{
		Init: init,
		Cond: cond,
		Post: post,
		Body: body,
	}
}

// String gets the string for this constuct.
func (s *ForStat) String() string {
	if s == nil {
		return nilStr
	}
	init := ""
	if s.Init != nil {
		init = ToString(s.Init)
	}
	cond := ""
	if s.Cond != nil {
		cond = expressions.ToString(s.Cond)
	}
	post := ""
	if s.Post != nil {
		post = ToString(s.Post)
	}
	return fmt.Sprint("for(", init, "; ", cond, "; ", post, ") ", ToString(s.Body))
}
