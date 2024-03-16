package reader

import (
	"context"
	"go/ast"
	"go/token"
)

// Config is the read and parse configuration.
type Config struct {

	// Verbose indicates that the parse logs
	// should be written to the standard out.
	Verbose bool

	// Path is the path to the main package or primary package.
	// The path should contain the mod file.
	// The path follows the standard patter for go tools.
	Path string

	// Context is the optional context to cancel a build with.
	Context context.Context

	// Tests indicates the testing files and packages should also be built.
	Tests bool

	// BuildFlags are the optional build flags to build with.
	// Example: // +build tag_name
	BuildFlags []string

	// AugmentFile is an optional hook to allow for alterations to the
	// individual files prior to being analyzed during the build.
	AugmentFile func(args *AugmentFileArgs) error
}

// AugmentFileArgs is the arguments given for augmenting a file.
type AugmentFileArgs struct {

	// Filename is the name of the file being augmented.
	Filename string

	// FileSet is the file set for read and arsed files.
	FileSet *token.FileSet

	// File is the AST tree to modify to augment the file.
	// Replacing this pointer will not affect the file.
	File *ast.File
}
