package branchKind

type BranchKind int

const (
	Invalid = BranchKind(iota)
	Break
	Continue
	Goto
	Fallthrough
)

func (b BranchKind) Valid() bool {
	return b >= Break && b <= Fallthrough
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
