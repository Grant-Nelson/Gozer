package common

import (
	"fmt"
)

var _ MessageProcessor = (*Counter)(nil)

// Counter is a message processor to count the number
// of a specific kind of message that has been logged.
type Counter struct {

	// Count is the number of messages of a specific kind logged.
	Count int

	// Kind is the specific kind of messages to count.
	Kind MessageKind
}

// NewCounter creates a new message counter.
func NewCounter(kind MessageKind) *Counter {
	return &Counter{
		Count: 0,
		Kind:  kind,
	}
}

// Process handles and/or modifies the given message
// and returns the same or new message.
func (c *Counter) Process(msg *Message) *Message {
	if (msg != nil) && (c.Kind == msg.Kind) {
		c.Count++
	}
	return msg
}

// String gets the string for the counter.
func (c *Counter) String() string {
	return fmt.Sprint(c.Kind.String(), ": ", c.Count)
}
