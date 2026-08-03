package branchKind

type BranchKind int

const (
	Break = BranchKind(iota)
	Continue
	Goto
	Fallthrough
)

func (b BranchKind) Valid() bool {
	switch b {
	case Break, Continue, Goto, Fallthrough:
		return true
	}
	return false
}

func (b BranchKind) String() string {
	switch b {
	case Break:
		return `break`
	case Continue:
		return `continue`
	case Goto:
		return `goto`
	case Fallthrough:
		return `fallthrough`
	default:
		return `invalid`
	}
}
