package artifacts

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// Parser is the interface for a file parser.
type Parser interface {

	// Parse will parse a file.
	// This is similar to [parser.ParseFile].
	Parse(fileSet *token.FileSet, filename string, src any) (*ast.File, error)
}

// DefaultParser will load a file using [parser.ParseFile].
var DefaultParser defaultParser

var _ Parser = DefaultParser

type defaultParser struct{}

func (defaultParser) Parse(fileSet *token.FileSet, filename string, src any) (*ast.File, error) {
	const mode = parser.AllErrors |
		parser.ParseComments |
		parser.DeclarationErrors |
		parser.SkipObjectResolution
	return parser.ParseFile(fileSet, filename, src, mode)
}
