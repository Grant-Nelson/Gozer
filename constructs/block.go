package constructs

import "strings"

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

// Equals determines if this block is the same as the other block.
func (s *BlockStatement) Equals(other interface{}) bool {
	if s == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	sother, ok := other.(*BlockStatement)
	if !ok {
		return false
	}
	if len(s.Statements) != len(sother.Statements) {
		return false
	}
	for i, stat := range s.Statements {
		if !Equals(stat, sother.Statements[i]) {
			return false
		}
	}
	return true
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
	return "{\n   " + strings.Join(parts, "\n   ") + "}"
}
