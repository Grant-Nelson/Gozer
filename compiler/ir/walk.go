package ir

import (
	"github.com/Grant-Nelson/Gozer/avail/iterator"
	"github.com/Grant-Nelson/Gozer/avail/stack"
)

// Parent is the interface any parent node should implement.
type Parent interface {
	// ChildCount returns the number of child nodes of this node.
	//
	// This value will be equal to or slightly larger than the number
	// of yields from `Children` because this count may include nil children.
	// This is useful to quickly preallocate memory for the children.
	//
	// To get an exact count of non-nil children, count the number of
	// yields from `Children`.
	ChildCount() int

	// Children will call yield for all children to this node in read order.
	// This will not return references to nodes that are not directly
	// child node of this node. This must not yield a nil node.
	//
	// This returns false if yield returns false to exit the yield
	// early. If yield never returns false, this will return true.
	Children(yield func(Node) bool)
}

type NodeConstraint interface {
	Node
	comparable
}

// YieldNode will call the given yield for a non-nil node.
// Returns false if yield returns false for this node, otherwise true.
func YieldNode[T NodeConstraint](n T, yield func(Node) bool) bool {
	var zero T
	return n == zero || yield(n)
}

// YieldSlice will call the given yield for all nodes in the given slice.
// This will not yield nil nodes in the given slice.
// Returns false if yield returns false for any value, otherwise true.
func YieldSlice[T NodeConstraint, S ~[]T](s S, yield func(Node) bool) bool {
	var zero T
	for _, n := range s {
		if n != zero && !yield(n) {
			return false
		}
	}
	return true
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

func WalkNodes(roots ...Node) iterator.Iterator[Node] {
	return Walk(roots...).Select(func(s *WalkStep) Node { return s.Node })
}

func Walk(roots ...Node) iterator.Iterator[*WalkStep] {
	return walkStack(stack.New[Node]().Push(roots...))
}

func walkStack(s stack.Stack[Node]) iterator.Iterator[*WalkStep] {
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
					s.Grow(p.ChildCount())
					s.PushSeq(p.Children)
				}
			}
		}
	}
}
