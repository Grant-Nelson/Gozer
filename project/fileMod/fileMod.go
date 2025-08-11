package fileMod

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"path/filepath"
	"slices"
	"strings"
)

// FileMod that is being loaded and modified.
type FileMod struct {

	// filename is the name of the first file being loaded.
	filename string

	// pkgName is the name of the package this file belongs to.
	pkgName string

	// pkgPath is the cache of the directory for this file that is also
	// the import path for the package this file belong to.
	pkgPath string

	// fileSet is used to set the tracing for the file.
	// This is a temporary file set specific to storing this file
	// and additional files while loading it.
	fileSet *token.FileSet

	// Doc are the document level comments.
	Doc []string

	// imports is all the imports in this file.
	// The imports should be merged and named such that there is only one
	// import spec for any specific path. All Decls need to be updated
	// to use these imports.
	Imports []*ast.ImportSpec

	// decls are the top-level declarations not including any imports.
	Decls []ast.Decl
}

// New creates a new FileMod.
func New(filename string) *FileMod {
	return &FileMod{
		filename: filename,
		fileSet:  token.NewFileSet(),
	}
}

const parseMode = parser.AllErrors | parser.ParseComments | parser.SkipObjectResolution

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

// Filename is the name of the first file being loaded.
// This should be the whole path including the package import path.
func (fm *FileMod) Filename() string { return fm.filename }

// PackageName is the name of the package this file belongs too.
func (fm *FileMod) PackageName() string { return fm.pkgName }

// FileSet is used to set the tracing for the file.
// This is a temporary file set specific to storing this file
// and additional files while loading it.
func (fm *FileMod) FileSet() *token.FileSet { return fm.fileSet }

// PackagePath is the import path for the package this file belong to.
func (fm *FileMod) PackagePath() string {
	if len(fm.pkgPath) <= 0 {
		fm.pkgPath = filepath.Dir(fm.filename)
	}
	return fm.pkgPath
}

func (fm *FileMod) AddDoc(text string) {
	if cut, ok := strings.CutPrefix(text, `// `); ok {
		text = cut
	}

	// TODO: Normalize /* comment
	fm.Doc = append(fm.Doc, text)
}

func (fm *FileMod) AddDocGroup(doc *ast.CommentGroup) {
	if doc != nil {
		fm.Doc = slices.Grow(fm.Doc, len(doc.List))
		for _, c := range doc.List {
			fm.AddDoc(c.Text)
		}
	}
}

func (fm *FileMod) mergeImport(imp *ast.ImportSpec) {
	// TODO: Need to normalize imports by removing duplicates and setting unique names.
	fm.Imports = append(fm.Imports, imp)
}

func (fm *FileMod) addDecl(decl ast.Decl) {
	// TODO: Need to normalize decls to have correct imports.
	fm.Decls = append(fm.Decls, decl)
}

// AddDecls adds declarations to this file.
// The declarations must use positions in the file set.
func (fm *FileMod) AddDecls(decls []ast.Decl) {
	for _, decl := range decls {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.IMPORT {
			for _, spec := range gen.Specs {
				fm.mergeImport(spec.(*ast.ImportSpec))
			}
		} else {
			fm.addDecl(decl)
		}
	}
}

func (fm *FileMod) AddFile(filename string, src []byte) error {
	if len(fm.filename) <= 0 {
		fm.filename = filename
	}

	f, err := parser.ParseFile(fm.FileSet(), filename, src, parseMode)
	if err != nil {
		return err
	}

	if len(fm.pkgName) <= 0 {
		fm.pkgName = f.Name.Name
	}

	fm.AddDocGroup(f.Doc)
	fm.AddDecls(f.Decls)
	return nil
}

func (fm *FileMod) Write(out io.Writer) (err error) {
	defer recoverError(`error writing file`, &err)

	write := func(text string) {
		if _, err := out.Write([]byte(text)); err != nil {
			panic(err)
		}
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

func (fm *FileMod) Finalize(fileSet *token.FileSet) (f *ast.File, err error) {
	defer recoverError(`error finalizing file`, &err)
	buf := &bytes.Buffer{}
	if err := fm.Write(buf); err != nil {
		return nil, err
	}

	f, err = parser.ParseFile(fileSet, fm.Filename(), buf.Bytes(), parseMode)
	return f, err
}
