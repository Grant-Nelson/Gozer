package ir

// Parent is the interface any parent node should implement.
type Parent interface {
	// Children will call yield for all children to this node in read order.
	// This will not return references to nodes that are not directly
	// child node of this node. This may yield a null node.
	//
	// This returns false if yield returns false to exit the yield
	// early. If yield never returns false, this will return true.
	Children(yield func(Node) bool)
}

// YieldSlice will call the given yield for all nodes in the given slice.
func YieldSlice[T Node, S ~[]T](s S, yield func(Node) bool) bool {
	for _, n := range s {
		if !yield(n) {
			return false
		}
	}
	return true
}
