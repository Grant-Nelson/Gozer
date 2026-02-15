package faults

import (
	"errors"
	"sync"
)

// Group collects multiple errors into an error containing all of those errors.
// The error group is designed to work concurrently.
type Group struct {
	limit int
	err   []error
	lock  *sync.Mutex
}

// NewGroup creates a new error group with the given limit.
//
// If the limit is greater than zero then the error will be created
// when that number of errors have been added.
// If the limit is zero or negative, then the error count is unlimited.
func NewGroup(limit int) *Group {
	return &Group{
		limit: limit,
		lock:  &sync.Mutex{},
	}
}

type wrappedError interface{ Unwrap() []error }

func (g *Group) addErr(err error) bool {
	if g.isPriorWrap(err) || g.isLastError(err) {
		// Nil or duplicate error(s), skip it.
		return false
	}

	g.err = append(g.err, err)
	return true
}

func (g *Group) isPriorWrap(err error) bool {
	if wErr, ok := err.(wrappedError); ok {
		inner := wErr.Unwrap()
		count := len(inner)
		return count > 0 && count <= len(g.err) && inner[0] == g.err[0]
	}
	return false
}

func (g *Group) isLastError(err error) bool {
	count := len(g.err)
	return count > 0 && g.err[count-1] == err
}

func (g *Group) wrapErrs() error {
	if len(g.err) <= 0 {
		return nil
	}
	return errors.Join(g.err...)
}

func (g *Group) atLimit() bool {
	return g.limit > 0 && len(g.err) >= g.limit
}

// Empty indicates if there has been one or more errors collected.
func (g *Group) Empty() bool {
	if g == nil {
		return true
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	return len(g.err) <= 0
}

// Add adds the given error to this group.
//
// If the given error is nil, the last error, or a prior wrap of this group,
// then it is skipped.
// If this error is added and the limit is reached, this will return
// the wrapped error for this group. Otherwise, this returns nil.
func (g *Group) Add(err error) error {
	if err == nil || g == nil {
		return err
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	if g.addErr(err) && g.atLimit() {
		return g.wrapErrs()
	}
	return nil
}

// Panic adds the given error to this group. If the group has reached
// its limit this will panic the wrapped error for this group.
//
// If the given error is nil, the last error, or a prior wrap of this group,
// then it is skipped.
// If the error is added but the limit is not reached, then this will not
// panic and simply add the error.
func (g *Group) Panic(err error) {
	if wErr := g.Add(err); wErr != nil {
		panic(wErr)
	}
}

// Fatal adds the given error to this group. This will always return
// a wrapped group as if the limit was reached. It will also set the limit
// so that any added error after a fatal will still return a wrapped group.
//
// If the given error is nil, then it is skipped.
// If the given error is the last error, or a prior wrap of this group,
// then the error will not be added but the errors will still be wrapped
// and returned.
// This will return nil if the group is empty and the given error is empty,
// but will still set the limit so the next non-nil error will return
// a wrapped group.
func (g *Group) Fatal(err error) error {
	if err == nil || g == nil {
		return err
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	g.addErr(err)
	// Set limit to the current length so
	// that all following [Add] calls error too.
	g.limit = max(1, len(g.err))
	return g.wrapErrs()
}

// Wrap returns a wrapped error with all the errors in this group
// or nil if there were no errors in this group.
func (g *Group) Wrap() error {
	if g == nil {
		return nil
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	return g.wrapErrs()
}
