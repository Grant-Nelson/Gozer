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

func (g *Group) Add(err error) error {
	if err == nil {
		return nil
	}
	g.err = append(g.err, err)
	if len(g.err) >= g.limit {
		return g.Wrap()
	}
	return nil
}

func (g *Group) Fatal(err error) error {
	if err == nil {
		return nil
	}
	g.err = append(g.err, err)
	return g.Wrap()
}

func (g *Group) Wrap() error {
	return errors.Join(g.err...)
}
