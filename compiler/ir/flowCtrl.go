package ir

// FlowCtrl is a flow control node such as a return or branch.
type FlowCtrl interface {
	Node

	// FlowCtrlNode is an empty method used for duck-typing flow control node.
	FlowCtrlNode()
}
