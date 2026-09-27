package unaryOp

type UnaryOp int

const (
	Invalid       = UnaryOp(iota)
	Negate        // -x
	Dereference   // *x
	Reference     // &x
	BitwiseInvert // ^x
	Not           // !x
	PreInc        // ++x
	PreDec        // --x
	PostInc       // x++
	PostDec       // x--
)

func (u UnaryOp) Valid() bool {
	return u > Invalid && u <= PostDec
}

// IsAssignment determines if the variable is assigned
// in this binary operation, e.g. `++x`, `x--`.
func (u UnaryOp) IsAssignment() bool {
	switch u {
	case PreInc, PreDec, PostInc, PostDec:
		return true
	default:
		return false
	}
}

func (u UnaryOp) String() string {
	switch u {
	case Negate:
		return `-`
	case Dereference:
		return `*`
	case Reference:
		return `&`
	case BitwiseInvert:
		return `^`
	case Not:
		return `!`
	case PreInc:
		return `++`
	case PreDec:
		return `--`
	case PostInc:
		return `++`
	case PostDec:
		return `--`
	default:
		return `invalid`
	}
}
