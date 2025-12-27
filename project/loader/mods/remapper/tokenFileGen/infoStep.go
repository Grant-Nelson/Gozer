package tokenFileGen

import "go/token"

type infoStep struct {
	offset   int
	filename string
	line     int
	column   int
}

var _ genStep = (*infoStep)(nil)

func newInfoStep(offset int, filename string, line, column int) *infoStep {
	return &infoStep{
		offset:   offset,
		filename: filename,
		line:     line,
		column:   column,
	}
}

func (s *infoStep) Apply(f *token.File) {
	f.AddLineColumnInfo(s.offset, s.filename, s.line, s.column)
}
