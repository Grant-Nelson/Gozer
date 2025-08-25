package directives

import (
	"errors"
	"go/ast"
	"go/token"
	"strings"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/file"
)

const (
	directiveGroup       = `gozer`
	directiveAdd         = `add`
	directiveDelete      = `delete`
	directiveDeleteAll   = `deleteAll`
	directiveReplace     = `replace`
	directiveReplaceSig  = `replaceSig`
	directiveRename      = `rename`
	directiveReplaceRecv = `replaceRecv`
	directiveIgnore      = `ignore`
)

// Directives contains the gozer directives to indicate what to do
// with code from the augmenter in relation to the original code.
type Directives struct {
	add         bool
	delete      bool
	deleteAll   bool
	replace     bool
	replaceSig  bool
	rename      string
	replaceRecv string
	ignore      bool
}

// Add indicates the code is being added and doesn't exist in the original.
// This will be true for only `gozer:add`.
func (d *Directives) Add() bool { return d.add }

// Delete indicates the code exists in the original and needs to be deleted.
// This will be true for `gozer:delete` and `gozer:deleteAll`.
func (d *Directives) Delete() bool { return d.delete }

// DeleteAll indicates the type exists in the original and needs to be
// deleted along with any function that uses that type for a receiver.
// This will be true for only `gozer:deleteAll`.
func (d *Directives) DeleteAll() bool { return d.deleteAll }

// Replace indicates the code exists in the original and needs to be
// replaced with code being added.
// This will be true for `gozer:replace`, `gozer:replaceSig`,
// `gozer:rename`, and `gozer:replaceRecv`.
func (d *Directives) Replace() bool { return d.replace }

// ReplaceSig indicates the function exists in the original and
// needs to have its signature replaced with the one being added.
// This will be true for only `gozer:replaceSig`.
func (d *Directives) ReplaceSig() bool { return d.replaceSig }

// Rename indicates the code exists in the original and needs to
// have its identifier replaced with the given identifier.
// This will be non-empty for `gozer:rename <name>`.

func (d *Directives) Rename() string { return d.rename }

// HasRename indicate that there is a rename directive.
func (d *Directives) HasRename() bool { return len(d.rename) > 0 }

// ReplaceRecv indicates the function exists in the original and needs to
// have the receiver type replaced with the given type.
// This will be non-empty for `gozer:replaceRecv <type>`.
func (d *Directives) ReplaceRecv() string { return d.replaceRecv }

// HasReplaceRecv indicates that there is a replace receiver directive.
func (d *Directives) HasReplaceRecv() bool { return len(d.replaceRecv) > 0 }

// Ignore indicates the code should be ignored.
// This will be true for only `gozer:ignore`.
func (d *Directives) Ignore() bool { return d.ignore }

// None indicates that there were no directives of any kind on the code.
func (d *Directives) None() bool {
	return !(d.add || d.delete || d.replace || d.ignore)
}

// Copy creates a copy of the directive.
func (d *Directives) Copy() *Directives {
	return &Directives{
		add:         d.add,
		delete:      d.delete,
		deleteAll:   d.deleteAll,
		replace:     d.replace,
		replaceSig:  d.replaceSig,
		rename:      d.rename,
		replaceRecv: d.replaceRecv,
		ignore:      d.ignore,
	}
}

// String returns a human-readable version of the directive for debugging.
func (d *Directives) String() string {
	var parts []string
	concat := func(test bool, s string) {
		if test {
			parts = append(parts, s)
		}
	}
	concat(d.add, directiveAdd)
	concat(d.delete, directiveDelete)
	concat(d.deleteAll, directiveDeleteAll)
	concat(d.replace, directiveReplace)
	concat(d.replaceSig, directiveReplaceSig)
	concat(d.HasRename(), directiveRename+`(`+d.rename+`)`)
	concat(d.HasReplaceRecv(), directiveReplaceRecv+`(`+d.replaceRecv+`)`)
	concat(d.ignore, directiveIgnore)
	concat(d.None(), `none`)
	return `[` + strings.Join(parts, `|`) + `]`
}

func (d *Directives) Join(d2 *Directives, pkgPath string, pos token.Position, errGroup *faults.Group) (dv *Directives, err error) {
	defer faults.Recover(&err)
	mod := &directiveMod{
		dv:       d.Copy(),
		pkgPath:  pkgPath,
		pos:      pos,
		errGroup: errGroup,
	}

	//if d.Add && !d2.Add {

	//}

	// TODO: Implement

	return nil, nil
}

// Read will read the given comments and gather the directives in those comments.
func Read(comments []*ast.Comment, pkgPath string, pos token.Position, errGroup *faults.Group) (dv *Directives, err error) {
	defer faults.Recover(&err)
	mod := &directiveMod{
		dm:       file.Directives(comments, directiveGroup),
		dv:       &Directives{},
		pkgPath:  pkgPath,
		pos:      pos,
		errGroup: errGroup,
	}
	mod.readAdd()
	mod.readDelete()
	mod.readDeleteAll()
	mod.readReplace()
	mod.readReplaceSig()
	mod.readRename()
	mod.readReplaceRecv()
	mod.readIgnore()
	return mod.dv, nil
}

type directiveMod struct {
	dm       map[string][]string
	dv       *Directives
	pkgPath  string
	pos      token.Position
	errGroup *faults.Group
}

func (mod *directiveMod) postErr(err error) error {
	return mod.errGroup.Add(faults.From(err).
		With(`package`, mod.pkgPath).
		With(`pos`, mod.pos))
}

func (mod *directiveMod) check(test bool, errMsg error) {
	if test {
		if err := mod.postErr(errMsg); err != nil {
			panic(err)
		}
	}
}

var (
	ErrAugAddOnlyOne    = errors.New(`may only define gozer:add once per construct`)
	ErrAugAddWithArgs   = errors.New(`may not define gozer:add with arguments`)
	ErrAugAddAndDelete  = errors.New(`may not define an add with a delete`)
	ErrAugAddAndReplace = errors.New(`may not define an add with a replace`)
	ErrAugAddAndIgnore  = errors.New(`may not define an add with an ignore`)
)

func (mod *directiveMod) readAdd() {
	if args, ok := mod.dm[directiveAdd]; ok {
		mod.check(len(args) != 1, ErrAugAddOnlyOne)
		mod.check(len(args[0]) > 0, ErrAugAddWithArgs)
		mod.setAdd()
		delete(mod.dm, directiveAdd)
	}
}

func (mod *directiveMod) setAdd() {
	mod.check(mod.dv.delete, ErrAugAddAndDelete)
	mod.check(mod.dv.replace, ErrAugAddAndReplace)
	mod.check(mod.dv.ignore, ErrAugAddAndIgnore)
	mod.dv.add = true
}

var (
	ErrAugDeleteOnlyOne    = errors.New(`may only define gozer:delete once per construct`)
	ErrAugDeleteWithArgs   = errors.New(`may not define gozer:delete with arguments`)
	ErrAugDeleteAndReplace = errors.New(`may not define a delete with a replace`)
	ErrAugDeleteAndIgnore  = errors.New(`may not define a delete with an ignore`)
)

func (mod *directiveMod) readDelete() {
	if args, ok := mod.dm[directiveDelete]; ok {
		mod.check(len(args) != 1, ErrAugDeleteOnlyOne)
		mod.check(len(args[0]) > 0, ErrAugDeleteWithArgs)
		mod.setDelete()
		delete(mod.dm, directiveDelete)
	}
}

func (mod *directiveMod) setDelete() {
	mod.check(mod.dv.add, ErrAugAddAndDelete)
	mod.check(mod.dv.replace, ErrAugDeleteAndReplace)
	mod.check(mod.dv.ignore, ErrAugDeleteAndIgnore)
	mod.dv.delete = true
}

var (
	ErrAugDeleteAllOnlyOne  = errors.New(`may only define gozer:deleteAll once per construct`)
	ErrAugDeleteAllWithArgs = errors.New(`may not define gozer:deleteAll with arguments`)
)

func (mod *directiveMod) readDeleteAll() {
	if args, ok := mod.dm[directiveDeleteAll]; ok {
		mod.check(len(args) != 1, ErrAugDeleteAllOnlyOne)
		mod.check(len(args[0]) > 0, ErrAugDeleteAllWithArgs)
		mod.setDeleteAll()
		delete(mod.dm, directiveDeleteAll)
	}
}

func (mod *directiveMod) setDeleteAll() {
	mod.setDelete()
	mod.dv.deleteAll = true
}

var (
	ErrAugReplaceOnlyOne   = errors.New(`may only define gozer:replace once per construct`)
	ErrAugReplaceWithArgs  = errors.New(`may not define gozer:replace with arguments`)
	ErrAugReplaceAndIgnore = errors.New(`may not define a replace with an ignore`)
)

func (mod *directiveMod) readReplace() {
	if args, ok := mod.dm[directiveReplace]; ok {
		mod.check(len(args) != 1, ErrAugReplaceOnlyOne)
		mod.check(len(args[0]) > 0, ErrAugReplaceWithArgs)
		mod.setReplace()
		delete(mod.dm, directiveReplace)
	}
}

func (mod *directiveMod) setReplace() {
	mod.check(mod.dv.add, ErrAugAddAndReplace)
	mod.check(mod.dv.delete, ErrAugDeleteAndReplace)
	mod.check(mod.dv.ignore, ErrAugReplaceAndIgnore)
	mod.dv.replace = true
}

var (
	ErrAugReplaceSigOnlyOne  = errors.New(`may only define gozer:replaceSig once per construct`)
	ErrAugReplaceSigWithArgs = errors.New(`may not define gozer:replaceSig with arguments`)
)

func (mod *directiveMod) readReplaceSig() {
	if args, ok := mod.dm[directiveReplaceSig]; ok {
		mod.check(len(args) != 1, ErrAugReplaceSigOnlyOne)
		mod.check(len(args[0]) > 0, ErrAugReplaceSigWithArgs)
		mod.setReplaceSig()
		delete(mod.dm, directiveReplaceSig)
	}
}

func (mod *directiveMod) setReplaceSig() {
	mod.setReplace()
	mod.dv.replaceSig = true
}

var (
	ErrAugRenameOnlyOne     = errors.New(`may only define gozer:rename once per construct`)
	ErrAugRenameWithoutArgs = errors.New(`may not define gozer:rename without an argument`)
	ErrAugRenameArgOnlyOne  = errors.New(`may not define gozer:rename more than one argument`)
)

func (mod *directiveMod) readRename() {
	if args, ok := mod.dm[directiveRename]; ok {
		mod.check(len(args) != 1, ErrAugRenameOnlyOne)
		mod.check(len(args[0]) <= 0, ErrAugRenameWithoutArgs)
		mod.setRename(args[0])
		delete(mod.dm, directiveRename)
	}
}

func (mod *directiveMod) setRename(s string) {
	mod.check(mod.dv.HasRename(), ErrAugRenameOnlyOne)
	parts := strings.Fields(s)
	mod.check(len(parts[0]) != 1, ErrAugRenameArgOnlyOne)
	mod.setReplace()
	mod.dv.rename = parts[0]
}

var (
	ErrAugReplaceRecvOnlyOne     = errors.New(`may only define gozer:replaceRecv once per construct`)
	ErrAugReplaceRecvWithoutArgs = errors.New(`may not define gozer:replaceRecv without an argument`)
	ErrAugReplaceRecvArgOnlyOne  = errors.New(`may not define gozer:replaceRecv more than one argument`)
)

func (mod *directiveMod) readReplaceRecv() {
	if args, ok := mod.dm[directiveReplaceRecv]; ok {
		mod.check(len(args) != 1, ErrAugReplaceOnlyOne)
		mod.check(len(args[0]) <= 0, ErrAugReplaceRecvWithoutArgs)
		mod.setReplaceRecv(args[0])
		delete(mod.dm, directiveReplaceRecv)
	}
}

func (mod *directiveMod) setReplaceRecv(s string) {
	mod.check(mod.dv.HasReplaceRecv(), ErrAugReplaceRecvOnlyOne)
	parts := strings.Fields(s)
	mod.check(len(parts[0]) != 1, ErrAugReplaceRecvArgOnlyOne)
	mod.setReplace()
	mod.dv.replaceRecv = parts[0]
}

var (
	ErrAugIgnoreOnlyOne  = errors.New(`may only define gozer:ignore once per construct`)
	ErrAugIgnoreWithArgs = errors.New(`may not define gozer:ignore with arguments`)
)

func (mod *directiveMod) readIgnore() {
	if args, ok := mod.dm[directiveIgnore]; ok {
		mod.check(len(args) != 1, ErrAugIgnoreOnlyOne)
		mod.check(len(args[0]) <= 0, ErrAugIgnoreWithArgs)
		mod.setIgnore()
		delete(mod.dm, directiveIgnore)
	}
}

func (mod *directiveMod) setIgnore() {
	mod.check(mod.dv.add, ErrAugAddAndIgnore)
	mod.check(mod.dv.delete, ErrAugDeleteAndIgnore)
	mod.check(mod.dv.replace, ErrAugReplaceAndIgnore)
	mod.dv.ignore = true
}
