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

	// TODO: Calculate widths beyond size of a node to get the whitespace
	// between the nodes, then store that so that we can change the identifiers
	// (and identifier lengths) and still be able to remap.
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

	next := fs.FindNext(pos)
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

func (fs *FileSet) getNodePositions(pos token.Pos) []int {
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

func (fs *FileSet) FindPrevious(pos token.Pos) token.Pos {
	nPos := fs.getNodePositions(pos)
	if prev := sort.SearchInts(nPos, int(pos)) - 1; prev > 0 {
		return token.Pos(nPos[prev])
	}
	return token.Pos(nPos[0])
}

func (fs *FileSet) FindNext(pos token.Pos) token.Pos {
	nPos := fs.getNodePositions(pos)
	next := sort.SearchInts(nPos, int(pos)) + 1
	if max := len(nPos) - 1; next >= max {
		return token.Pos(nPos[max])
	}
	return token.Pos(nPos[next])
}

// RegisterFile adds the extra file information to the FileSet for the given
// file. The file must be an unmodified or remapped file so that all the
// positions are part of the same file entry in the FileSet.
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
	var prior int
	for pt := range WalkPos(f.File, false) {
		if !pt.Pos.IsValid() {
			panic(fmt.Errorf(`file for %d (%s) got invalid position for %s`, filePos, f.File.Name.String(), pt.String()))
		}

		pos := int(*pt.Pos)
		if pos < prior {
			panic(fmt.Errorf(`file for %d (%s) got bad order for prior %d and %d for %s`, filePos, f.File.Name.String(), prior, pos, pt.String()))
		}
		nPos = append(nPos, pos)
		prior = pos
	}
	// Should already be sorted, but sort anyway
	sort.Ints(nPos)
	nPos = slices.Compact(nPos)
	fs.nodePos[filePos] = nPos
}
