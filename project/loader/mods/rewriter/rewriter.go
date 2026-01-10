package rewriter

import (
	"fmt"
	"go/ast"
	"go/token"

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

func (rw *Rewriter) Parser(fileSet *token.FileSet, filename string, src any) (*ast.File, error) {
	f, err := rw.parser(fileSet, filename, src)
	if err != nil {
		return nil, err
	}
	if err := rw.recordFile(fileSet, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (rw *Rewriter) recordFile(fileSet *token.FileSet, f *ast.File) error {
	files, found := rw.fSetData[fileSet]
	if !found {
		files = fileSetData{}
		rw.fSetData[fileSet] = files
	}
	fs := fileSet.File(f.FileStart)
	if fs == nil {
		return fmt.Errorf(`failed to find the FileSet file with the FileStart of %d`, f.FileStart)
	}
	if _, found := files[fs.Base()]; found {
		return fmt.Errorf(`failed to record file with the FileStart of %d. An entry already existed`, f.FileStart)
	}
	fd := createFileData(fs, f)
	files[fs.Base()] = fd
	return nil
}

/*
func (fd *fileData) getNodePositions(pos token.Pos) []int {
	fsFile := fs.fileSet.File(pos)
	if fsFile == nil {
		panic(fmt.Errorf(`failed to find fileSet file for %d`, int(pos)))
	}
	nPos, exists := fs.nodePos[fsFile.Base()]
	if !exists || len(nPos) <= 0 {
		panic(fmt.Errorf(`failed to find fileSet node-positions for %d`, int(pos)))
	}
	return nPos
}
*/
