package constructs

import (
	"strings"

	"github.com/grant-nelson/Gozer/common"
)

var _ Statement = (*BlockStatement)(nil)

// BlockStatement is a scoped block of statements.
type BlockStatement struct {

	// Statements are the set of statements in this block.
	Statements []Statement
}

// Block creates a new block statment.
func Block(statements ...Statement) *BlockStatement {
	return &BlockStatement{
		Statements: statements,
	}
}

// String gets the string for this block.
func (s *BlockStatement) String() string {
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
