package tokenFileGen

import "go/token"

type lineStep struct {
	offset int
}

var _ genStep = (*lineStep)(nil)

func newLineStep(offset int) *lineStep {
	return &lineStep{
		offset: offset,
	}
}

func (s *lineStep) Apply(f *token.File) {
	f.AddLine(s.offset)
}
