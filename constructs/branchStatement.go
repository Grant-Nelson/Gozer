package constructs

var _ Statement = (*BranchStat)(nil)

// BranchStat is a branch statement for break or continue.
type BranchStat struct {

	// Break indicates if the branch is a break or a continue.
	Break bool
}

// Branch creates a new branch statement for break or continue.
func Branch(Break bool) *BranchStat {
	return &BranchStat{
		Break: Break,
	}
}

// String gets the string for this constuct.
func (s *BranchStat) String() string {
	if s == nil {
		return nilStr
	}
	if s.Break {
		return "break"
	}
	return "continue"
}
