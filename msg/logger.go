package msg

import (
	"github.com/grant-nelson/Gozer/common"
)

// Logger designed for managing messages and logging data.
type Logger struct {

	// msgs is the set of messages.
	msgs []*Message

	// procs is the stack of message processors.
	procs []Processor

	// errCount is the number of errors which have been logged.
	errCount *Counter
}

// NewLogger creates a new logger to a buffer.
func NewLogger() *Logger {
	errCount := NewCounter(Error)
	return &Logger{
		msgs:     []*Message{},
		procs:    []Processor{errCount},
		errCount: errCount,
	}
}

// Messages gets the set of messages.
func (log *Logger) Messages() []*Message {
	if log == nil {
		return []*Message{}
	}
	return log.msgs
}

// Clear removes all the messages from the log
// and resets the error count to zero.
func (log *Logger) Clear() {
	if log != nil {
		log.msgs = []*Message{}
		log.errCount.Count = 0
	}
}

// Process processes all the messages in the log with the given process.
func (log *Logger) Process(proc Processor) {
	if (log != nil) && (proc != nil) {
		log.errCount.Count = 0
		msgs := make([]*Message, 0, len(log.procs))
		for _, msg := range log.msgs {
			msg = proc.Process(msg)
			if msg != nil {
				log.errCount.Process(msg)
				msgs = append(msgs, msg)
			}
		}
		log.msgs = msgs
	}
}

// PushData pushes a new data setter onto the message processor stack.
func (log *Logger) PushData(key string, value ...interface{}) *Logger {
	return log.Push(NewDataSetter(key, value...))
}

// PushPretext pushes a new prepended text onto the message processor stack.
func (log *Logger) PushPretext(args ...interface{}) *Logger {
	return log.Push(NewPrepender(args...))
}

// Push adds a new message processor onto the processor stack.
// The top of the stack processes a message first.
func (log *Logger) Push(proc Processor) *Logger {
	if log != nil {
		log.procs = append(log.procs, proc)
	}
	return log
}

// Pop removes the top message processor from the processor stack.
func (log *Logger) Pop() {
	if log != nil {
		if count := len(log.procs); count > 0 {
			log.procs = log.procs[:count-1]
		}
	}
}

// ErrorCount gets the number of error messages that have been logged.
func (log *Logger) ErrorCount() int {
	if log == nil {
		return 0
	}
	return log.errCount.Count
}

// Add adds a new message of the given kind to the logger.
func (log *Logger) Add(msg *Message) *Message {
	if (log == nil) || (msg == nil) {
		return msg
	}
	for i := len(log.procs) - 1; i >= 0; i-- {
		if proc := log.procs[i]; proc != nil {
			if msg = proc.Process(msg); msg == nil {
				return nil
			}
		}
	}
	log.msgs = append(log.msgs, msg)
	return msg
}

// Error will log an error to the current output or console.
func (log *Logger) Error(args ...interface{}) *Message {
	return log.Add(NewError(args...))
}

// Warning will log a warning to the current output or console.
func (log *Logger) Warning(args ...interface{}) *Message {
	return log.Add(NewWarning(args...))
}

// Info will log information to the current output or console.
func (log *Logger) Info(args ...interface{}) *Message {
	return log.Add(NewInfo(args...))
}

// Debug will log debugging information to the current output or console.
func (log *Logger) Debug(args ...interface{}) *Message {
	msg := NewDebug(args...)
	msg.Add("stack", common.StackTrace(1, 30))
	return log.Add(msg)
}

// String gets the set of messages in the log.
func (log *Logger) String() string {
	return MessagesToString(log.Messages()...)
}

// errorMessage is an error which contains a message which can be thrown.
type errorMessage struct {

	// msg is the internal error message being thrown.
	msg *Message
}

// Error gets the string for the error message.
func (err errorMessage) Error() string {
	return err.msg.String()
}

// ThrowError creates an error. The error message is added
// to the log and then panics with that message.
func (log *Logger) ThrowError(args ...interface{}) {
	msg := log.Error(args...)
	panic(errorMessage{msg: msg})
}

// RethrowError will recover and rethrow a panicked error.
// This must be called via a defer.
// If the recovered error is not a message it will be logged before being rethrown.
func (log *Logger) RethrowError() {
	if r := recover(); r != nil {
		if err, ok := r.(errorMessage); ok {
			panic(err)
		} else {
			log.ThrowError(r)
		}
	}
}

// RecoverError will recover a panicked error.
// This must be called via a defer.
// If the recovered error is not a message it will be logged.
func (log *Logger) RecoverError() {
	if r := recover(); r != nil {
		if _, ok := r.(errorMessage); !ok {
			log.Error(r)
		}
	}
}
