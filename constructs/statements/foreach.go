package statements

import (
	"fmt"

	"github.com/grant-nelson/Gozer/constructs/expressions"
)

var _ Statement = (*ForeachStat)(nil)

// ForeachStat is an for-each-statement.
type ForeachStat struct {

	// Key is the variable to put the mak key or list index into.
	Key expressions.Expression

	// Value is the variable to put the map value or list element into.
	Value expressions.Expression

	// Range is the value to iterate through.
	Range expressions.Expression

	// Body is the statement to call for each value in the map.
	Body Statement
}

// Foreach creates a new for-each-statement.
func Foreach(key expressions.Expression, value expressions.Expression,
	rangeExp expressions.Expression, body Statement) *ForeachStat {
	return &ForeachStat{
		Key:   key,
		Value: value,
		Range: rangeExp,
		Body:  body,
	}
}

// String gets the string for this constuct.
func (s *ForeachStat) String() string {
	if s == nil {
		return nilStr
	}
	return fmt.Sprint("foreach(",
		expressions.ToString(s.Key), ", ",
		expressions.ToString(s.Value), " in ",
		expressions.ToString(s.Range), ") ", ToString(s.Body))
}
