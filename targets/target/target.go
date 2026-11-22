package target

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/compiler/project"
)

type ListConfig struct {

	// Lang is the language to transpile into, e.g. `ts` or `typescript`.
	// This is case insensitive.
	Lang string

	// Logger to log verbose messages with. Has no affect if verbose was false.
	Logger *logger.Logger

	// ErrGroup is the collector to handle multiple errors.
	// This will be returned from the [Build] method.
	// If nil, a error group will be created with the default limit.
	ErrGroup *faults.ErrGroup

	// Dir is the directory in which to run the build system's query tool
	// that provides information about the packages.
	// If Dir is empty, the tool is run in the current directory.
	Dir string

	// Build is a list of command-line flags to be passed through to
	// the build system's query tool.
	Build []string

	// Patterns are the file patterns to use for the project root.
	Patterns []string

	// Tests indicates the test files will also be read for the patterns.
	Tests bool
}

type BuildConfig struct {

	// Lang is the language to transpile into, e.g. `ts` or `typescript`.
	// This is case insensitive.
	Lang string

	// Logger to log verbose messages with. Has no affect if verbose was false.
	Logger *logger.Logger

	// ErrGroup is the collector to handle multiple errors.
	// This will be returned from the [Build] method.
	// If nil, a error group will be created with the default limit.
	ErrGroup *faults.ErrGroup

	// Output is the directory to copy the resulting application out to.
	//
	// If empty, then the compiled code is left in the default output directory
	// that is used for cache information as well.
	Output string

	// Dir is the directory in which to run the build system's query tool
	// that provides information about the packages.
	// If Dir is empty, the tool is run in the current directory.
	Dir string

	// Build is a list of command-line flags to be passed through to
	// the build system's query tool.
	Build []string

	// Patterns are the file patterns to use for the project root.
	Patterns []string

	// Tests indicates the test files will also be read for the patterns.
	Tests bool

	// Overlay is a mapping from absolute file paths to file contents.
	Overlay map[string][]byte

	// Parallel indicates that the packages should be loaded
	// in parallel when possible based on dependencies and
	// remodelled in parallel.
	// Otherwise, the packages are loaded one at a time in the order
	// that the packages are defined in the project.
	Parallel bool

	// DisableCacheRead disables reading from cache files.
	DisableCacheRead bool

	// DisableCacheWrite disables writing new cache files.
	DisableCacheWrite bool
}

// Target is the interface for all target languages.
type Target interface {

	// Language is the name of the targeted language, e.g. `typescript`.
	Language() string

	// Aliases is all the names that can be used (including the language name)
	// to select this builder, e.g. {`cs`, `csharp`, `c#`}. These will be matched
	// ignore letter case, e.g. `JS`, `Js`, `jS`, `js` will all match the alias `js`.
	Aliases() []string

	// List populates the project with filenames and
	// packages without parsing files or checking types.
	List(cfg *ListConfig) (*project.Project, error)

	// Build performs a build to the target language.
	Build(cfg *BuildConfig) error

	// TODO: Add methods to Run, Test, Serve, etc to support the tools.
}
