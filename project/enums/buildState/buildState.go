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
	// reading files, parsing, augmenting, and type checking the package.
	// This is filling out the AST part of the package information.
	Loading

	// Loaded indicates that the package has finished [Loading] and is waiting
	// for other packages to finish loading before moving to the next state.
	Loaded

	// Finished indicates the package has finished building.
	Finished
)

func (b BuildState) Valid() bool {
	switch b {
	case Listed, Loading, Loaded, Finished:
		return true
	}
	return false
}

func (b BuildState) String() string {
	switch b {
	case Listed:
		return `Listed`
	case Loading:
		return `Loading`
	case Loaded:
		return `Loaded`
	case Finished:
		return `Finished`
	default:
		return `UnknownState`
	}
}
