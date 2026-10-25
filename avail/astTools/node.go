package astTools

import (
	"errors"
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/iterator"
)

type NodeIteratorValue struct {
	// Node is the node that is being opened and nil when a node is closed.
	Node ast.Node

	// Skip will be set to false when returned via the iterator.
	// Set to true to skip all the children of the node being opened.
	// There is no effect when set to true when a node is closing.
	Skip bool
}

type NodeWithStackIteratorValue struct {
	// Stack is the current stack of parent nodes including the one being opened
	// but excluding the one being closed. Do not modify the slice.
	Stack []ast.Node

	// Node is the node that is being opened or closed.
	Node ast.Node

	// Closing indicates this node has already been called when opened and
	// it is now being closed.
	Closing bool

	// Skip will be set to false when returned via the iterator.
	// Set to true to skip all the children of the node being opened.
	// There is no effect when set to true when a node is closing.
	Skip bool
}

var errEndInspect = errors.New(`end inspect`)

func newNodeIteratorValue(node ast.Node) *NodeIteratorValue {
	return &NodeIteratorValue{Node: node}
}

func newNodeWithStackIteratorValue(stack []ast.Node, node ast.Node, closing bool) *NodeWithStackIteratorValue {
	return &NodeWithStackIteratorValue{Stack: stack, Node: node, Closing: closing}
}

// Nodes iterates through all the nodes in the file.
// It moves depth first and returns a nil Node when a node is being
// closed and all of its children have been iterated.
// No nil Node is returned for a skipped node.
func Nodes(node ast.Node) iterator.Iterator[*NodeIteratorValue] {
	return func(yield func(v *NodeIteratorValue) bool) {
		defer func() {
			if r := recover(); r != nil && r != errEndInspect {
				panic(r)
			}
		}()

		ast.Inspect(node, func(n ast.Node) bool {
			v := newNodeIteratorValue(n)
			if !yield(v) {
				panic(errEndInspect)
			}
			return !v.Skip
		})
	}
}

// NodesWithStack iterates through all the nodes in the file.
// It moves depth first and returns a nil Node when a node is being
// closed and all of its children have been iterated.
func NodesWithStack(node ast.Node) iterator.Iterator[*NodeWithStackIteratorValue] {
	return func(yield func(v *NodeWithStackIteratorValue) bool) {
		defer func() {
			if r := recover(); r != nil && r != errEndInspect {
				panic(r)
			}
		}()

		stack := []ast.Node{}
		push := func(n ast.Node) {
			stack = append(stack, n)
		}
		pop := func() ast.Node {
			var close ast.Node
			if high := len(stack) - 1; high >= 0 {
				close = stack[high]
				stack[high] = nil
				stack = stack[:high]
			}
			return close
		}

		ast.Inspect(node, func(n ast.Node) bool {
			if n == nil {
				close := pop()
				v := newNodeWithStackIteratorValue(stack, close, true)
				if !yield(v) {
					panic(errEndInspect)
				}
				return true
			}

			push(n)
			v := newNodeWithStackIteratorValue(stack, n, false)
			if !yield(v) {
				panic(errEndInspect)
			}
			if v.Skip {
				// the n == nil path for popping the stack is skipped by ast.Inspect
				// when `false` is returned, so pop the value from the stack now.
				_ = pop()
				return false
			}
			return true
		})
	}
}
