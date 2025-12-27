package tokenFileGen

import "go/token"

type TokenFileGen struct {
	filename string
	base     int
	offset   int
	steps    []genStep
}

type genStep interface {
	Apply(f *token.File)
}

func New(filename string, base int) *TokenFileGen {
	return &TokenFileGen{
		filename: filename,
		base:     base,
		offset:   base,
	}
}

func (f *TokenFileGen) Current() token.Pos {
	return token.Pos(f.offset)
}

func (f *TokenFileGen) Add(offset int) {
	f.offset += offset
}

func (f *TokenFileGen) AddLine() {
	f.steps = append(f.steps, newLineStep(f.offset))
}

func (f *TokenFileGen) AddInfo(filename string, line, column int) {
	f.steps = append(f.steps, newInfoStep(f.offset, filename, line, column))
}

func (f *TokenFileGen) Write(outFileSet *token.FileSet) {
	tf := outFileSet.AddFile(f.filename, f.base, f.offset-f.base)
	for _, step := range f.steps {
		step.Apply(tf)
	}
}
