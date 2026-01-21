package rewriter

import (
	"go/ast"
	"go/token"
	"slices"
	"strings"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/rewriter/posFile"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/rewriter/walkPos"
)

type remapper struct {
	fileSet  *token.FileSet
	f        *ast.File
	fileData *fileData
	floaters map[int][]*ast.CommentGroup
	pf       *posFile.PosFile
}

func (rm *remapper) perform(targetFileSet *token.FileSet) {
	if rm.checkForUnmodifiedFile() {
		// TODO: If unmodified but targetFileSet is different,
		// do a quick copy over to new fileSet.
		return
	}
	rm.calculateFloaters()
	rm.remapFile()
	rm.outputPosFile(targetFileSet)
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
	posSorted, posMap := getSolidPosOrder(rm.fileData.fs, rm.f)
	floating := getFloatingComments(rm.f, posMap)
	rm.floaters = getFloaterAttachments(posSorted, floating)
}

// getSolidPosOrder collects all the current positions of all the nodes
// in the file while skipping file comments. We skip the file comments since
// the floating ones in those will have to be positioned specially for them
// to be close to the correct location in the modified file as they were
// in the original files that were used to create the modified file.
func getSolidPosOrder(fs *token.File, f *ast.File) ([]int, map[int]bool) {
	posSorted := []int{}
	posMap := map[int]bool{}
	for pt := range walkPos.WalkPos(fs, f, walkPos.SkipFileComments) {
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

func (rm *remapper) remapFile() {
	rm.pf = posFile.New(rm.fileData.Name())

	// TODO: Need to ensure lines between decls
	// TODO: Write line info directives for decls and specs

	for pt := range walkPos.WalkPos(rm.fileData.fs, rm.f, walkPos.SkipFileComments) {
		rm.placePos(pt)
		rm.placeFloaters(pt)
	}
}

func (rm *remapper) outputPosFile(targetFileSet *token.FileSet) {
	base := rm.pf.Write(targetFileSet)
	for pt := range walkPos.WalkPos(rm.fileData.fs, rm.f, walkPos.SkipPseudoPos) {
		*pt.Pos = token.Pos(base + int(*pt.Pos))
	}
}

func (rm *remapper) placeFloaters(prior walkPos.PosTuple) {
	cgs, found := rm.floaters[int(*prior.Pos)]
	if !found {
		return
	}
	for _, cg := range cgs {
		for pt := range walkPos.WalkPos(rm.fileData.fs, cg) {
			rm.placePos(pt)
		}
	}
}

func (rm *remapper) placePos(pt walkPos.PosTuple) {
	cur := rm.pf.Offset()

	rm.placeText(pt.Text)

	tail := rm.fileData.TailWithWidth(*pt.Pos, pt.Width)
	rm.placeWhitespace(tail)

	*pt.Pos = token.Pos(cur)
}

func (rm *remapper) placeText(text string) {
	rest := false
	for line := range strings.Lines(text) {
		if rest {
			rm.pf.AddLine()
		}
		rm.pf.Add(len(line))
		rest = true
	}
}

func (rm *remapper) placeWhitespace(lines []int) {
	for i, width := range lines {
		if i > 0 {
			rm.pf.AddLine()
		}
		rm.pf.Add(width)
	}
}
