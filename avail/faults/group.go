package faults

import (
	"errors"
	"sync"
)

type Group struct {
	limit int
	err   []error
	lock  *sync.Mutex
}

func NewGroup(limit int) *Group {
	return &Group{
		limit: limit,
		lock:  &sync.Mutex{},
	}
}

func (g *Group) addErr(err error) bool {
	if err == nil {
		return false
	}
	if count := len(g.err); count > 0 && g.err[count-1] == err {
		// Duplicate error, skip it.
		return false
	}
	g.err = append(g.err, err)
	return true
}

func (g *Group) wrapErrs() error {
	return errors.Join(g.err...)
}

func (g *Group) Add(err error) error {
	if g == nil {
		return err
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	if g.addErr(err) && g.limit > 0 && len(g.err) >= g.limit {
		return g.wrapErrs()
	}
	return nil
}

func (g *Group) Panic(err error) {
	if err == nil {
		return
	}
	if g == nil {
		panic(err)
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	if err2 := g.Add(err); err2 != nil {
		panic(err2)
	}
}

func (g *Group) Fatal(err error) error {
	if g == nil {
		return err
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	if g.addErr(err) {
		// Set limit to the current length so
		// that all following [Add] calls error too.
		g.limit = len(g.err)
		return g.wrapErrs()
	}
	return nil
}

func (g *Group) Wrap() error {
	if g == nil {
		return nil
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	return g.wrapErrs()
}
