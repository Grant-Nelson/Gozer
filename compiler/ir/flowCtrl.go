package ir

// FlowCtrl is a flow control node such as a return or branch.
type FlowCtrl interface {
	Node

	// FlowCtrlNode is an empty method used for duck-typing flow control node.
	FlowCtrlNode()
}

// IsFlowControlStatement determines if the given node is a FlowCtrl node,
// such as a return or branch.
func IsFlowCtrl(n Node) bool { return Is[FlowCtrl](n) }
