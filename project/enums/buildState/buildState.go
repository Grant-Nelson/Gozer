package buildState

// BuildState is the status of a package during a build.
type BuildState int

const (
	Listed BuildState = iota
)
