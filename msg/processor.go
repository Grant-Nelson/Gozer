package msg

// Processor is an interface for processing and modifying a message.
type Processor interface {

	// Process handles and/or modifies the given message
	// and returns the same or new message.
	// To remove/stop the message return nil.
	Process(msg *Message) *Message
}
