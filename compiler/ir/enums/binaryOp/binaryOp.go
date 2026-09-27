package binaryOp

type BinaryOp int

const (
	Invalid = 0

	Assign           = 1  // x  =  y (assign if x is ref, define if y is def)
	Add              = 2  // x  +  y
	AddAssign        = 3  // x  += y, Assign | Add
	Subtract         = 4  // x  -  y
	SubtractAssign   = 5  // x  -= y, Assign | Subtract
	Multiply         = 6  // x  *  y
	MultiplyAssign   = 7  // x  *= y, Assign | Multiply
	Divide           = 8  // x  /  y
	DivideAssign     = 9  // x  /= y, Assign | Divide
	Modulo           = 10 // x  %  y
	ModuloAssign     = 11 // x  %= y. Assign | Modulo
	BitwiseAnd       = 12 // x  &  y
	BitwiseAndAssign = 13 // x  &= y, Assign | BitwiseAnd
	BitwiseOr        = 14 // x  |  y
	BitwiseOrAssign  = 15 // x  |= y, Assign | BitwiseOr
	BitwiseXor       = 16 // x  ^  y
	BitwiseXorAssign = 17 // x  ^= y, Assign | BitwiseXor
	ShiftLeft        = 18 // x <<  y
	ShiftLeftAssign  = 19 // x <<= y, Assign | ShiftLeft
	ShiftRight       = 20 // x >>  y
	ShiftRightAssign = 21 // x >>= y, Assign | ShiftRight
	AndNot           = 22 // x &^  y
	AndNotAssign     = 23 // x &^= y, Assign | AndNot

	NotEqual           = 32 // x != y
	Equal              = 33 // x == y
	LessThan           = 34 // x <  y
	LessThanOrEqual    = 35 // x <= y, Equal | LessThan
	GreaterThan        = 36 // x >  y
	GreaterThanOrEqual = 37 // x >= y, Equal | GreaterThan

	LogicalAnd = 38
	LogicalOr  = 39
)

func (b BinaryOp) Valid() bool {
	return b > Invalid && b <= LogicalOr
}

// IsAssignment determines if the left hand side (lhs) is assigned
// in this binary operation, e.g. `x = y`, `x += y`, `x *= y`.
func (b BinaryOp) IsAssignment() bool {
	switch b {
	case Assign,
		AddAssign,
		SubtractAssign,
		MultiplyAssign,
		DivideAssign,
		ModuloAssign,
		BitwiseAndAssign,
		BitwiseOrAssign,
		BitwiseXorAssign,
		ShiftLeftAssign,
		ShiftRightAssign,
		AndNotAssign:
		return true
	default:
		return false
	}
}

// IsComparator determine if the operation is a comparator,
// e.g. `x == y`, `x < y`, `x >= y`.
func (b BinaryOp) IsComparator() bool {
	switch b {
	case Equal,
		NotEqual,
		LessThan,
		LessThanOrEqual,
		GreaterThan,
		GreaterThanOrEqual:
		return true
	default:
		return false
	}
}

func (b BinaryOp) String() string {
	switch b {
	case Assign:
		return `=`
	case Add:
		return `+`
	case Subtract:
		return `-`
	case Multiply:
		return `*`
	case Divide:
		return `/`
	case Modulo:
		return `%`
	case BitwiseAnd:
		return `&`
	case BitwiseOr:
		return `|`
	case BitwiseXor:
		return `^`
	case ShiftLeft:
		return `<<`
	case ShiftRight:
		return `>>`
	case AndNot:
		return `&^`
	case AddAssign:
		return `+=`
	case SubtractAssign:
		return `-=`
	case MultiplyAssign:
		return `*=`
	case DivideAssign:
		return `/=`
	case ModuloAssign:
		return `%=`
	case BitwiseAndAssign:
		return `&=`
	case BitwiseOrAssign:
		return `|=`
	case BitwiseXorAssign:
		return `^=`
	case ShiftLeftAssign:
		return `<<=`
	case ShiftRightAssign:
		return `>>=`
	case AndNotAssign:
		return `&^=`
	case Equal:
		return `==`
	case NotEqual:
		return `!=`
	case LessThan:
		return `<`
	case LessThanOrEqual:
		return `<=`
	case GreaterThan:
		return `>`
	case GreaterThanOrEqual:
		return `>=`
	case LogicalAnd:
		return `&&`
	case LogicalOr:
		return `||`
	default:
		return `invalid`
	}
}
