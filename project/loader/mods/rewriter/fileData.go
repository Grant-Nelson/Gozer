package rewriter

import (
	"go/ast"
	"go/token"
	"sort"

	"github.com/Grant-Nelson/Gozer/avail/faults"
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
	// The total number of characters this pos data covers including
	// all tailing white space. The sum of the widths in [lines] should
	// be equal to this total.
	total int

	// The total number of characters this pos data covers, not including
	// any tailing white space. This is the initial size of the token
	// when the file was read. If the token is changed, such as renaming a
	// method, this width and the [lines] helps determine how much tailing
	// whitespace to put after the new token.
	width int

	// lines represents the size of lines that this pos data covers.
	// The first value is the width to add to the current line.
	// All following lines indicate there is a newline at the end of the
	// prior line and then how much of the new line width is used.
	// If there is a newline at the end of a line but the pos data ends
	// at that newline, the following line will be 0.
	lines []int
}

func createFileData(fs *token.File, f *ast.File) *fileData {
	fd := &fileData{
		fs:      fs,
		posData: map[int]*posData{},
	}
	fd.collectPosOrder(f)
	fd.populatePosLines()
	return fd
}

func (fd *fileData) collectPosOrder(f *ast.File) {
	var nodePos []int
	var prior int
	for pt := range artifacts.WalkPos(f, false) {
		if !pt.Pos.IsValid() {
			panic(faults.New(`walking a file returned an invalid position`).
				With(`file base`, fd.fs.Base()).
				With(`file path`, artifacts.FilePath(fd.fs, f)).
				With(`position tuple`, pt.String()))
		}

		pos := int(*pt.Pos)
		if pos < prior {
			panic(faults.New(`walking a file returned positions in wrong order`).
				With(`file base`, fd.fs.Base()).
				With(`file path`, artifacts.FilePath(fd.fs, f)).
				With(`prior position`, prior).
				With(`current position`, pos).
				With(`position tuple`, pt.String()))
		}

		pd, has := fd.posData[pos]
		if !has {
			nodePos = append(nodePos, pos)
			fd.posData[pos] = &posData{width: pt.Width}
		} else {
			if pd.width != 0 && pt.Width != 0 {
				panic(faults.New(`walking a file got duplicate position with more than one non-zero width`).
					With(`file base`, fd.fs.Base()).
					With(`file path`, artifacts.FilePath(fd.fs, f)).
					With(`position`, pos).
					With(`prior width`, pd.width).
					With(`current width`, pt.Width).
					With(`position tuple`, pt.String()))
			}
			pd.width = max(pd.width, pt.Width)
		}

		prior = pos
	}
	fd.posOrder = nodePos
}

func (fd *fileData) getPosData(pos int) *posData {
	pd, has := fd.posData[int(pos)]
	if !has {
		panic(faults.New(`failed to find the position in the pos data`).
			With(`file base`, fd.fs.Base()).
			With(`file name`, fd.fs.Name()).
			With(`position`, pos))
	}
	return pd
}

func (fd *fileData) populatePosLines() {
	pos := fd.posOrder[0]
	startLine := fd.fs.Line(token.Pos(pos))
	for _, next := range fd.posOrder[1:] {
		endLine := fd.fs.Line(token.Pos(next))
		pd := fd.getPosData(pos)
		pd.total = int(next) - int(pos)
		pd.lines = fd.measureLines(pos, next, startLine, endLine)
		pos, startLine = next, endLine
	}

	// Set for the EOF position.
	pd := fd.getPosData(pos)
	pd.total = 0
	pd.lines = []int{0}
}

// measureLines measures the lines between the given pos and the prior pos.
func (fd *fileData) measureLines(pos, next, startLine, endLine int) []int {
	lines := []int{}
	line := int(fd.fs.LineStart(startLine))
	startCol := int(pos) - line
	for i := startLine + 1; i <= endLine; i++ {
		cur := int(fd.fs.LineStart(i))
		lines = append(lines, cur-line-startCol)
		startCol = 0
		line = cur
	}
	return append(lines, int(next)-line-startCol)
}

func (fd *fileData) PosOrder() []int {
	return fd.posOrder
}

func (fd *fileData) Total(pos token.Pos) int {
	return fd.getPosData(int(pos)).total
}

func (fd *fileData) Width(pos token.Pos) int {
	return fd.getPosData(int(pos)).width
}

func (fd *fileData) Lines(pos token.Pos) []int {
	return fd.getPosData(int(pos)).lines
}

func (fd *fileData) Previous(pos token.Pos) token.Pos {
	if prev := sort.SearchInts(fd.posOrder, int(pos)) - 1; prev > 0 {
		return token.Pos(fd.posOrder[prev])
	}
	return token.Pos(fd.posOrder[0])
}

func (fd *fileData) Next(pos token.Pos) token.Pos {
	next := sort.SearchInts(fd.posOrder, int(pos)) + 1
	if max := len(fd.posOrder) - 1; next >= max {
		return token.Pos(fd.posOrder[max])
	}
	return token.Pos(fd.posOrder[next])
}

func (fd *fileData) Tail(pos token.Pos) []int {
	return fd.TailWithWidth(pos, -1)
}

// TailWithWidth gets the tails with the given width,
// this width should be the same or smaller than the
// width of this position.
// If the given width is less than zero, the position's width is used.
func (fd *fileData) TailWithWidth(pos token.Pos, width int) []int {
	pd := fd.getPosData(int(pos))
	if width < 0 {
		width = pd.width
	}
	line, lineCount := 0, len(pd.lines)
	for {
		if line >= lineCount {
			return []int{}
		}
		lw := pd.lines[line]
		if lw >= width {
			break
		}
		width -= lw
		line++
	}

	tail := make([]int, lineCount-line)
	copy(tail, pd.lines[line:])
	tail[0] -= width
	return tail
}
