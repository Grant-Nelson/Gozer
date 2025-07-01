package project

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
)

// File that is being loaded and modified.
type File struct {

	// Filename is the name of the first file being loaded.
	Filename string

	// Package is the name of the package this file belongs too.
	Package string

	// FileSet is used to set the tracing for the file.
	// This is a temporary file set specific to storing this file
	// and additional files while loading it.
	FileSet *token.FileSet

	// Doc are the document level comments.
	Doc []*ast.Comment

	// Imports is all the imports in this file.
	// The imports should be merged and named such that there is only one
	// import spec for any specific path. All Decls need to be updated
	// to use these imports.
	Imports []*ast.ImportSpec

	// Decls are the top-level declarations not including any imports.
	Decls []ast.Decl
}

const (
	parseMode = parser.AllErrors | parser.ParseComments | parser.SkipObjectResolution
)

func comments(cg *ast.CommentGroup) ([]*ast.Comment, bool) {
	if cg != nil && len(cg.List) > 0 {
		return cg.List, true
	}
	return []*ast.Comment{}, false
}

func initFile(filename string, src []byte) (*File, error) {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, filename, src, parseMode)
	if err != nil {
		return nil, err
	}

	fo := &File{
		Filename: filename,
		FileSet:  fileSet,
		Package:  f.Name.Name,
		Imports:  f.Imports,
	}

	// Get all declarations without any imports.
	fo.Decls = make([]ast.Decl, 0, len(f.Decls))
	for _, decl := range f.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.IMPORT {
			// If the import decl has a comment, move it onto the first import spec for that group.
			if genComments, ok := comments(gen.Doc); ok && len(gen.Specs) > 0 {
				firstImp := gen.Specs[0].(*ast.ImportSpec)
				impComment, _ := comments(firstImp.Comment)
				firstImp.Comment = &ast.CommentGroup{
					List: append(genComments, impComment...),
				}
			}
		} else {
			fo.Decls = append(fo.Decls, decl)
		}
	}

	fo.Doc, _ = comments(f.Doc)
	return fo, nil
}

func recoverError(msg string, err *error) {
	if r := recover(); r != nil {
		switch t := r.(type) {
		case error:
			*err = fmt.Errorf(msg+`: %w`, t)
		default:
			*err = fmt.Errorf(msg+`: %v`, r)
		}
	}
}

func (file *File) write(out io.Writer) (err error) {
	defer recoverError(`error writing file`, &err)

	write := func(text string) {
		if _, err := out.Write([]byte(text)); err != nil {
			panic(err)
		}
	}

	for _, doc := range file.Doc {
		write(doc.Text + "\n")
	}
	write(`package ` + file.Package + "\n\nimport (\n")
	for _, im := range file.Imports {
		if err := printer.Fprint(out, file.FileSet, im); err != nil {
			panic(err)
		}
	}
	write(")\n\n")
	for _, decl := range file.Decls {
		if p := file.FileSet.Position(decl.Pos()); p.IsValid() {
			write(`//line ` + p.String() + "\n")
		}
		if err := printer.Fprint(out, file.FileSet, decl); err != nil {
			panic(err)
		}
		write("\n")
	}
	return nil
}

func (file *File) finalize(fileSet *token.FileSet) (f *ast.File, err error) {
	defer recoverError(`error finalizing file`, &err)
	buf := &bytes.Buffer{}
	if err := file.write(buf); err != nil {
		return nil, err
	}

	f, err = parser.ParseFile(fileSet, file.Filename, buf.Bytes(), parseMode)
	return f, err
}
