package unaryOp

type UnaryOp int

const (
	Invalid       = UnaryOp(iota)
	Negate        // -
	Dereference   // *
	Reference     // &
	BitwiseInvert // ^
	Increment     // ++
	Decrement     // --
	Not           // !
)

func (u UnaryOp) Valid() bool {
	return u >= Negate && u <= Not
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
	case Increment:
		return `++`
	case Decrement:
		return `--`
	case Not:
		return `!`
	default:
		return `invalid`
	}
}
