package blocks

import "go/ast"

type Func struct {
	Ast    *ast.FuncDecl
	Blocks []*Block
}
