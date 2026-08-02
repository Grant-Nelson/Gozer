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
	t := stack.New[Node]()

	yield := func(n Node) bool {
		if n != nil {
			t.Push(n)
		}
		return true
	}

	finish := func() {
		s.PushSeq(t.Iterate(), t.Count())
	}

	YieldSlice(s, roots)
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
				if p, ok := node.(Parent); ok {

					pushChildren(s, node)
				}
			}
		}
	}
}
