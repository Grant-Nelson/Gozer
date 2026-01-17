package rewriter

import (
	"go/ast"
	"go/token"
	"slices"
	"strings"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/rewriter/posFile"
)

type remapper struct {
	fileSet  *token.FileSet
	f        *ast.File
	fileData *fileData
	floaters map[int][]*ast.CommentGroup
}

func (rm *remapper) perform(targetFileSet *token.FileSet) {
	if rm.checkForUnmodifiedFile() {
		return
	}

	rm.calculateFloaters()
	rm.remapFile(targetFileSet)
}

// checkForUnmodifiedFile quickly checks if the file is consistent with
// an unmodified file. If it isn't modified then most of remap can be skipped.
// If fileSet and targetFileSet are the same, then the remap can be completely
// skipped otherwise a quick remap can be used to map the unmodified file
// to the new file set.
func (rm *remapper) checkForUnmodifiedFile() bool {
	// TODO: If a fast simple way can be found, then implement, otherwise remove.
	return false
}

// calculateFloaters sets the [floaters] field with a map that contains
// the sorted comments that need to be added after the node with the position
// used as the key.
func (rm *remapper) calculateFloaters() {
	posSorted, posMap := getSolidPosOrder(rm.f)
	floating := getFloatingComments(rm.f, posMap)
	rm.floaters = getFloaterAttachments(posSorted, floating)
}

// getSolidPosOrder collects all the current positions of all the nodes
// in the file while skipping file comments. We skip the file comments since
// the floating ones in those will have to be positioned specially for them
// to be close to the correct location in the modified file as they were
// in the original files that were used to create the modified file.
func getSolidPosOrder(f *ast.File) ([]int, map[int]bool) {
	posSorted := []int{}
	posMap := map[int]bool{}
	for pt := range artifacts.WalkPos(f, artifacts.SkipFileComments) {
		pos := int(*pt.Pos)
		posSorted = append(posSorted, pos)
		posMap[pos] = true
	}
	slices.Sort(posSorted)
	posSorted = slices.Compact(posSorted)
	return posSorted, posMap
}

// getFloatingComments collects all comments that are floating in the file and
// not tied to any specific node's Comment nor Doc fields.
func getFloatingComments(f *ast.File, posMap map[int]bool) map[int]*ast.CommentGroup {
	floating := map[int]*ast.CommentGroup{}
	for _, c := range f.Comments {
		pos := int(c.Pos())
		if !posMap[pos] {
			floating[pos] = c
		}
	}
	return floating
}

// getFloaterAttachments finds the previous positions to each of the floating
// text with the possibility that the previous position is from another floating
// comment. The result is the list of floating comment positions in the order
// they should be added keyed by the position the comments follow after.
func getFloaterAttachments(posSorted []int, floaters map[int]*ast.CommentGroup) map[int][]*ast.CommentGroup {
	follows := map[int][]*ast.CommentGroup{}
	for _, floater := range floaters {
		index, exists := slices.BinarySearch(posSorted, int(floater.Pos()))
		if exists {
			panic(faults.New(`the floating comment appeared in the position sorted set`).
				With(`comment pos`, floater))
		}
		index--
		if index < 0 {
			panic(faults.New(`the floating comment is before any position in the sorted set`).
				With(`comment pos`, floater).
				With(`index`, index))
		}
		prior := posSorted[index]

		comments := append(follows[prior], floater)
		slices.SortFunc(comments, commentCmp)
		follows[prior] = comments
	}
	return follows
}

// commentCmp is a comparer for two comments that compares them by
// the position of the front of the comment.
func commentCmp(a, b *ast.CommentGroup) int {
	return int(a.Pos()) - int(b.Pos())
}

func (rm *remapper) remapFile(targetFileSet *token.FileSet) {
	posFile := posFile.New(rm.fileData.Name())

	// TODO: Need to ensure lines between decls
	for pt := range artifacts.WalkPos(rm.f, artifacts.SkipFileComments) {
		rm.placePos(posFile, pt)
		rm.placeFloaters(posFile, pt)
	}

	base := posFile.Write(targetFileSet)
	for pt := range artifacts.WalkPos(rm.f, artifacts.SkipFileComments) {
		*pt.Pos = token.Pos(base + int(*pt.Pos))
	}
}

func (rm *remapper) placeFloaters(posFile *posFile.PosFile, prior artifacts.PosTuple) {
	cgs, found := rm.floaters[int(*prior.Pos)]
	if !found {
		return
	}
	for _, cg := range cgs {
		for pt := range artifacts.WalkPos(cg) {
			rm.placePos(posFile, pt)
		}
	}
}

func (rm *remapper) placePos(posFile *posFile.PosFile, pt artifacts.PosTuple) {
	cur := posFile.Offset()

	rest := false
	for line := range strings.Lines(pt.Text) {
		if rest {
			posFile.AddLine()
		}
		posFile.Add(len(line))
		rest = true
	}

	// TODO: Figure out a way to ensure there is space if needed to separate identifiers,
	// insert a ".", "]", etc.

	tail := rm.fileData.TailWithWidth(*pt.Pos, pt.Width)
	for i, width := range tail {
		if i > 0 {
			posFile.AddLine()
		}
		posFile.Add(width)
	}

	*pt.Pos = token.Pos(cur)
}
