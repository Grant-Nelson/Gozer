package msg

import (
	"fmt"
	"strings"

	"github.com/grant-nelson/Gozer/common"
)

// nilStr is the string to use for nil values.
const nilStr = "nil"

// Message contains an error, warning, info, or debug message for the log.
type Message struct {

	// Kind is the type of message.
	Kind MessageKind

	// Text is the non-optional actual message content.
	Text string

	// Data is any additional information about the message
	// such as line number and file name.
	Data common.Map
}

// NewMessage creates a new message for a log.
func NewMessage(kind MessageKind, args ...interface{}) *Message {
	return &Message{
		Kind: kind,
		Text: fmt.Sprint(args...),
		Data: common.NewMap(),
	}
}

// NewError creates a new error message.
func NewError(args ...interface{}) *Message {
	return NewMessage(Error, args...)
}

// NewWarning creates a new warning message.
func NewWarning(args ...interface{}) *Message {
	return NewMessage(Warning, args...)
}

// NewInfo creates a new information message.
func NewInfo(args ...interface{}) *Message {
	return NewMessage(Info, args...)
}

// NewDebug creates a new debug message.
func NewDebug(args ...interface{}) *Message {
	return NewMessage(Debug, args...)
}

// Add adds a new peice of additional information/data to the message.
func (msg *Message) Add(key string, val ...interface{}) *Message {
	if msg != nil {
		msg.Data.Add(key, val...)
	}
	return msg
}

// String gets the string for this message.
func (msg *Message) String() string {
	if msg == nil {
		return nilStr
	}
	result := fmt.Sprintf("%s: %s", msg.Kind.String(), msg.Text)
	if !msg.Data.Empty() {
		result += ":\n  " + msg.Data.FormatMap("  ")
	}
	return result
}

// MessagesToString gets the string for a set of messages.
func MessagesToString(msgs ...*Message) string {
	parts := make([]string, len(msgs))
	for i, msg := range msgs {
		parts[i] = msg.String()
	}
	return strings.Join(parts, "\n")
}
