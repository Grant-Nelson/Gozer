package file

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"path/filepath"
	"strings"
)

// File that is being modified or inspected.
type File struct {

	// FileSet that is associated with this file.
	// This file set may be unique for this file during loading.
	FileSet *token.FileSet

	// File is the file's ast being modified.
	File *ast.File
}

// New creates a new file mod.
func New(fileSet *token.FileSet, file *ast.File) *File {
	return &File{
		FileSet: fileSet,
		File:    file,
	}
}

// Load will load a file using parser.ParseFile.
func Load(fileSet *token.FileSet, filename string, src any) (*File, error) {
	const mode = parser.AllErrors |
		parser.ParseComments |
		parser.DeclarationErrors |
		parser.SkipObjectResolution
	f, err := parser.ParseFile(fileSet, filename, src, mode)
	return New(fileSet, f), err
}

// PackageName is the name of the package this file belongs too.
func (f *File) PackageName() string {
	return f.File.Name.Name
}

// PackagePath is the path of the package this file belongs too.
// This is the directly containing the file and may not always
// match the import path.
func (f *File) PackagePath() string {
	return filepath.Dir(f.FilePath())
}

// File Path is the path to the file being modified.
// This should be the whole path including the package import path.
func (f *File) FilePath() string {
	return f.FileSet.Position(f.File.Pos()).Filename
}

// Write will write the modified file to the given writer.
//
// This will not use the error group and returns any errors that occurred.
func (f *File) Write(out io.Writer) error {
	cfg := &printer.Config{
		Mode:     printer.TabIndent | printer.SourcePos,
		Tabwidth: 4,
	}
	return cfg.Fprint(out, f.FileSet, f.File)
}

// Reload will write the file to a temporary buffer and reload it
// with the given file set to normalize the file information.
func (f *File) Reload(fileSet *token.FileSet) (*File, error) {
	buf := &bytes.Buffer{}
	if err := f.Write(buf); err != nil {
		return nil, err
	}
	return Load(fileSet, f.FilePath(), buf.Bytes())
}

// Directives finds all the directives with the given prefix.
func Directives(comments []*ast.Comment, prefix string) map[string][]string {
	prefix += `:`
	result := map[string][]string{}
	for _, c := range comments {
		if tail, ok := strings.CutPrefix(c.Text, prefix); ok {
			if i := strings.Index(tail, ` `); i > 0 {
				key := strings.TrimSpace(tail[:i])
				value := strings.TrimSpace(tail[i:])
				result[key] = append(result[key], value)
			}
		}
	}
	return result
}
