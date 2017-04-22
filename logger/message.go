package common

import (
	"fmt"
	"strings"
)

// Message contains an error, warning, info, or debug message for the log.
type Message struct {

	// Kind is the type of message.
	Kind MessageKind

	// Text is the non-optional actual message content.
	Text string

	// Data is any additional information about the message
	// such as line number and file name.
	Data map[string]interface{}
}

// NewMessage creates a new message for a log.
func NewMessage(kind MessageKind, args ...interface{}) Message {
	return Message{
		Kind: kind,
		Text: fmt.Sprint(args...),
		Data: map[string]interface{}{},
	}
}

// String gets the string for this message.
func (e Message) String() string {
	return fmt.Sprint(
	//
	// "Kind: ", e.Kind, "\n",
	// "Text: ", fmt.Sprintf("%q", e.Text), "\n",
	//
	// "File: ", e.File, "\n",
	// "Line: ", e.LineNo, "\n",
	// "Stack: ", strings.Replace(e.Stack, "\n", "\n        ", -1)
	)
}
