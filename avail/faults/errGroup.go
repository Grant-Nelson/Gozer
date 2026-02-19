package faults

import (
	"fmt"
	"strings"
	"sync"
)

// ErrGroup collects multiple errors into an error containing all of those errors.
// The error group is designed to work concurrently.
type ErrGroup struct {
	limit     int
	errs      []error
	remainder int
	lock      *sync.Mutex
}

var _ error = (*ErrGroup)(nil)

// NewErrGroup creates a new error group with the given limit.
//
// If the limit is zero or negative, then the error count is unlimited.
func NewErrGroup(limit int, initialErrors ...error) *ErrGroup {
	g := &ErrGroup{
		limit: limit,
		lock:  &sync.Mutex{},
	}
	g.Add(initialErrors...)
	return g
}

// Error gets the error string for this error.
func (g *ErrGroup) Error() string {
	if g == nil {
		return `<nil>`
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	count := len(g.errs)
	switch count {
	case 0:
		return `no errors`
	case 1:
		return g.errs[0].Error()
	default:
		buf := &strings.Builder{}
		for i, err := range g.errs {
			if i > 0 {
				buf.WriteByte('\n')
			}
			fmt.Fprintf(buf, `%d. %s`, i+1, err.Error())
		}
		if g.remainder > 0 {
			fmt.Fprintf(buf, "\n%d. too many errors (%d remaining)", count, g.remainder)
		}
		return buf.String()
	}
}

// Unwrap returns all the errors in this group.
func (g *ErrGroup) Unwrap() []error {
	if g == nil {
		return []error{}
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	max := len(g.errs) - 1
	return g.errs[:max:max]
}

// Len gets the number of errors in this group including any
// discarded because they are over the limit.
func (g *ErrGroup) Len() int {
	if g == nil {
		return 0
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	return len(g.errs) + g.remainder
}

// Full indicates if there has been enough errors that the limit has
// been reached. Any new errors added will be counted and discarded.
func (g *ErrGroup) Full() bool {
	if g == nil {
		return true
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	return g.limit > 0 && len(g.errs) >= g.limit
}

// Empty indicates if there has been one or more errors collected.
func (g *ErrGroup) Empty() bool {
	if g == nil {
		return true
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	return len(g.errs) <= 0
}

// ErrorOrNil will return this error if not empty or nil if empty.
func (g *ErrGroup) ErrorOrNil() error {
	if g.Empty() {
		return nil
	}
	return g
}

func crunch(errs []error) []error {
	count, dest := len(errs), 0
	for src := range count {
		if errs[src] != nil {
			errs[dest], errs[src] = errs[src], errs[dest]
			dest++
		}
	}
	return errs[:dest]
}

func expand(skip *ErrGroup, errs []error) []error {
	errs2 := make([]error, 0, len(errs))
	for _, err := range errs {
		if g, ok := err.(*ErrGroup); ok {
			if g != skip {
				errs2 = append(errs2, g.Unwrap()...)
			}
		} else {
			errs2 = append(errs2, err)
		}
	}
	return errs2
}

// Add adds the given error(s) to this group.
//
// If this group has reached its limit, this will return this group
// as an error, otherwise it will return nil.
func (g *ErrGroup) Add(errs ...error) error {
	errs = expand(g, crunch(errs))
	count := len(errs)
	if count <= 0 {
		return nil
	}
	if g == nil {
		if count == 1 {
			return errs[0]
		}
		return NewErrGroup(-1, errs...)
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	newSize := len(g.errs) + count
	if g.limit <= 0 || newSize <= g.limit {
		g.errs = append(g.errs, errs...)
		return nil
	}

	extra := newSize - g.limit - 1
	g.remainder += extra
	g.errs = append(g.errs, errs[:extra]...)
	return g
}

// Recover handles catching a panic and adding it to the group,
// then setting the given pointer, [pe], to the result error to the result
// of the [Add]. If [pe] is nil then it will not be set.
func (g *ErrGroup) Recover(pe *error) {
	if disableRecovers {
		return
	}
	if r := recover(); r != nil {
		err := g.Add(From(r))
		if pe != nil {
			*pe = err
		}
	}
}
