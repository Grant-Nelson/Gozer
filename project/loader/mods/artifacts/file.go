package artifacts

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

// PackageName is the name of the package this file belongs too.
func PackageName(f *ast.File) string {
	return f.Name.Name
}

// PackagePath is the path of the package this file belongs too.
// This is the directly containing the file and may not always
// match the import path.
func PackagePath(fSet *token.FileSet, f *ast.File) string {
	return filepath.Dir(FilePath(fSet, f))
}

// File Path is the path to the file being modified.
// This should be the whole path including the package import path.
func FilePath(fSet *token.FileSet, f *ast.File) string {
	return fSet.Position(f.Pos()).Filename
}

// IsTest indicates this file is part of an package test,
// i.e. the file path ends with `_test.go`.
func IsTest(fSet *token.FileSet, f *ast.File) bool {
	return strings.HasSuffix(FilePath(fSet, f), `_test.go`)
}

// IsXTest indicates this file is part of an extra-package test,
// i.e. the package name ends with `_test`.
func IsXTest(f *ast.File) bool {
	return strings.HasSuffix(PackageName(f), `_test`)
}

// Empty indicates if the file was empty.
func Empty(f *ast.File) bool {
	return f.FileStart == f.FileEnd ||
		(len(f.Comments) <= 0 && len(f.Decls) <= 0)
}

// PackageKey gets the key for a package based on the package path and test flags.
func PackageKey(fSet *token.FileSet, f *ast.File) string {
	switch {
	case IsXTest(f):
		return PackagePath(fSet, f) + `#_XTest`
	case IsTest(fSet, f):
		return PackagePath(fSet, f) + `#_Test`
	default:
		return PackagePath(fSet, f)
	}
}
