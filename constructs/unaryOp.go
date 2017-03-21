package constructs

var _ Expression = (*BinaryOpExp)(nil)

const (
	// PosOp is the unary operation which performs no change to a value.
	PosOp = "+"

	// NegateOp is the unary operation which negates a value.
	NegateOp = "-"

	// IncrementOp is the unary operation which increments a value.
	IncrementOp = "++"

	// DecrementOp is the unary operation which decrements a value.
	DecrementOp = "--"

	// NotOp is the unary operation which NOTs a value.
	NotOp = "!"

	// DereferanceOp is the unary operation which dereferences the value.
	DereferanceOp = "*"

	// ReferanceOp is the unary operation which references the value.
	ReferanceOp = "&"
)

// UnaryOpExp is an expression which applies an operation to a single value.
type UnaryOpExp struct {

	// Exp is the value to operate on.
	Exp Expression

	// Operand is the operations to apply to the expression.
	Operand string

	// Type is the resulting type after the operation.
	Type Type
}

// UnaryOp creates a new unary operation expression.
func UnaryOp(exp Expression, operand string, t Type) *UnaryOpExp {
	return &UnaryOpExp{
		Exp:     exp,
		Operand: operand,
		Type:    t,
	}
}

// ReturnTypes are the types this expression resolves to.
func (e *UnaryOpExp) ReturnTypes() []Type {
	return []Type{e.Type}
}

// String gets the string for this constuct.
func (e *UnaryOpExp) String() string {
	if e == nil {
		return nilStr
	}
	return e.Operand + ToString(e.Exp)
}
