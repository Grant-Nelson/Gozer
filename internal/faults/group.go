package faults

import "errors"

var defaultErrorLimit = 100

type Group struct {
	limit int
	err   []error
}

func NewGroup(limit int) *Group {
	if limit < 1 {
		limit = defaultErrorLimit
	}
	return &Group{limit: limit}
}

func (g *Group) addErr(err error) bool {
	if g == nil || err == nil {
		return false
	}
	if count := len(g.err); count > 0 && g.err[count-1] == err {
		// Duplicate error, skip it.
		return false
	}
	g.err = append(g.err, err)
	return true
}

func (g *Group) Add(err error) error {
	if g == nil {
		return nil
	}
	g.addErr(err)
	if len(g.err) >= g.limit {
		return g.Wrap()
	}
	return nil
}

func (g *Group) Fatal(err error) error {
	if g == nil {
		return nil
	}
	g.addErr(err)
	return g.Wrap()
}

func (g *Group) Wrap() error {
	if g == nil {
		return nil
	}
	return errors.Join(g.err...)
}
