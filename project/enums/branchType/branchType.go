package branchType

import "go/token"

type BranchType int

const (
	undefined BranchType = iota

	Break
	Continue
	Goto
	Fallthrough
)

const max = Fallthrough

var (
	names  []string
	tokens map[token.Token]BranchType
)

func init() {
	nameMap := map[BranchType]string{
		undefined: `UndefinedBranchType`,

		Break:       `break`,
		Continue:    `continue`,
		Goto:        `goto`,
		Fallthrough: `fallthrough`,
	}
	tokenMap := map[BranchType]token.Token{
		undefined: token.ILLEGAL,

		Break:       token.BREAK,
		Continue:    token.CONTINUE,
		Goto:        token.GOTO,
		Fallthrough: token.FALLTHROUGH,
	}
	names = make([]string, max+1)
	tokens = make(map[token.Token]BranchType, max+1)
	for i := undefined; i <= max; i++ {
		names[i] = nameMap[i]
		tokens[tokenMap[i]] = i
	}
}

func (op BranchType) Valid() bool {
	// undefined is not valid
	return op > undefined && op <= max
}

func (op BranchType) String() string {
	if op >= undefined && op <= max {
		return names[op]
	}
	return `UnknownBranchType`
}

func FromToken(t token.Token) BranchType {
	if b, ok := tokens[t]; ok {
		return b
	}
	return undefined
}
