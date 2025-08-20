package augmenter

import (
	"errors"
	"go/ast"
	"go/token"
	"strings"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/file"
)

type directives struct {
	add        bool
	delete     bool
	replace    bool
	replaceSig bool
	rename     string
	renameRecv string
	ignore     bool
	none       bool
}

func readDirectives(comments []*ast.Comment, pkgPath string, pos token.Position, errGroup *faults.Group) (dv *directives, err error) {
	defer faults.Recover(&err)
	const directiveGroup = `gozer`
	dr := &directiveReader{
		dm:       file.Directives(comments, directiveGroup),
		dv:       &directives{none: true},
		pkgPath:  pkgPath,
		pos:      pos,
		errGroup: errGroup,
	}
	dr.readAdd()
	dr.readDelete()
	dr.readReplace()
	dr.readReplaceSig()
	dr.readRename()
	dr.readRenameRecv()
	dr.readIgnore()
	return dr.dv, nil
}

type directiveReader struct {
	dm       map[string][]string
	dv       *directives
	pkgPath  string
	pos      token.Position
	errGroup *faults.Group
}

func (dr *directiveReader) postErr(err error) error {
	return dr.errGroup.Add(faults.From(err).
		With(`package`, dr.pkgPath).
		With(`pos`, dr.pos))
}

func (dr *directiveReader) check(test bool, errMsg error) {
	if test {
		if err := dr.postErr(errMsg); err != nil {
			panic(err)
		}
	}
}

var (
	ErrAugAddOnlyOne  = errors.New(`may only define gozer:add once per construct`)
	ErrAugAddWithArgs = errors.New(`may not define gozer:add with arguments`)
)

func (dr *directiveReader) readAdd() {
	const directiveAdd = `add`
	if args, ok := dr.dm[directiveAdd]; ok {
		dr.check(len(args) != 1, ErrAugAddOnlyOne)
		dr.check(len(args[0]) > 0, ErrAugAddWithArgs)
		dr.dv.add = true
		dr.dv.none = false
		delete(dr.dm, directiveAdd)
	}
}

var (
	ErrAugDeleteOnlyOne  = errors.New(`may only define gozer:delete once per construct`)
	ErrAugDeleteWithArgs = errors.New(`may not define gozer:delete with arguments`)
	ErrAugAddAndDelete   = errors.New(`may not define gozer:add with gozer:delete`)
)

func (dr *directiveReader) readDelete() {
	const directiveDelete = `delete`
	if args, ok := dr.dm[directiveDelete]; ok {
		dr.check(len(args) != 1, ErrAugDeleteOnlyOne)
		dr.check(len(args[0]) > 0, ErrAugDeleteWithArgs)
		dr.check(dr.dv.add, ErrAugAddAndDelete)
		dr.dv.delete = true
		dr.dv.none = false
		delete(dr.dm, directiveDelete)
	}
}

var (
	ErrAugReplaceOnlyOne   = errors.New(`may only define gozer:replace once per construct`)
	ErrAugReplaceWithArgs  = errors.New(`may not define gozer:replace with arguments`)
	ErrAugAddAndReplace    = errors.New(`may not define gozer:add with gozer:replace`)
	ErrAugDeleteAndReplace = errors.New(`may not define gozer:delete with gozer:replace`)
)

func (dr *directiveReader) readReplace() {
	const directiveReplace = `replace`
	if args, ok := dr.dm[directiveReplace]; ok {
		dr.check(len(args) != 1, ErrAugReplaceOnlyOne)
		dr.check(len(args[0]) > 0, ErrAugReplaceWithArgs)
		dr.check(dr.dv.add, ErrAugAddAndReplace)
		dr.check(dr.dv.delete, ErrAugDeleteAndReplace)
		dr.dv.replace = true
		dr.dv.none = false
		delete(dr.dm, directiveReplace)
	}
}

var (
	ErrAugReplaceSigOnlyOne   = errors.New(`may only define gozer:replaceSig once per construct`)
	ErrAugReplaceSigWithArgs  = errors.New(`may not define gozer:replaceSig with arguments`)
	ErrAugAddAndReplaceSig    = errors.New(`may not define gozer:add with gozer:replaceSig`)
	ErrAugDeleteAndReplaceSig = errors.New(`may not define gozer:delete with gozer:replaceSig`)
)

func (dr *directiveReader) readReplaceSig() {
	const directiveReplaceSig = `replaceSig`
	if args, ok := dr.dm[directiveReplaceSig]; ok {
		dr.check(len(args) != 1, ErrAugReplaceSigOnlyOne)
		dr.check(len(args[0]) > 0, ErrAugReplaceSigWithArgs)
		dr.check(dr.dv.add, ErrAugAddAndReplaceSig)
		dr.check(dr.dv.delete, ErrAugDeleteAndReplaceSig)
		dr.dv.replace = true
		dr.dv.replaceSig = true
		dr.dv.none = false
		delete(dr.dm, directiveReplaceSig)
	}
}

var (
	ErrAugRenameOnlyOne     = errors.New(`may only define gozer:rename once per construct`)
	ErrAugRenameWithoutArgs = errors.New(`may not define gozer:rename without an argument`)
	ErrAugRenameArgOnlyOne  = errors.New(`may not define gozer:rename more than one argument`)
	ErrAugAddAndRename      = errors.New(`may not define gozer:add with gozer:rename`)
	ErrAugDeleteAndRename   = errors.New(`may not define gozer:delete with gozer:rename`)
)

func (dr *directiveReader) readRename() {
	const directiveRename = `rename`
	if args, ok := dr.dm[directiveRename]; ok {
		dr.check(len(args) != 1, ErrAugRenameOnlyOne)
		dr.check(len(args[0]) <= 0, ErrAugRenameWithoutArgs)
		parts := strings.Fields(args[0])
		dr.check(len(parts[0]) != 1, ErrAugRenameArgOnlyOne)
		dr.check(dr.dv.add, ErrAugAddAndRename)
		dr.check(dr.dv.delete, ErrAugDeleteAndRename)
		dr.dv.replace = true
		dr.dv.rename = parts[0]
		dr.dv.none = false
		delete(dr.dm, directiveRename)
	}
}

var (
	ErrAugRenameRecvOnlyOne     = errors.New(`may only define gozer:renameRecv once per construct`)
	ErrAugRenameRecvWithoutArgs = errors.New(`may not define gozer:renameRecv without an argument`)
	ErrAugRenameRecvArgOnlyOne  = errors.New(`may not define gozer:renameRecv more than one argument`)
	ErrAugAddAndRenameRecv      = errors.New(`may not define gozer:add with gozer:renameRecv`)
	ErrAugDeleteAndRenameRecv   = errors.New(`may not define gozer:delete with gozer:renameRecv`)
)

func (dr *directiveReader) readRenameRecv() {
	const directiveRenameRecv = `renameRecv`
	if args, ok := dr.dm[directiveRenameRecv]; ok {
		dr.check(len(args) != 1, ErrAugRenameOnlyOne)
		dr.check(len(args[0]) <= 0, ErrAugRenameWithoutArgs)
		parts := strings.Fields(args[0])
		dr.check(len(parts[0]) != 1, ErrAugRenameArgOnlyOne)
		dr.check(dr.dv.add, ErrAugAddAndRename)
		dr.check(dr.dv.delete, ErrAugDeleteAndRename)
		dr.dv.replace = true
		dr.dv.renameRecv = parts[0]
		dr.dv.none = false
		delete(dr.dm, directiveRenameRecv)
	}
}

var (
	ErrAugIgnoreOnlyOne  = errors.New(`may only define gozer:renameRecv once per construct`)
	ErrAugIgnoreWithArgs = errors.New(`may not define gozer:renameRecv with arguments`)
	ErrAugIgnoreNotAlone = errors.New(`may not define gozer:ignore with other directives`)
)

func (dr *directiveReader) readIgnore() {
	const directiveIgnore = `ignore`
	if args, ok := dr.dm[directiveIgnore]; ok {
		dr.check(len(args) != 1, ErrAugIgnoreOnlyOne)
		dr.check(len(args[0]) <= 0, ErrAugIgnoreWithArgs)
		dr.check(dr.dv.add || dr.dv.delete || dr.dv.replace || dr.dv.replaceSig ||
			len(dr.dv.rename) > 0 || len(dr.dv.renameRecv) > 0, ErrAugIgnoreNotAlone)
		dr.dv.ignore = true
		dr.dv.none = false
		delete(dr.dm, directiveIgnore)
	}
}
