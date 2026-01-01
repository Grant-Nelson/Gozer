package artifacts

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"io"
	"path/filepath"
	"strings"
)

// File that is being modified or inspected.
type File struct {

	// Package that this file is part of.
	Package *Package

	// TempFileSet that is associated with this file.
	// This file set may be unique for this file during loading.
	TempFileSet *FileSet

	// File is the file's ast being modified.
	File *ast.File
}

// New creates a new file mod.
func New(tempFileSet *FileSet, file *ast.File) *File {
	f := &File{
		TempFileSet: tempFileSet,
		File:        file,
	}
	f.Package = NewPackageForFile(f)
	f.TempFileSet.RegisterFile(f)
	return f
}

// Load will load a file using parser.ParseFile.
func Load(fileSet *FileSet, filename string, src any) (*File, error) {
	const mode = parser.AllErrors |
		parser.ParseComments |
		parser.DeclarationErrors |
		parser.SkipObjectResolution
	f, err := parser.ParseFile(fileSet.fileSet, filename, src, mode)
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
	return f.TempFileSet.Position(f.File.Pos()).Filename
}

// IsTest indicates this file is part of an package test,
// i.e. the file path ends with `_test.go`.
func (f *File) IsTest() bool {
	return strings.HasSuffix(f.FilePath(), `_test.go`)
}

// IsXTest indicates this file is part of an extra-package test,
// i.e. the package name ends with `_test`.
func (f *File) IsXTest() bool {
	return strings.HasSuffix(f.PackageName(), `_test`)
}

// Empty indicates if the file was empty.
func (f *File) Empty() bool {
	return f.File.FileStart == f.File.FileEnd ||
		(len(f.File.Comments) <= 0 && len(f.File.Decls) <= 0)
}

// Write will write the modified file to the given writer.
//
// This will not use the error group and returns any errors that occurred.
func (f *File) Write(out io.Writer) error {
	cfg := &printer.Config{
		Mode:     printer.TabIndent, // | printer.SourcePos,
		Tabwidth: 4,
	}
	return cfg.Fprint(out, f.TempFileSet.fileSet, f.File)
}

// Reload will write the file to a temporary buffer and reload it
// with the given file set to normalize the file information.
func (f *File) Reload(fileSet *FileSet) (*File, error) {
	buf := &bytes.Buffer{}
	if err := f.Write(buf); err != nil {
		return nil, err
	}
	return Load(fileSet, f.FilePath(), buf.Bytes())
}

// PackageKey gets the key for a package based on the package path and test flags.
func (f *File) PackageKey() string {
	switch {
	case f.IsXTest():
		return f.PackagePath() + `#_XTest`
	case f.IsTest():
		return f.PackagePath() + `#_Test`
	default:
		return f.PackagePath()
	}
}

// Directives finds all the directives with the given prefix.
func Directives(comments []*ast.Comment, prefix string) map[string][]string {
	prefix = `//` + prefix + `:`
	result := map[string][]string{}
	for _, c := range comments {
		if tail, ok := strings.CutPrefix(c.Text, prefix); ok {
			var key, value string
			if i := strings.Index(tail, ` `); i > 0 {
				key = strings.TrimSpace(tail[:i])
				value = strings.TrimSpace(tail[i:])
			} else {
				key = tail
				value = ``
			}
			result[key] = append(result[key], value)
		}
	}
	return result
}

func RemoveDirectives(comments []*ast.Comment, prefix string) []*ast.Comment {
	prefix = `//` + prefix + `:`
	result := make([]*ast.Comment, 0, len(comments))
	for _, c := range comments {
		if !strings.HasPrefix(c.Text, prefix) {
			result = append(result, c)
		}
	}
	return result
}
