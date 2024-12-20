package args

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput       = ArgsError{inner: errors.New(`invalid input`)}
	ErrTooFewArgs         = ErrInvalidInput.with(`expected at least one argument containing the command name`)
	ErrNotStructPointer   = ErrInvalidInput.with(`expected a pointer to a struct`)
	ErrNilPointer         = ErrInvalidInput.with(`pointer to a struct may not be nil`)
	ErrWriteOutFailure    = ErrInvalidInput.with(`write failure to out`)
	ErrWriteErrFailure    = ErrInvalidInput.with(`write failure to err`)
	ErrExpectedStructType = ErrInvalidInput.with(`expected a struct type`)
	ErrInvalidHelpTag     = ErrInvalidInput.with(`help tag is only valid on string fields`)
	ErrUnknownTag         = ErrInvalidInput.with(`unknown tag`)
	ErrGroupTagWrongType  = ErrInvalidInput.with(`group tag is only valid on struct or pointer to struct fields`)
	ErrGroupTagRequired   = ErrInvalidInput.with(`group tag cannot be required`)
	ErrInvalidGroupName   = ErrInvalidInput.with(`group name is not valid`)
	ErrGroupAlreadyExists = ErrInvalidInput.with(`group name is already defined`)
	ErrFlagTagWrongType   = ErrInvalidInput.with(`flag tag is only valid on basic type fields`)
	ErrInvalidFlagName    = ErrInvalidInput.with(`flag name is not valid`)
	ErrFlagNameReserved   = ErrInvalidInput.with(`flag name is reserved for help`)
	ErrFlagAlreadyExists  = ErrInvalidInput.with(`flag name is already defined`)
	ErrPosTagWrongType    = ErrInvalidInput.with(`positional argument tag is only valid on basic type fields`)
	ErrInvalidPosName     = ErrInvalidInput.with(`positional argument name is not valid`)
	ErrPosAlreadyExists   = ErrInvalidInput.with(`positional argument name is already defined`)
	ErrVarPosRequired     = ErrInvalidInput.with(`variadic positional argument cannot be required`)
	ErrPosRequiredAfterOp = ErrInvalidInput.with(`all positional arguments after an optional argument must also be optional`)
	ErrVarPosNotLast      = ErrInvalidInput.with(`variadic positional argument must be the last positional argument`)
)

type ArgsError struct {
	inner error
}

func (err ArgsError) Error() string {
	return err.inner.Error()
}

func (err ArgsError) Unwrap() error {
	return errors.Unwrap(err.inner)
}

func (err ArgsError) with(format string, a ...any) ArgsError {
	return ArgsError{inner: fmt.Errorf(`%w: `+format, append([]any{err}, a...)...)}
}
