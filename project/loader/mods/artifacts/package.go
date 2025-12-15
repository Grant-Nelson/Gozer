package artifacts

import "go/token"

type Package struct {
	Name string
	Path string

	TempFileSet *token.FileSet

	// TODO: Need to indicate if this package includes test in the same package.
	// TODO: Need to indicate if this package is a test package.
}
