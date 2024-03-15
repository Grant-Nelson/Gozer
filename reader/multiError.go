package reader

import (
	"strconv"
)

type MultiError struct {
	errs []error
}

func newError(errs []error) error {
	count := len(errs)
	switch {
	case count > 0:
		return &MultiError{errs: errs}
	case count == 1:
		return errs[0]
	default:
		return nil
	}
}

func (me *MultiError) Error() string {
	return strconv.Itoa(len(me.errs)) + `%d errors have occurred`
}

func (me *MultiError) Unwrap() []error {
	return me.errs
}
