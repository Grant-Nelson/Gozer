package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// Parser will parse a file.
// This is similar to [parser.ParseFile].
//
// This parser function can be used to modify how code is looked up,
// read, and parsed for the compiler. For example, the parser could be
// configured to attempt to read virtual files first then fall back to
// checking the OS for actual files. Or maybe the files are downloaded
// from a served source (although it would be slow to do that per file
// and instead a modifier should be used to shortcut parsing or to batch
// load the files.)
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
