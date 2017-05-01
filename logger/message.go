package common

import (
	"fmt"

	"github.com/grant-nelson/Gozer/common"
)

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

// AddData adds a new peice of additional information to the message.
func (msg *Message) AddData(key string, val ...interface{}) {
	msg.Data.Add(key, val...)
}

// String gets the string for this message.
func (msg *Message) String() string {
	result := fmt.Sprintf("%s: %q", msg.Kind.String(), msg.Text)
	if !msg.Data.Empty() {
		result += ":\n   " + msg.Data.FormatMap("   ")
	}
	return result
}
