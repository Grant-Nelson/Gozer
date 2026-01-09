package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// Parser will parse a file.
// This is similar to [parser.ParseFile].
type Parser func(fileSet *token.FileSet, filename string, src any) (*ast.File, error)

var _ Parser = Default

// Default will load a file using [parser.ParseFile].
func Default(fileSet *token.FileSet, filename string, src any) (*ast.File, error) {
	const mode = parser.AllErrors |
		parser.ParseComments |
		parser.DeclarationErrors |
		parser.SkipObjectResolution
	return parser.ParseFile(fileSet, filename, src, mode)
}
