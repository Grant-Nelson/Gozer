package opDefs

import (
	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/binaryOp"
	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/unaryOp"
)

// TODO: Determine if this should be kept or thrown out.

type OpDefs struct {
	UnaryOps  map[unaryOp.UnaryOp]string
	BinaryOps map[binaryOp.BinaryOp]string
}

// See: https://www.geeksforgeeks.org/typescript/typescript-operators/

var Int32 = &OpDefs{
	UnaryOps: map[unaryOp.UnaryOp]string{
		unaryOp.Negate:        `-%[1]s`,
		unaryOp.BitwiseInvert: `~^%[1]s`,
		unaryOp.PreInc:        `++%[1]s`,
		unaryOp.PreDec:        `--%[1]s`,
		unaryOp.PostInc:       `%[1]s++`,
		unaryOp.PostDec:       `%[1]s--`,
	},
	BinaryOps: map[binaryOp.BinaryOp]string{
		binaryOp.Assign:             `%[1]s=%[2]s`,
		binaryOp.Add:                `%[1]s+%[2]s`,
		binaryOp.AddAssign:          `%[1]s+=%[2]s`,
		binaryOp.Subtract:           `%[1]s-%[2]s`,
		binaryOp.SubtractAssign:     `%[1]s-=%[2]s`,
		binaryOp.Multiply:           `%[1]s*%[2]s`,
		binaryOp.MultiplyAssign:     `%[1]s*=%[2]s`,
		binaryOp.Divide:             `%[1]s/%[2]s`,
		binaryOp.DivideAssign:       `%[1]s/=%[2]s`,
		binaryOp.Modulo:             `%[1]s%%%[2]s`,
		binaryOp.ModuloAssign:       `%[1]s%=%[2]s`,
		binaryOp.BitwiseAnd:         `%[1]s&%[2]s`,
		binaryOp.BitwiseAndAssign:   `%[1]s&=%[2]s`,
		binaryOp.BitwiseOr:          `%[1]s|%[2]s`,
		binaryOp.BitwiseOrAssign:    `%[1]s|=%[2]s`,
		binaryOp.BitwiseXor:         `%[1]s^%[2]s`,
		binaryOp.BitwiseXorAssign:   `%[1]s^=%[2]s`,
		binaryOp.ShiftLeft:          `%[1]s<<%[2]s`,
		binaryOp.ShiftLeftAssign:    `%[1]s<<=%[2]s`,
		binaryOp.ShiftRight:         `%[1]s>>%[2]s`,
		binaryOp.ShiftRightAssign:   `%[1]s>>=%[2]s`,
		binaryOp.AndNot:             `%[1]s&^%[2]s`,
		binaryOp.AndNotAssign:       `%[1]s&^=%[2]s`,
		binaryOp.NotEqual:           `%[1]s!=%[2]s`,
		binaryOp.Equal:              `%[1]s==%[2]s`,
		binaryOp.LessThan:           `%[1]s<%[2]s`,
		binaryOp.LessThanOrEqual:    `%[1]s<=%[2]s`,
		binaryOp.GreaterThan:        `%[1]s>%[2]s`,
		binaryOp.GreaterThanOrEqual: `%[1]s>=%[2]s`,
		binaryOp.LogicalAnd:         `%[1]s&&%[2]s`,
		binaryOp.LogicalOr:          `%[1]s||%[2]s`,
	},
}

var Uint32 = &OpDefs{
	// TODO: Finish
}

var Int = Int32
var Uint = Uint32
