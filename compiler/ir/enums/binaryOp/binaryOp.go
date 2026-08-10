package binaryOp

type BinaryOp int

const (
	Invalid = BinaryOp(iota)
	Add
	Subtract
	Multiply
	Divide
	Modulo
	BitwiseAnd
	BitwiseOr
	BitwiseXor
	ShiftLeft
	ShiftRight
	AndNot
	AddAssign
	SubtractAssign
	MultiplyAssign
	DivideAssign
	ModuloAssign
	BitwiseAndAssign
	BitwiseOrAssign
	BitwiseXorAssign
	ShiftLeftAssign
	ShiftRightAssign
	AndNotAssign
	Equal
	NotEqual
	LessThan
	LessThanOrEqual
	GreaterThan
	GreaterThanOrEqual
	LogicalAnd
	LogicalOr
	Assign
	Define
)

func (b BinaryOp) Valid() bool {
	return b >= Add && b <= Define
}

func (b BinaryOp) String() string {
	switch b {
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
	case Assign:
		return `=`
	case Define:
		return `:=`
	default:
		return `invalid`
	}
}
