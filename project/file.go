package project

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// File that is being loaded and modified.
type File struct {

	// Filename is the name of the file being loaded.
	Filename string

	// FileSet is used to set the tracing for the file.
	// This is a temporary file set specific to storing this file
	// and additional files while loading it.
	FileSet *token.FileSet

	// Doc are the document level comments.
	Doc []*ast.Comment

	// Name is the name of the package this file belongs too.
	Name *ast.Ident

	// All the imports in this file.
	Imports []*ast.ImportSpec

	// Decls are the top-level declarations not including any imports.
	Decls []ast.Decl
}

func comments(cg *ast.CommentGroup) ([]*ast.Comment, bool) {
	if cg != nil && len(cg.List) > 0 {
		return cg.List, true
	}
	return []*ast.Comment{}, false
}

func initFile(filename string, src []byte) (*File, error) {
	fileSet := token.NewFileSet()
	const mode = parser.AllErrors | parser.ParseComments | parser.SkipObjectResolution
	f, err := parser.ParseFile(fileSet, filename, src, mode)
	if err != nil {
		return nil, err
	}

	// Get all declarations without any imports.
	decls := make([]ast.Decl, 0, len(f.Decls))
	for _, decl := range f.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.IMPORT {
			// If the import decl has a comment, move it onto the first import.
			if genComments, ok := comments(gen.Doc); ok && len(gen.Specs) > 0 {
				firstImp := gen.Specs[0].(*ast.ImportSpec)
				impComment, _ := comments(firstImp.Comment)
				firstImp.Comment = &ast.CommentGroup{
					List: append(genComments, impComment...),
				}
			}
		} else {
			decls = append(decls, decl)
		}
	}

	doc, _ := comments(f.Doc)

	return &File{
		Filename: filename,
		FileSet:  fileSet,
		Doc:      doc,
		Name:     f.Name,
		Imports:  f.Imports,
		Decls:    decls,
	}, nil
}

func (file *File) finalize(fileSet *token.FileSet) (*ast.File, error) {
	// Doc *ast.CommentGroup
	// Name *ast.Ident
	// Imports []*ast.ImportSpec
	// Decls []ast.Decl

	totalSize := 0
	for _, dc := range file.Doc.List {
		//dc.
	}

	f := &ast.File{}

	return f, nil
}
