package msg

import (
	"fmt"
	"strings"
)

var _ Processor = (*Prepender)(nil)

// Prepender prepends text to any message which is logged.
type Prepender struct {

	// Args is the text to prepend to the message's text.
	Args []interface{}
}

// NewPrepender creates a new text message prepender.
func NewPrepender(args ...interface{}) *Prepender {
	return &Prepender{
		Args: args,
	}
}

// Process will prepend to the given message's text and return the message.
func (ds *Prepender) Process(msg *Message) *Message {
	if (ds != nil) && (msg != nil) {
		text := fmt.Sprint(ds.Args...)
		parts := []string{}
		if len(text) > 0 {
			parts = append(parts, text)
		}
		if len(msg.Text) > 0 {
			parts = append(parts, msg.Text)
		}
		msg.Text = strings.Join(parts, ": ")
	}
	return msg
}

// String gets the string for this message processor.
func (ds *Prepender) String() string {
	if ds == nil {
		return nilStr
	}
	return fmt.Sprint("Prepender(", fmt.Sprint(ds.Args...), ")")
}
