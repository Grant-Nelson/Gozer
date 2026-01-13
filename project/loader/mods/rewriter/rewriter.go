package rewriter

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/parser"
)

type Rewriter struct {
	parser   parser.Parser
	fSetData map[*token.FileSet]fileSetData
}

type fileSetData map[int]*fileData

func New(p parser.Parser) *Rewriter {
	if p == nil {
		p = parser.Default
	}
	return &Rewriter{
		parser:   p,
		fSetData: map[*token.FileSet]fileSetData{},
	}
}

func (rw *Rewriter) Parser(fileSet *token.FileSet, filename string, src any) (f *ast.File, err error) {
	defer faults.Recover(&err)

	f, err = rw.parser(fileSet, filename, src)
	if err != nil {
		return nil, err
	}

	rw.recordFile(fileSet, f)
	return f, nil
}

func (rw *Rewriter) recordFile(fileSet *token.FileSet, f *ast.File) {
	files, found := rw.fSetData[fileSet]
	if !found {
		files = fileSetData{}
		rw.fSetData[fileSet] = files
	}

	fs := fileSet.File(f.FileStart)
	if fs == nil {
		panic(faults.New(`failed to find the token.File in the token.FileSet with the file start`).
			With(`file start`, f.FileStart))
	}

	if _, found := files[fs.Base()]; found {
		panic(faults.New(`failed to record file since an entry already existed`).
			With(`file start`, f.FileStart).
			With(`file base`, fs.Base()).
			With(`file name`, fs.Name()))
	}

	fd := createFileData(fs, f)
	files[fs.Base()] = fd
}

func (rw *Rewriter) getFileData(fileSet *token.FileSet, f *ast.File) *fileData {
	files, found := rw.fSetData[fileSet]
	if !found {
		panic(faults.New(`failed to find the file set data file with the file start`).
			With(`file start`, f.FileStart))
	}

	fs := fileSet.File(f.FileStart)
	if fs == nil {
		panic(faults.New(`failed to find the token.File in the token.FileSet with the file start`).
			With(`file start`, f.FileStart))
	}

	fd, found := files[fs.Base()]
	if !found {
		panic(faults.New(`failed to find file data`).
			With(`file start`, f.FileStart).
			With(`file base`, fs.Base()).
			With(`file name`, fs.Name()))
	}

	return fd
}
