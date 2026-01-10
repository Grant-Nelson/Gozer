package rewriter

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"
	"sort"

	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

// TODO: Calculate widths beyond size of a node to get the whitespace
// between the nodes, then store that so that we can change the identifiers
// (and identifier lengths) and still be able to remap.

type fileData struct {
	fs *token.File

	// posOrder is all the positions used in the file including floating comments.
	// Positions are sorted into ascending order, i.e. from the top left of
	// the file to the bottom right in read order.
	posOrder []int

	// posData is additional layout information per position in posOrder.
	posData map[int]*posData
}

type posData struct {
	lines []int

	tail []int
}

func createFileData(fs *token.File, f *ast.File) *fileData {
	fd := &fileData{fs: fs}
	fd.collectPosOrder(f)
	fd.populatePosData()

	// TODO: Finish

	return fd
}

func (fd *fileData) collectPosOrder(f *ast.File) {
	var nodePos []int
	var prior int
	for pt := range artifacts.WalkPos(f, false) {
		if !pt.Pos.IsValid() {
			panic(fmt.Errorf(`file for %d (%s) got invalid position for %s`, fd.fs.Base(), f.Name.String(), pt.String()))
		}
		pos := int(*pt.Pos)
		if pos < prior {
			panic(fmt.Errorf(`file for %d (%s) got bad order for prior %d and %d for %s`, fd.fs.Base(), f.Name.String(), prior, pos, pt.String()))
		}
		nodePos = append(nodePos, pos)
		prior = pos
	}

	// The nodes should already be in sorted order, but double check to be safe.
	// Also remove duplicates if any exist since FileStart may be the same as another position.
	sort.Ints(nodePos)
	nodePos = slices.Compact(nodePos)

	fd.posOrder = nodePos
}

func (fd *fileData) FindPrevious(pos token.Pos) token.Pos {
	if prev := sort.SearchInts(fd.posOrder, int(pos)) - 1; prev > 0 {
		return token.Pos(fd.posOrder[prev])
	}
	return token.Pos(fd.posOrder[0])
}

func (fd *fileData) FindNext(pos token.Pos) token.Pos {
	next := sort.SearchInts(fd.posOrder, int(pos)) + 1
	if max := len(fd.posOrder) - 1; next >= max {
		return token.Pos(fd.posOrder[max])
	}
	return token.Pos(fd.posOrder[next])
}

func (fd *fileData) populatePosData() {
	fd.posData = make(map[int]*posData, len(fd.posOrder))
	for _, pos := range fd.posOrder[1:] {
		lines := fd.measureLines(token.Pos(pos))
		fd.posData[pos] = &posData{
			lines: lines,
		}
	}
}

// measureLines measures the lines between the given pos and the prior pos.
//
// TODO: Update this method to make it do less work. The pos, startLine,
// and others could be used in the next call since the given pos should
// be called in order.
func (fd *fileData) measureLines(pos token.Pos) []int {
	next := fd.FindNext(pos)
	if next <= pos {
		// This occurs when pos is the eof.
		return []int{0}
	}

	startLine := fd.fs.Line(pos)
	endLine := fd.fs.Line(next)

	lines := []int{}
	line := int(fd.fs.LineStart(startLine))
	startCol := int(pos) - line
	for i := startLine + 1; i <= endLine; i++ {
		cur := int(fd.fs.LineStart(i))
		lines = append(lines, cur-line-startCol)
		startCol = 0
		line = cur
	}
	lines = append(lines, int(next)-line-startCol)
	return lines
}
