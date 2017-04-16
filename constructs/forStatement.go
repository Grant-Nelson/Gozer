package constructs

import (
	"fmt"
	"strings"
)

var _ Statement = (*ForStat)(nil)

// ForStat is an for statement.
type ForStat struct {

	// Init is the statement to call when the For loop starts.
	Init Statement

	// Cond is the expession for the conditional of this for.
	Cond Expression

	// Post is the statement to call when the For loop is about to re-evaluate.
	// Usually used for incrementing some state.
	Post []Statement

	// Body is the statement to call when Cond evaluates to true.
	Body Statement
}

// For creates a new for-statement.
func For(init Statement, cond Expression, post []Statement, body Statement) *ForStat {
	return &ForStat{
		Init: init,
		Cond: cond,
		Post: post,
		Body: body,
	}
}

// statsToString gets the string for a slice of statements.
func (s *ForStat) statsToString(stats []Statement) string {
	if (stats != nil) && (len(stats) > 0) {
		if len(stats) <= 1 {
			return ToString(stats[0])
		}
		parts := make([]string, len(stats))
		for i, part := range stats {
			parts[i] = ToString(part)
		}
		return "{" + strings.Join(parts, "; ") + "}"
	}
	return ""
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
		cond = ToString(s.Cond)
	}
	post := s.statsToString(s.Post)
	return fmt.Sprint("for(", init, "; ", cond, "; ", post, ") ", ToString(s.Body))
}
