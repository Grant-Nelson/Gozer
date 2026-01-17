package posFile

import (
	"go/token"
)

type PosFile struct {
	filename string
	offset   int
	steps    []genStep
	hadLine  bool
}

type genStep interface {
	Apply(base int, f *token.File)
}

func New(filename string) *PosFile {
	return &PosFile{filename: filename}
}

func (f *PosFile) Offset() int { return f.offset }

func (f *PosFile) Add(offset int) {
	if offset > 0 {
		f.offset += offset
		f.hadLine = false
	}
}

func (f *PosFile) HadLine() bool { return f.hadLine }

func (f *PosFile) AddLine() {
	f.steps = append(f.steps, newLineStep(f.offset))
	f.hadLine = true
}

func (f *PosFile) AddInfo(filename string, line, column int) {
	f.steps = append(f.steps, newInfoStep(f.offset, filename, line, column))
	f.hadLine = false
}

func (f *PosFile) Write(fileSet *token.FileSet) int {
	tf := fileSet.AddFile(f.filename, -1, f.offset)
	base := tf.Base()
	for _, step := range f.steps {
		step.Apply(base, tf)
	}
	return base
}
