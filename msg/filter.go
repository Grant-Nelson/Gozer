package msg

import "fmt"

var _ Processor = (*Filter)(nil)

// Filter removes messages from the log.
type Filter struct {

	// Kind of message to remove from the logs.
	Kind MessageKind
}

// NewFilter creates a new message filter message processor.
func NewFilter(kind MessageKind) *Filter {
	return &Filter{
		Kind: kind,
	}
}

// Process removes any message of the specific message kind.
// If nil is returned the message is filtered.
func (f *Filter) Process(msg *Message) *Message {
	if (f != nil) && (msg != nil) && (msg.Kind == f.Kind) {
		return nil
	}
	return msg
}

// String gets the string for this message processor.
func (f *Filter) String() string {
	if f == nil {
		return nilStr
	}
	return fmt.Sprint("Filter(", f.Kind, ")")
}
