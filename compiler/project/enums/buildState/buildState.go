package buildState

// BuildState is the status of a package during a build.
type BuildState int

const (

	// Listed is the initial state.
	//
	// Listed indicates that the packages have been created and
	// the list of file names per package has been determined.
	// No syntax nor types have been determined and no augmentation
	// has occurred when the build is in this state.
	Listed BuildState = iota

	// Loading indicates that the package is in the process of checking caches,
	// reading files, parsing, augmenting, type checking, and creating the
	// IR in the Go AST shape for the package.
	Loading

	// Loaded indicates that the package has finished [Loading] and is waiting
	// for other packages to finish loading before moving to the next state.
	Loaded

	// Remodelling indicates that the package is being converted from the IR
	// in the Go AST shape into the IR in the shape of the target language.
	Remodelling

	// Remodelled indicates that the package has finished [Remodelling]
	// and has the IR ready for writing the target language.
	Remodelled

	// Writing indicates that the package is being written to the target language.
	Writing

	// Finished indicates the package has finished [Writing] and
	// therefore the package has finished building.
	Finished
)

func (b BuildState) Valid() bool {
	return b >= Listed && b <= Finished
}

func (b BuildState) String() string {
	switch b {
	case Listed:
		return `Listed`
	case Loading:
		return `Loading`
	case Loaded:
		return `Loaded`
	case Remodelling:
		return `Remodelling`
	case Remodelled:
		return `Remodelled`
	case Finished:
		return `Finished`
	case Writing:
		return `Writing`
	default:
		return `UnknownState`
	}
}
