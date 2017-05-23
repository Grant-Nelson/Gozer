package statements

import (
	"strings"

	"github.com/grant-nelson/Gozer/common"
)

var _ Statement = (*BlockStat)(nil)

// BlockStat is a scoped block of statements.
type BlockStat struct {

	// Statements are the set of statements in this block.
	Statements []Statement
}

// Block creates a new block statment.
func Block(statements ...Statement) *BlockStat {
	return &BlockStat{
		Statements: statements,
	}
}

// String gets the string for this block.
func (s *BlockStat) String() string {
	if s == nil {
		return nilStr
	}
	statLen := len(s.Statements)
	if statLen <= 0 {
		return "{}"
	}
	parts := make([]string, statLen)
	for i, stat := range s.Statements {
		parts[i] = ToString(stat)
	}
	return "{\n  " + common.Indent(strings.Join(parts, "\n"), "  ") + "\n}"
}
