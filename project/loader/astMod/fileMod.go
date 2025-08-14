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

func (fm *FileMod) Write(out io.Writer) (err error) {
	write := func(text string) error {
		_, err := out.Write([]byte(text))
		return err
	}

	for _, doc := range fm.Doc {
		write(`// ` + doc + "\n")
	}
	write(`package ` + fm.pkgName + "\n\n")

	if len(fm.Imports) == 1 {
		write("import ")
		if err := printer.Fprint(out, fm.FileSet(), fm.Imports[0]); err != nil {
			panic(err)
		}
		write("\n\n")
	} else if len(fm.Imports) > 1 {
		write("import (\n")
		for _, im := range fm.Imports {
			if err := printer.Fprint(out, fm.FileSet(), im); err != nil {
				panic(err)
			}
		}
		write(")\n\n")
	}

	// TODO: Need to handle nil decl/spec values.
	// TODO: Ensure that iota is being handled correctly.
	for _, decl := range fm.Decls {
		if p := fm.FileSet().Position(decl.Pos()); p.IsValid() {
			write(`//line ` + p.String() + "\n")
		}
		// TODO: Need to handle rename and replace signature by adding the line:column at the
		//       start of a body if the body doesn't offset correctly from the signature.
		if err := printer.Fprint(out, fm.FileSet(), decl); err != nil {
			panic(err)
		}
		write("\n") // TODO: Use error group
	}
	return nil
}

func (fm *FileMod) Finalize(fileSet *token.FileSet, mode parser.Mode) (f *ast.File, err error) {
	buf := &bytes.Buffer{}
	if err := fm.Write(buf); err != nil {
		return nil, fm.Package().ErrorGroup().Add(err)
	}

	f, err = parser.ParseFile(fileSet, fm.Path(), buf.Bytes(), mode)
	if err != nil {
		return nil, fm.Package().ErrorGroup().Add(err)
	}

	return f, nil
}
