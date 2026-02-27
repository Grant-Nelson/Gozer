package unaryOp

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
	Recv   // <-X
)

const max = Recv

var (
	names   []string
	formats []string
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
		Recv:      `Receive`,
	}
	formatMap := map[UnaryOp]string{
		undefined: `Undefined(%s)`,
		Neg:       `(-%s)`,
		Not:       `(!%s)`,
		Inc:       `(%s++)`,
		Dec:       `(%s--)`,
		Expand:    `(%s...)`,
		Addr:      `(&%s)`,
		Deref:     `(*%s)`,
		BitNot:    `(^%s)`,
		Recv:      `(<-%s)`,
	}
	names = make([]string, max+1)
	formats = make([]string, max+1)
	for i := undefined; i <= max; i++ {
		names[i] = nameMap[i]
		formats[i] = formatMap[i]
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
	return `Unknown(%s)`
}
