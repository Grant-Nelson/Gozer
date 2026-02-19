package builder

import (
	"strings"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/logger"
)

const defaultErrorGroupLimit = 10

type Config struct {

	// Lang is the language to transpile into, e.g. `ts` or `typescript`.
	// This is case insensitive.
	Lang string

	// Logger to log verbose messages with. Has no affect if verbose was false.
	Logger logger.Logger

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
	// in parallel when possible based on dependencies.
	// Otherwise, the packages are loaded one at a time in the order
	// that the packages are defined in the project.
	Parallel bool
}

func Build(cfg *Config) (err error) {
	if cfg.ErrGroup == nil {
		cfg.ErrGroup = faults.NewErrGroup(defaultErrorGroupLimit)
	}
	defer cfg.ErrGroup.Recover(&err)

	switch strings.ToLower(cfg.Lang) {
	case `ts`, `typescript`:
		return buildTS(cfg)
	default:
		return faults.New(`Unsupported target language`).
			With(`target`, cfg.Lang)
	}
}

// Languages returns a list of all supported languages.
func Languages() []string {
	return []string{`typescript`}
}
