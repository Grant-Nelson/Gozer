package unaryOp

import "go/token"

type UnaryOp int

const (
	undefined UnaryOp = iota

	Neg    // -X
	Not    // !X
	Inc    // X++
	Dec    // X--
	Expand // X...
	Addr   // &X
	Deref  // *X
	BitNot // ^X
)

const max = BitNot

var (
	names   []string
	formats []string
	tokens  map[token.Token]UnaryOp
)

func init() {
	nameMap := map[UnaryOp]string{
		undefined: `UndefinedUnaryOp`,
		Neg:       `Negate`,
		Not:       `BooleanNot`,
		Inc:       `Increment`,
		Dec:       `Decrement`,
		Expand:    `Expand`,
		Addr:      `Address`,
		Deref:     `Dereference`,
		BitNot:    `BitwiseNot`,
	}
	formatMap := map[UnaryOp]string{
		undefined: `Undefined(%v)`,
		Neg:       `(-%v)`,
		Not:       `(!%v)`,
		Inc:       `(%v++)`,
		Dec:       `(%v--)`,
		Expand:    `(%v...)`,
		Addr:      `(&%v)`,
		Deref:     `(*%v)`,
		BitNot:    `(^%v)`,
	}
	tokenMap := map[UnaryOp]token.Token{
		undefined: token.ILLEGAL,
		Neg:       token.SUB,
		Not:       token.NOT,
		Inc:       token.INC,
		Dec:       token.DEC,
		Expand:    token.ELLIPSIS,
		Addr:      token.AND,
		Deref:     token.MUL,
		BitNot:    token.XOR,
	}
	names = make([]string, max+1)
	formats = make([]string, max+1)
	tokens = make(map[token.Token]UnaryOp, max+1)
	for i := undefined; i <= max; i++ {
		names[i] = nameMap[i]
		formats[i] = formatMap[i]
		tokens[tokenMap[i]] = i
	}
}

func (op UnaryOp) Valid() bool {
	// undefined is not valid
	return op > undefined && op <= max
}

func (op UnaryOp) String() string {
	if op >= undefined && op <= max {
		return names[op]
	}
	return `UnknownUnaryOp`
}

func (op UnaryOp) Format() string {
	if op >= undefined && op <= max {
		return formats[op]
	}
	return `Unknown(%v)`
}
