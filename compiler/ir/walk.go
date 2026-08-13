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

func WalkPackage(pkg *Package) iterator.Iterator[*WalkStep] {
	s := stack.New[Node]()
	s.PushSeq(pkg.Children, 0)
	return walk(s)
}

func WalkNodes(roots ...Node) iterator.Iterator[*WalkStep] {
	s := stack.New[Node]()
	s.Push(roots...)
	return walk(s)
}

func walk(s stack.Stack[Node]) iterator.Iterator[*WalkStep] {
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
					s.PushSeq(p.Children, 0)
				}
			}
		}
	}
}
