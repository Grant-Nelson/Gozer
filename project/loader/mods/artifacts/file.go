package artifacts

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

// File that is being modified or inspected.
type File struct {

	// Package that this file is part of.
	Package *Package

	// TempFileSet that is associated with this file.
	// This file set may be unique for this file during loading.
	TempFileSet *token.FileSet

	// File is the file's ast being modified.
	File *ast.File
}

// NewFile creates a new file for the modifier.
// This will create a temporary package object for this file.
func NewFile(tempFileSet *token.FileSet, file *ast.File) *File {
	f := &File{
		TempFileSet: tempFileSet,
		File:        file,
	}
	f.Package = NewPackageForFile(f)
	return f
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
