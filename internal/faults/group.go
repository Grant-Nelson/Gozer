package faults

import "errors"

var ErrorLimit = 100

type Group struct {
	err []error
}

func NewGroup() *Group {
	return &Group{}
}

func (g *Group) Add(err error) error {
	if err == nil {
		return nil
	}
	g.err = append(g.err, err)
	if len(g.err) >= ErrorLimit {
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
