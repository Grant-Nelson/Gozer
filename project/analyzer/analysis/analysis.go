package analysis

type Analysis interface {
}

type (
	FuncAnalysis struct {

		// IsBlocking indicates that the function contains blocking calls
		// such as a send, receive, function call, sleep, or lock.
		// For function calls, if the call to a function that is blocking,
		// then this function is also blocking.
		IsBlocking bool

		// IsSimple indicates that the function does not contain function calls,
		// sends, receives, for-loops, jump backs, etc.
		IsSimple bool
	}
)
