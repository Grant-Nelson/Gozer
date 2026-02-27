package binaryOp

type BinaryOp int

const (
	undefined BinaryOp = iota

	Add // X + Y
	Sub // X - Y
	Mul // X * Y
	Div // X / Y
	Rem // X % Y

	BitOr      // X | Y
	BitAnd     // X & Y
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

	Send     // X <- Y
	Selector // X.Y
	Indexer  // X[Y]
)

const max = Indexer

var (
	names   []string
	formats []string
)

func init() {
	nameMap := map[BinaryOp]string{
		undefined: `UndefinedBinaryOp`,

		Add: `Add`,
		Sub: `Subtract`,
		Mul: `Multiply`,
		Div: `Divide`,
		Rem: `Remainder`,

		BitOr:      `BitwiseOr`,
		BitAnd:     `BitwiseAnd`,
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

		Send:     `Send`,
		Selector: `Selector`,
		Indexer:  `Indexer`,
	}
	formatMap := map[BinaryOp]string{
		undefined: `Undefined(%s, %s)`,

		Add: `(%s + %s)`,
		Sub: `(%s - %s)`,
		Mul: `(%s * %s)`,
		Div: `(%s / %s)`,
		Rem: `(%s %% %s)`,

		BitOr:      `(%s | %s)`,
		BitAnd:     `(%s & %s)`,
		BitXor:     `(%s ^ %s)`,
		BitAndNot:  `(%s &^ %s)`,
		ShiftLeft:  `(%s << %s)`,
		ShiftRight: `(%s >> %s)`,
		LogicalAnd: `(%s && %s)`,
		LogicalOr:  `(%s || %s)`,

		Define:           `(%s := %s)`,
		Assign:           `(%s = %s)`,
		AddAssign:        `(%s += %s)`,
		SubAssign:        `(%s -= %s)`,
		MulAssign:        `(%s *= %s)`,
		DivAssign:        `(%s /= %s)`,
		RemAssign:        `(%s %= %s)`,
		BitAndAssign:     `(%s &= %s)`,
		BitOrAssign:      `(%s |= %s)`,
		BitXorAssign:     `(%s ^= %s)`,
		LogicalAndAssign: `(%s &&= %s)`,
		LogicalOrAssign:  `(%s ||= %s)`,
		ShiftLeftAssign:  `(%s <<= %s)`,
		ShiftRightAssign: `(%s >>= %s)`,
		AndNotAssign:     `(%s &^= %s)`,

		Eq:    `(%s == %s)`,
		Ls:    `(%s < %s)`,
		Gt:    `(%s > %s)`,
		NotEq: `(%s != %s)`,
		LsEq:  `(%s <= %s)`,
		GtEq:  `(%s >= %s)`,

		Send:     `(%s <- %s)`,
		Selector: `(%s.%s)`,
		Indexer:  `(%s[%s])`,
	}

	names = make([]string, max+1)
	formats = make([]string, max+1)
	for i := undefined; i <= max; i++ {
		names[i] = nameMap[i]
		formats[i] = formatMap[i]
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
	return `Unknown(%s, %s)`
}
