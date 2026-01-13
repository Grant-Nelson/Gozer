package artifacts

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

// PositionResolver is for either a `token.FileSet` or a `token.File`.
type PositionResolver interface {
	Position(token.Pos) token.Position
}

// PackageName is the name of the package this file belongs too.
func PackageName(f *ast.File) string {
	return f.Name.Name
}

// PackagePath is the path of the package this file belongs too.
// This is the directly containing the file and may not always
// match the import path.
// The node should typically be an *ast.File but may be any node in that file.
func PackagePath(pr PositionResolver, n ast.Node) string {
	return filepath.Dir(FilePath(pr, n))
}

// File Path is the path to the file being modified.
// This should be the whole path including the package import path.
// The node should typically be an *ast.File but may be any node in that file.
func FilePath(pr PositionResolver, n ast.Node) string {
	return pr.Position(n.Pos()).Filename
}

// IsTest indicates this file is part of an package test,
// i.e. the file path ends with `_test.go`.
// The node should typically be an *ast.File but may be any node in that file.
func IsTest(pr PositionResolver, n ast.Node) bool {
	return strings.HasSuffix(FilePath(pr, n), `_test.go`)
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
func PackageKey(pr PositionResolver, f *ast.File) string {
	switch {
	case IsXTest(f):
		return PackagePath(pr, f) + `#_XTest`
	case IsTest(pr, f):
		return PackagePath(pr, f) + `#_Test`
	default:
		return PackagePath(pr, f)
	}
}
