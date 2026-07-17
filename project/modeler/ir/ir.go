package ir

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"

	"github.com/Grant-Nelson/Gozer/avail/iterator"
	"github.com/Grant-Nelson/Gozer/avail/stack"
)

const (
	directiveGroup      = `gozer`
	directiveAtomicFunc = `atomic`
)

type Node interface {
	Pos() token.Pos
}

func csvString[E any, S []E](s S) string {
	const sep = `, `
	elems := make([]string, len(s))
	for i, elem := range s {
		elems[i] = toString(elem)
	}
	return strings.Join(elems, sep)
}

func linesString[E any, S []E](s S) string {
	const indent = `  `
	const nl = "\n"
	elems := make([]string, len(s))
	for i, elem := range s {
		eStr := indent + toString(elem)
		elems[i] = strings.ReplaceAll(eStr, nl, nl+indent)
	}
	return strings.Join(elems, nl)
}

func toString(t any) string {
	switch t := t.(type) {
	case nil:
		return `<nil>`
	case ast.Node:
		return nodeString(t)
	case interface{ String() string }:
		return t.String()
	default:
		return fmt.Sprintf(`%v`, t)
	}
}

func nodeString(n ast.Node) string {
	buf := &bytes.Buffer{}
	fSet := token.NewFileSet()
	if err := format.Node(buf, fSet, n); err != nil {
		panic(err)
	}
	return buf.String()
}

func astPos(n ast.Node) token.Pos {
	if n == nil {
		return token.NoPos
	}
	return n.Pos()
}

// WalkStep will be returned from a walk. One instance is reused
// for each step so don't hold onto the instance.
type WalkStep struct {
	Node Node

	// By setting skip to true, the children of Node will not be walked.
	// This will be defaulted to false for each node. To stop all iteration,
	// return false from the iterator or break/return from the for-range.
	Skip bool
}

func WalkNodes(root Node) iterator.Iterator[*WalkStep] {
	if root == nil {
		return iterator.Empty[*WalkStep]()
	}
	stack := stack.New[Node]()
	stack.Push(root)
	step := &WalkStep{}
	return func(yield func(*WalkStep) bool) {
		for !stack.Empty() {
			node := stack.Pop()
			step.Node = node
			step.Skip = false
			if !yield(step) {
				return
			}
			if !step.Skip {
				pushChildren(stack, node)
			}
		}
	}
}

func pushAll[T Node, S ~[]T](stack stack.Stack[Node], s S) {
	for i := len(s) - 1; i >= 0; i-- {
		stack.PushOne(s[i])
	}
}

// TODO: Finish

func pushChildren(stack stack.Stack[Node], node Node) {
	switch n := node.(type) {
	case *Block:
		pushAll(stack, n.Body)
	case *GotoBlockStmt:
		//Block  *BlockRef
	case *FuncCallStmt:
		//Fun    ast.Expr
		//Args   []ast.Expr
		//Follow *BlockRef
	case *DeclStmt:
		//Decl ast.Decl
	case *LabeledStmt:
		//Label *ast.Ident
		//Stmt  Stmt
	case *ExprStmt:
		//X   ast.Expr
	case *SendStmt:
		//Chan  ast.Expr
		//Value ast.Expr
	case *AssignStmt:
		//Lhs []ast.Expr
		//Tok token.Token
		//Rhs []ast.Expr
	case *GoStmt:
		//Call *ast.CallExpr
	case *DeferStmt:
		//Call *ast.CallExpr
	case *ReturnStmt:
		//Results []ast.Expr
	case *BranchStmt:
		//Tok   token.Token
		//Label *ast.Ident
	case *StmtListStmt:
		//List []Stmt
	case *IfStmt:
		//Init Stmt
		//Cond ast.Expr
		//Body []Stmt
		//Else []Stmt
	case *SwitchStmt:
		//Init Stmt
		//Tag  ast.Expr
		//Body []*CaseClause
	case *TypeSwitchStmt:
		//Init   Stmt
		//Assign Stmt
		//Body   []*CaseClause
	case *SelectStmt:
		//Body []*CommClause
	case *ForStmt:
		//Init Stmt
		//Cond ast.Expr
		//Post Stmt
		//Body []Stmt
	case *RangeStmt:
		//Key   ast.Expr
		//Value ast.Expr
		//Tok   token.Token
		//X     ast.Expr
		//Body  []Stmt
	case *CaseClause:
		//List []ast.Expr
		//Body []Stmt
	case *CommClause:
		//Comm Stmt
		//Body []Stmt
	}
}

/*
// visitExprIdents walks an expression and records identifier reads
// into use (unless they have already been defined locally).
func visitExprIdents(e ast.Expr, info *types.Info, use, def objectSet) {
	if e == nil {
		return
	}
	ast.Inspect(e, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		obj := info.Uses[id]
		if obj == nil {
			return true
		}
		if _, isVar := obj.(*types.Var); !isVar {
			return true
		}
		if !def.has(obj) {
			use.add(obj)
		}
		return true
	})
}
*/
