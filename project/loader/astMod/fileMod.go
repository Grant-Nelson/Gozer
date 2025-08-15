package astMod

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
)

// FileMod that is being loaded and modified.
type FileMod struct {

	// pkg is te package this file will belong to.
	pkg *PackageMod

	// path is the path of the first file being loaded.
	path string

	// file is the file's ast being modified.
	file *ast.File
}

// NewFile creates a new temporary file for modifying the file's AST.
func NewFile(path string, file *ast.File, pkg *PackageMod) *FileMod {
	return &FileMod{
		pkg:  pkg,
		path: path,
		file: file,
	}
}

// Package is the name of the package this file belongs too.
func (fm *FileMod) Package() *PackageMod { return fm.pkg }

// Path is the path to the file being modified.
// This should be the whole path including the package import path.
func (fm *FileMod) Path() string { return fm.path }

// File is the AST of the file being modified.
func (fm *FileMod) File() *ast.File { return fm.file }

// Write will write the modified file to the given writer.
//
// This will not use the error group and returns any errors that occurred.
func (fm *FileMod) Write(out io.Writer) error {
	cfg := &printer.Config{
		Mode:     printer.TabIndent | printer.SourcePos,
		Tabwidth: 4,
	}
	return cfg.Fprint(out, fm.Package().FileSet(), fm.file)
}

func (fm *FileMod) Finalize(fileSet *token.FileSet, mode parser.Mode) (*ast.File, error) {
	buf := &bytes.Buffer{}
	if err := fm.Write(buf); err != nil {
		return nil, fm.Package().ErrorGroup().Add(err)
	}

	f, err := parser.ParseFile(fileSet, fm.Path(), buf.Bytes(), mode)
	if err != nil {
		return nil, fm.Package().ErrorGroup().Add(err)
	}
	return f, nil
}
