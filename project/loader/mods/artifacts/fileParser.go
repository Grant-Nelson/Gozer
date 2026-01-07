package artifacts

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// FileParser is the interface for a file parser.
type FileParser interface {

	// Parse will parse a file.
	// This is similar to [parser.ParseFile].
	Parse(fileSet *token.FileSet, filename string, src any) (*ast.File, error)
}

// DefaultFileParser will load a file using [parser.ParseFile].
var DefaultFileParser defaultFileParser

var _ FileParser = DefaultFileParser

type defaultFileParser struct{}

func (defaultFileParser) Parse(fileSet *token.FileSet, filename string, src any) (*ast.File, error) {
	const mode = parser.AllErrors |
		parser.ParseComments |
		parser.DeclarationErrors |
		parser.SkipObjectResolution
	f, err := parser.ParseFile(fileSet, filename, src, mode)
	return f, err
}
