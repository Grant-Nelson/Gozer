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

func WalkNodes(roots ...Node) iterator.Iterator[*WalkStep] {
	s := stack.New[Node]()
	pushAll(s, roots)
	if s.Empty() {
		return iterator.Empty[*WalkStep]()
	}
	step := &WalkStep{}
	return func(yield func(*WalkStep) bool) {
		for !s.Empty() {
			node := s.Pop()
			step.Node = node
			step.Skip = false
			if !yield(step) {
				return
			}
			if !step.Skip {
				pushChildren(s, node)
			}
		}
	}
}

func pushNode(s stack.Stack[Node], node Node) {
	if node != nil {
		s.PushOne(node)
	}
}

func pushAll[T Node, S ~[]T](s stack.Stack[Node], nodes S) {
	for i := len(nodes) - 1; i >= 0; i-- {
		pushNode(s, nodes[i])
	}
}

func pushChildren(s stack.Stack[Node], node Node) {
	switch n := node.(type) {
	case *Block:
		pushAll(s, n.Body)
	case *BlockRef:
		pushAll(s, n.Args)
	case *GotoBlockStmt:
		pushNode(s, n.Block)
	case *FuncCallStmt:
		pushNode(s, n.Follow)
		pushAll(s, n.Args)
		pushNode(s, n.Fun)
	case *DeclStmt:
		pushNode(s, n.Decl)
	case *LabeledStmt:
		pushNode(s, n.Stmt)
		pushNode(s, n.Label)
	case *ExprStmt:
		pushNode(s, n.X)
	case *SendStmt:
		pushNode(s, n.Value)
		pushNode(s, n.Chan)
	case *AssignStmt:
		pushAll(s, n.Rhs)
		pushAll(s, n.Lhs)
	case *GoStmt:
		pushNode(s, n.Call)
	case *DeferStmt:
		pushNode(s, n.Call)
	case *ReturnStmt:
		pushAll(s, n.Results)
	case *BranchStmt:
		pushNode(s, n.Label)
	case *StmtListStmt:
		pushAll(s, n.List)
	case *IfStmt:
		pushAll(s, n.Else)
		pushAll(s, n.Body)
		pushNode(s, n.Cond)
		pushNode(s, n.Init)
	case *SwitchStmt:
		pushAll(s, n.Body)
		pushNode(s, n.Tag)
		pushNode(s, n.Init)
	case *TypeSwitchStmt:
		pushAll(s, n.Body)
		pushNode(s, n.Assign)
		pushNode(s, n.Init)
	case *SelectStmt:
		pushAll(s, n.Body)
	case *ForStmt:
		pushAll(s, n.Body)
		pushNode(s, n.Post)
		pushNode(s, n.Cond)
		pushNode(s, n.Init)
	case *RangeStmt:
		pushAll(s, n.Body)
		pushNode(s, n.X)
		pushNode(s, n.Value)
		pushNode(s, n.Key)
	case *CaseClause:
		pushAll(s, n.Body)
		pushAll(s, n.List)
	case *CommClause:
		pushAll(s, n.Body)
		pushNode(s, n.Comm)
	case ast.Node:
		cs := stack.New[Node]()
		ast.Inspect(n, func(child ast.Node) bool {
			if child != nil {
				cs.Push(child)
			}
			return false
		})
		s.PushSeq(cs.Iterate(), cs.Count())
	default:
		panic(fmt.Errorf(`unexpected node type in walk: %T`, n))
	}
}
