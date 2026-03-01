package binaryOp

import "go/token"

type BinaryOp int

const (
	undefined BinaryOp = iota

	Add // X + Y
	Sub // X - Y
	Mul // X * Y
	Div // X / Y
	Rem // X % Y

	BitAnd     // X & Y
	BitOr      // X | Y
	BitXor     // X ^ Y
	BitAndNot  // X &^ Y
	ShiftLeft  // X << Y
	ShiftRight // X >> Y
	LogicalAnd // X && Y
	LogicalOr  // X || Y

	Define           // X := Y
	Assign           // X = Y
	AddAssign        // X += Y
	SubAssign        // X -= Y
	MulAssign        // X *= Y
	DivAssign        // X /= Y
	RemAssign        // X %= Y
	BitAndAssign     // X &= Y
	BitOrAssign      // X |= Y
	BitXorAssign     // X ^= Y
	LogicalAndAssign // X &&= Y
	LogicalOrAssign  // X ||= Y
	ShiftLeftAssign  // X <<= Y
	ShiftRightAssign // X >>= Y
	AndNotAssign     // X &^= Y

	Eq    // X == Y
	Ls    // X < Y
	Gt    // X > Y
	NotEq // X != Y
	LsEq  // X <= Y
	GtEq  // X >= Y

	Selector // X.Y
	Indexer  // X[Y]
)

const max = Indexer

var (
	names   []string
	formats []string
	tokens  map[token.Token]BinaryOp
)

func init() {
	nameMap := map[BinaryOp]string{
		undefined: `UndefinedBinaryOp`,

		Add: `Add`,
		Sub: `Subtract`,
		Mul: `Multiply`,
		Div: `Divide`,
		Rem: `Remainder`,

		BitAnd:     `BitwiseAnd`,
		BitOr:      `BitwiseOr`,
		BitXor:     `BitwiseXor`,
		BitAndNot:  `BitwiseAndNot`,
		ShiftLeft:  `ShiftLeft`,
		ShiftRight: `ShiftRight`,
		LogicalAnd: `LogicalAnd`,
		LogicalOr:  `LogicalOr`,

		Define:           `Define`,
		Assign:           `Assign`,
		AddAssign:        `AddAssign`,
		SubAssign:        `SubAssign`,
		MulAssign:        `MulAssign`,
		DivAssign:        `DivAssign`,
		RemAssign:        `RemAssign`,
		BitAndAssign:     `BitwiseAndAssign`,
		BitOrAssign:      `BitwiseOrAssign`,
		BitXorAssign:     `BitwiseXorAssign`,
		LogicalAndAssign: `LogicalAndAssign`,
		LogicalOrAssign:  `LogicalOrAssign`,
		ShiftLeftAssign:  `ShiftLeftAssign`,
		ShiftRightAssign: `ShiftRightAssign`,
		AndNotAssign:     `AndNotAssign`,

		Eq:    `Equal`,
		Ls:    `LessThan`,
		Gt:    `GreaterThan`,
		NotEq: `NotEqual`,
		LsEq:  `LessThanOrEqual`,
		GtEq:  `GreaterThanOrEqual`,

		Selector: `Selector`,
		Indexer:  `Indexer`,
	}
	formatMap := map[BinaryOp]string{
		undefined: `Undefined(%v, %v)`,

		Add: `(%v + %v)`,
		Sub: `(%v - %v)`,
		Mul: `(%v * %v)`,
		Div: `(%v / %v)`,
		Rem: `(%v %% %v)`,

		BitAnd:     `(%v & %v)`,
		BitOr:      `(%v | %v)`,
		BitXor:     `(%v ^ %v)`,
		BitAndNot:  `(%v &^ %v)`,
		ShiftLeft:  `(%v << %v)`,
		ShiftRight: `(%v >> %v)`,
		LogicalAnd: `(%v && %v)`,
		LogicalOr:  `(%v || %v)`,

		Define:           `(%v := %v)`,
		Assign:           `(%v = %v)`,
		AddAssign:        `(%v += %v)`,
		SubAssign:        `(%v -= %v)`,
		MulAssign:        `(%v *= %v)`,
		DivAssign:        `(%v /= %v)`,
		RemAssign:        `(%v %= %v)`,
		BitAndAssign:     `(%v &= %v)`,
		BitOrAssign:      `(%v |= %v)`,
		BitXorAssign:     `(%v ^= %v)`,
		LogicalAndAssign: `(%v &&= %v)`,
		LogicalOrAssign:  `(%v ||= %v)`,
		ShiftLeftAssign:  `(%v <<= %v)`,
		ShiftRightAssign: `(%v >>= %v)`,
		AndNotAssign:     `(%v &^= %v)`,

		Eq:    `(%v == %v)`,
		Ls:    `(%v < %v)`,
		Gt:    `(%v > %v)`,
		NotEq: `(%v != %v)`,
		LsEq:  `(%v <= %v)`,
		GtEq:  `(%v >= %v)`,

		Selector: `(%v.%v)`,
		Indexer:  `(%v[%v])`,
	}
	tokenMap := map[BinaryOp]token.Token{
		undefined: token.ILLEGAL,

		Add: token.ADD,
		Sub: token.SUB,
		Mul: token.MUL,
		Div: token.QUO,
		Rem: token.REM,

		BitAnd:     token.ADD,
		BitOr:      token.OR,
		BitXor:     token.XOR,
		BitAndNot:  token.AND_NOT,
		ShiftLeft:  token.SHL,
		ShiftRight: token.SHR,
		LogicalAnd: token.LAND,
		LogicalOr:  token.LOR,

		Define:       token.DEFINE,
		Assign:       token.ASSIGN,
		AddAssign:    token.ADD_ASSIGN,
		SubAssign:    token.SUB_ASSIGN,
		MulAssign:    token.MUL_ASSIGN,
		DivAssign:    token.QUO_ASSIGN,
		RemAssign:    token.REM_ASSIGN,
		BitAndAssign: token.AND_ASSIGN,
		BitOrAssign:  token.OR_ASSIGN,
		BitXorAssign: token.XOR_ASSIGN,
		// LogicalAndAssign: No go equivalent token,
		// LogicalOrAssign:  No go equivalent token,
		ShiftLeftAssign:  token.SHL_ASSIGN,
		ShiftRightAssign: token.SHR_ASSIGN,
		AndNotAssign:     token.AND_NOT_ASSIGN,

		Eq:    token.EQL,
		Ls:    token.LSS,
		Gt:    token.GTR,
		NotEq: token.NEQ,
		LsEq:  token.LEQ,
		GtEq:  token.GEQ,

		// Selector: No go equivalent token,
		// Indexer:  No go equivalent token,
	}
	names = make([]string, max+1)
	formats = make([]string, max+1)
	tokens = make(map[token.Token]BinaryOp, max+1)
	for i := undefined; i <= max; i++ {
		names[i] = nameMap[i]
		formats[i] = formatMap[i]
		tokens[tokenMap[i]] = i
	}
}

func (op BinaryOp) Valid() bool {
	// undefined is not valid
	return op > undefined && op <= max
}

func (op BinaryOp) String() string {
	if op >= undefined && op <= max {
		return names[op]
	}
	return `UnknownBinaryOp`
}

func (op BinaryOp) Format() string {
	if op >= undefined && op < max {
		return formats[op]
	}
	return `Unknown(%v, %v)`
}

func FromToken(t token.Token) BinaryOp {
	if b, ok := tokens[t]; ok {
		return b
	}
	return undefined
}
