package remapper

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

// TODO: Convert to be a modifier

// Remap will rewrite the file to a new file set to normalize the file information.
// This is required to be done prior to writing the file so that the file
// will output correctly.
func Remap(f *artifacts.File, finalFileSet *artifacts.FileSet, errGroup *faults.Group) error {
	start := int(f.File.FileStart)
	frm := &fileRemapper{
		f:       f,
		offset:  1,
		prior:   start,
		expNext: start,
	}
	for n, off := range artifacts.WalkPos(f.File) {
		frm.mapPos(n, off)
	}
	frm.finish(fileSet)
	return nil
}

type fileRemapper struct {
	f       *artifacts.File
	offset  int
	edits   []remapperEdit
	prior   int
	expNext int
}

type remapperEdit func(f *token.File)

func (frm *fileRemapper) finish(fileSet *artifacts.FileSet) {
	p := frm.f.FileSet.Position(frm.f.File.FileStart)
	f := fileSet.FileSet().AddFile(p.Filename, 1, frm.offset)
	for _, e := range frm.edits {
		e(f)
	}
	frm.f.FileSet = fileSet
}

func (frm *fileRemapper) mapPos(n ast.Node, off *token.Pos) {
	if !off.IsValid() {
		return
	}

	cur := int(*off)
	fmt.Printf("mapPos: cur: %d, node: %T\n", cur, n)

	if cur == frm.prior {
		fmt.Printf("Duplicate\n\n")
		*off = token.Pos(frm.offset)
		return
	}

	if cur != frm.expNext {
		pos := frm.f.FileSet.Position(*off)
		offset := frm.offset
		fmt.Printf("AddLineColumnInfo: offset: %d, pos: %s:%d:%d\n", offset, pos.Filename, pos.Line, pos.Column)
		frm.edits = append(frm.edits, func(f *token.File) {

			fmt.Printf(">>> AddLineColumnInfo: offset: %d, pos: %s:%d:%d\n", offset, pos.Filename, pos.Line, pos.Column)
			f.AddLine(offset)
			f.AddLineColumnInfo(offset, pos.Filename, pos.Line, pos.Column)
		})
	}

	total, lines := frm.f.FileSet.Widths(*off)
	fmt.Printf("Widths: total: %d, lines: %v\n", total, lines)

	for i, ln := range lines {
		if i > 0 {
			offset := frm.offset
			fmt.Printf("AddLine(%d)\n", offset)
			frm.edits = append(frm.edits, func(f *token.File) {
				f.AddLine(offset)
			})
		}
		frm.offset += ln
	}

	frm.expNext = cur + total
	frm.prior = cur
	*off = token.Pos(frm.offset)
	fmt.Printf("Finish: expNext: %d, prior: %d, set: %d\n\n", frm.expNext, frm.prior, int(*off))
}
