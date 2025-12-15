package artifacts

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
)

type FileSet struct {
	fileSet *token.FileSet
	nodePos map[int][]int
}

func NewFileSet() *FileSet {
	return &FileSet{
		fileSet: token.NewFileSet(),
		nodePos: map[int][]int{},
	}
}

func (fs *FileSet) FileSet() *token.FileSet {
	return fs.fileSet
}

// Position gets the position information for the given position offset.
func (fs *FileSet) Position(pos token.Pos) token.Position {
	return fs.fileSet.Position(pos)
}

func (fs *FileSet) Widths(pos token.Pos) (total int, lines []int) {
	fsFile := fs.fileSet.File(pos)
	if fsFile == nil {
		panic(fmt.Errorf(`failed to find fileSet file when getting widths for %d`, pos))
	}

	next := fs.findNext(pos)
	if next == pos {
		panic(fmt.Errorf(`no width found for %d`, pos))
	}

	total = int(next) - int(pos) - 1
	startLine := fsFile.Line(pos)
	endLine := fsFile.Line(next)

	lines = []int{}
	line := int(fsFile.LineStart(startLine))
	startCol := int(pos) - line
	for i := startLine + 1; i <= endLine; i++ {
		line = int(fsFile.LineStart(i))
		lines = append(lines, int(line)-startCol)
		startCol = 0
	}
	lines = append(lines, int(next)-line)

	return total, lines
}

func (fs *FileSet) findNext(pos token.Pos) token.Pos {
	fsFile := fs.fileSet.File(pos)
	nPos, exists := fs.nodePos[fsFile.Base()]
	if !exists {
		panic(fmt.Errorf(`failed to find fileSet file for %d`, int(pos)))
	}

	next := sort.SearchInts(nPos, int(pos)) + 1
	if max := len(nPos) - 1; next <= max {
		return token.Pos(nPos[next])
	}
	return token.Pos(fsFile.Base() + fsFile.Size() - 1)
}

func (fs *FileSet) registerFile(f *ast.File) {
	filePos := int(f.FileStart)
	fsFile := fs.fileSet.File(f.FileStart)
	if fsFile == nil {
		panic(fmt.Errorf(`failed to find fileSet file when registering %d (%s)`, filePos, f.Name.String()))
	}
	if _, exists := fs.nodePos[filePos]; exists {
		panic(fmt.Errorf(`file for %d (%s) already registered`, filePos, f.Name.String()))
	}

	var prev int
	var nPos []int
	walkPos(f, func(n ast.Node, off *token.Pos) {
		cur := int(*off)
		if prev >= cur {
			panic(fmt.Errorf(`node position wasn't in expected order: prev=%d, cur=%d, node=%#v`, prev, cur, n))
		}
		nPos = append(nPos, cur)
		prev = cur
	})

	fs.nodePos[filePos] = nPos
}
