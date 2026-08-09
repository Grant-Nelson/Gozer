package converter

import (
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/binaryOp"
	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/branchKind"
	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/unaryOp"
)

func (c *Converter) FromBranchToken(t token.Token) branchKind.BranchKind {
	switch t {
	case token.BREAK:
		return branchKind.Break
	case token.CONTINUE:
		return branchKind.Continue
	case token.GOTO:
		return branchKind.Goto
	case token.FALLTHROUGH:
		return branchKind.Fallthrough
	default:
		c.addFault(faults.New(`unexpected token for a branch kind`).
			With(`token`, t.String()))
		return branchKind.Invalid
	}
}

func (c *Converter) FromUnaryOp(t token.Token) unaryOp.UnaryOp {
	switch t {
	case token.SUB:
		return unaryOp.Negate
	case token.MUL:
		return unaryOp.Dereference
	case token.AND:
		return unaryOp.Reference
	case token.XOR:
		return unaryOp.BitwiseInvert
	case token.INC:
		return unaryOp.Increment
	case token.DEC:
		return unaryOp.Decrement
	case token.NOT:
		return unaryOp.Not
	default:
		c.addFault(faults.New(`unexpected token for a unary op`).
			With(`token`, t.String()))
		return unaryOp.Invalid
	}
}

func (c *Converter) FromBinaryOp(t token.Token) binaryOp.BinaryOp {
	switch t {
	case token.ADD:
		return binaryOp.Add
	case token.SUB:
		return binaryOp.Subtract
	case token.MUL:
		return binaryOp.Multiply
	case token.QUO:
		return binaryOp.Divide
	case token.REM:
		return binaryOp.Modulo
	case token.ADD_ASSIGN:
		return binaryOp.AddAssign
	case token.SUB_ASSIGN:
		return binaryOp.SubtractAssign
	case token.MUL_ASSIGN:
		return binaryOp.MultiplyAssign
	case token.QUO_ASSIGN:
		return binaryOp.DivideAssign
	case token.REM_ASSIGN:
		return binaryOp.ModuloAssign
	case token.AND:
		return binaryOp.BitwiseAnd
	case token.OR:
		return binaryOp.BitwiseOr
	case token.XOR:
		return binaryOp.BitwiseXor
	case token.SHL:
		return binaryOp.ShiftLeft
	case token.SHR:
		return binaryOp.ShiftRight
	case token.AND_NOT:
		return binaryOp.AndNot
	case token.AND_ASSIGN:
		return binaryOp.BitwiseAndAssign
	case token.OR_ASSIGN:
		return binaryOp.BitwiseOrAssign
	case token.XOR_ASSIGN:
		return binaryOp.BitwiseXorAssign
	case token.SHL_ASSIGN:
		return binaryOp.ShiftLeftAssign
	case token.SHR_ASSIGN:
		return binaryOp.ShiftRightAssign
	case token.AND_NOT_ASSIGN:
		return binaryOp.AndNotAssign
	case token.LAND:
		return binaryOp.LogicalAnd
	case token.LOR:
		return binaryOp.LogicalOr
	case token.EQL:
		return binaryOp.Equal
	case token.NEQ:
		return binaryOp.NotEqual
	case token.LSS:
		return binaryOp.LessThan
	case token.LEQ:
		return binaryOp.LessThanOrEqual
	case token.GTR:
		return binaryOp.GreaterThan
	case token.GEQ:
		return binaryOp.GreaterThanOrEqual
	default:
		c.addFault(faults.New(`unexpected token for a binary op`).
			With(`token`, t.String()))
		return binaryOp.Invalid
	}
}
