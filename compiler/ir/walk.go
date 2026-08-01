package ir

import (
	"github.com/Grant-Nelson/Gozer/avail/iterator"
	"github.com/Grant-Nelson/Gozer/avail/stack"
)

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

func pushChildren(s stack.Stack[Node], node Node) {
	switch n := node.(type) {
	case *GotoBlockStmt:
		pushNode(s, n.Block)
	case *LabeledStmt:
		pushNode(s, n.Stmt)
		pushNode(s, n.Label)
	case *SendStmt:
		pushNode(s, n.Value)
		pushNode(s, n.Chan)
	case *GoStmt:
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
	case *RangeStmt:
		pushAll(s, n.Body)
		pushNode(s, n.X)
		pushNode(s, n.Value)
		pushNode(s, n.Key)
	}
}
