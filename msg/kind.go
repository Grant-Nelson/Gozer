package msg

// MessageKind is type for the kind of the message.
type MessageKind int

const (
	// Error indicates the message is an error during transpiling.
	Error MessageKind = iota

	// Warning indicates the message is a warning during transpiling.
	Warning

	// Info indicates the message is information during transpiling.
	Info

	// Debug indicates the message is part of the transpiler.
	Debug
)

// String gets the printable name for a message kind.
func (kind MessageKind) String() string {
	switch kind {
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	case Info:
		return "Info"
	case Debug:
		return "Debug"
	default:
		return "Unknown"
	}
}
