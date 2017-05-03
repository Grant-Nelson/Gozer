package msg

import (
	"fmt"
	"io"
	"strings"
)

var _ MessageProcessor = (*LogIO)(nil)

// LogIO is a message processor to write messages to an output.
type LogIO struct {

	// Output is the output to wrtie the logs to.
	// If the output is set to nil then the log is printed to the console.
	Output io.Writer

	// Errors indicates error messages should be written.
	Errors bool

	// Warnings indicates warning messages should be written.
	Warnings bool

	// Info indicates information messages should be written.
	Info bool

	// Debug indicates debugging messages should be written.
	Debug bool
}

// NewLogIO creates a new log output.
func NewLogIO(output io.Writer) *LogIO {
	return &LogIO{
		Output:   output,
		Errors:   true,
		Warnings: true,
		Info:     false,
		Debug:    false,
	}
}

// write prints the given message to the console or output.
func (log *LogIO) write(msg *Message) {
	str := msg.String()
	if log.Output != nil {
		io.WriteString(log.Output, str)
		io.WriteString(log.Output, "\n")
	} else {
		fmt.Println(str)
	}
}

// Process writes the message to the given output if it should be written.
// The given message is returned.
func (log *LogIO) Process(msg *Message) *Message {
	if msg != nil {
		switch msg.Kind {
		case Error:
			if log.Errors {
				log.write(msg)
			}
			break
		case Warning:
			if log.Warnings {
				log.write(msg)
			}
			break
		case Info:
			if log.Info {
				log.write(msg)
			}
			break
		case Debug:
			if log.Debug {
				log.write(msg)
			}
			break
		}
	}
	return msg
}

// String gets the string for this message processor.
func (log *LogIO) String() string {
	parts := []string{}
	if log.Errors {
		parts = append(parts, "Errors")
	}
	if log.Warnings {
		parts = append(parts, "Warnings")
	}
	if log.Info {
		parts = append(parts, "Info")
	}
	if log.Debug {
		parts = append(parts, "Debug")
	}
	return fmt.Sprint("LogIO(", strings.Join(parts, ", "), ")")
}
