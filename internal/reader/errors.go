package reader

import (
	"strconv"

	"golang.org/x/tools/go/packages"
)

// ProjectErrors is an error that contains multiple project errors.
type ProjectErrors struct {
	Errs []packages.Error
}

// Count is the number of package errors in this error.
func (e *ProjectErrors) Count() int {
	return len(e.Errs)
}

// Error gets the project error text.
func (e *ProjectErrors) Error() string {
	return strconv.Itoa(e.Count()) + ` errors have occurred`
}

// Unwrap gets all the package errors as errors to implement
// the standard multi-wrapped error interface.
func (e *ProjectErrors) Unwrap() []error {
	errs := make([]error, len(e.Errs))
	for i, err := range e.Errs {
		errs[i] = err
	}
	return errs
}
