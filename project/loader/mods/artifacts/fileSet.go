package artifacts

import (
	"fmt"
	"go/token"
	"slices"
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

// TODO: Future: Add method/type name field like in source maps

func (fs *FileSet) Widths(pos token.Pos) (total int, lines []int) {
	fsFile := fs.fileSet.File(pos)
	if fsFile == nil {
		panic(fmt.Errorf(`failed to find fileSet file when getting widths for %d`, pos))
	}

	next := fs.findNext(pos)
	if next <= pos {
		// This occurs when pos is the eof.
		return 0, []int{0}
	}

	total = int(next) - int(pos)
	startLine := fsFile.Line(pos)
	endLine := fsFile.Line(next)

	lines = []int{}
	line := int(fsFile.LineStart(startLine))
	startCol := int(pos) - line
	for i := startLine + 1; i <= endLine; i++ {
		cur := int(fsFile.LineStart(i))
		lines = append(lines, cur-line-startCol)
		startCol = 0
		line = cur
	}
	lines = append(lines, int(next)-line-startCol)

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
	return token.Pos(fsFile.Base() + fsFile.Size())
}

func (fs *FileSet) RegisterFile(f *File) {
	if f.Empty() {
		return
	}

	filePos := int(f.File.FileStart)
	fsFile := fs.fileSet.File(f.File.FileStart)
	if fsFile == nil {
		panic(fmt.Errorf(`failed to find fileSet file when registering %d (%s)`, filePos, f.File.Name.String()))
	}
	if _, exists := fs.nodePos[filePos]; exists {
		panic(fmt.Errorf(`file for %d (%s) already registered`, filePos, f.File.Name.String()))
	}

	var nPos []int
	for _, off := range WalkPos(f.File) {
		nPos = append(nPos, int(*off))
	}
	sort.Ints(nPos)
	nPos = slices.Compact(nPos)
	fs.nodePos[filePos] = nPos
}
