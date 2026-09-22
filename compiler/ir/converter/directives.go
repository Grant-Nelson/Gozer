package converter

import (
	"go/ast"
	"strings"

	"github.com/Grant-Nelson/Gozer/avail/assert"
	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
)

// FromCommentGroup reads directives out of the comment group.
//
// The given target is the node that the comment group is associated with.
// The target can be used to determine if a type of directive is allowed or not.
// Most directives (that aren't ignored as part of another process) will be returned
// and there are only a few that will apply directly to the target.
func (c *converter) FromCommentGroup(cg *ast.CommentGroup, target ir.Node) []ir.Directive {
	ds := []ir.Directive{}
	for d := range astTools.DirectivesFromGroup(cg) {
		if dt := c.FromDirective(d, target); dt != nil {
			ds = append(ds, dt)
		}
	}
	return ds
}

func (c *converter) FromDirective(d *ast.Directive, target ir.Node) ir.Directive {
	switch d.Tool {
	case `line`:
		// Ignore
		return nil
	case `go`:
		return c.fromGoDirective(d, target)
	case `gozer`:
		return c.fromGozerDirective(d, target)
	case `extern`:
		// TODO: Look into leveraging or add to ignored directives
		// See https://go.dev/doc/comment#directives
		return nil
	case `export`:
		// TODO: Look into leveraging or add to ignored directives
		// See https://go.dev/doc/comment#directives
		return nil
	default:
		c.addFault(faults.New(`unexpected directive`).
			With(`pos`, c.pos(d.Pos())).
			With(`tool`, d.Tool).
			With(`name`, d.Name).
			With(`args`, d.Args))
	}
	return nil
}

func (c *converter) fromGoDirective(d *ast.Directive, target ir.Node) ir.Directive {
	switch d.Name {
	case `generate`, `build`, `noescape`, `uintptrescapes`:
		// Ignore
		return nil
	case `linkname`:
		return c.fromLinkName(d)
	case `nosplit`:
		c.fromBoolFuncDirective(d, target, func(fn *ir.FuncDecl) { fn.Atomic = true })
		return nil
	case `norace`:
		c.fromBoolFuncDirective(d, target, func(fn *ir.FuncDecl) { fn.NoRace = true })
		return nil
	case `noinline`:
		c.fromBoolFuncDirective(d, target, func(fn *ir.FuncDecl) { fn.NoInline = true })
		return nil
	case `wasmimport`:
		// TODO: Look into leveraging or add to ignored directives
		// See https://pkg.go.dev/cmd/compile
		return nil
	case `wasmexport`:
		// TODO: Look into leveraging or add to ignored directives
		// See https://pkg.go.dev/cmd/compile
		return nil
	default:
		c.addFault(faults.New(`unexpected go directive`).
			With(`pos`, c.pos(d.Pos())).
			With(`name`, d.Name).
			With(`args`, d.Args))
	}
	return nil
}

func (c *converter) fromGozerDirective(d *ast.Directive, target ir.Node) ir.Directive {
	switch d.Name {
	case `atomic`:
		c.fromBoolFuncDirective(d, target, func(fn *ir.FuncDecl) { fn.Atomic = true })
		return nil
	default:
		c.addFault(faults.New(`unexpected gozer directive`).
			With(`pos`, c.pos(d.Pos())).
			With(`name`, d.Name).
			With(`args`, d.Args))
	}
	return nil
}

func (c *converter) fromLinkName(d *ast.Directive) *ir.LinkName {
	args, err := d.ParseArgs()
	if err != nil {
		c.addFault(faults.New(`unexpected go:linkname directive arguments`, err).
			With(`pos`, c.pos(d.Pos())).
			With(`args`, d.Args))
	}
	if len(args) != 1 && len(args) != 2 {
		c.addFault(faults.New(`expected go:linkname directive to have one or two arguments`).
			With(`pos`, c.pos(d.Pos())).
			With(`count`, len(args)).
			With(`args`, d.Args))
	}

	ln := &ir.LinkName{
		LinkPos:   d.Pos(),
		LocalName: args[0].Arg,
	}
	if len(args) == 2 {
		remote := args[1].Arg
		if path, name, ok := strings.CutLast(remote, `.`); ok {
			ln.RemotePath = path
			ln.RemoteName = name
		} else {
			ln.RemoteName = remote
		}
	}
	return ln
}

func (c *converter) fromBoolFuncDirective(d *ast.Directive, target ir.Node, setBool func(fn *ir.FuncDecl)) {
	if fn, ok := target.(*ir.FuncDecl); ok {
		assert.EmptyStr(d.Args)
		setBool(fn)
	}
	c.addFault(faults.New(`expected directive to be on a function declaration`).
		With(`tool`, d.Tool).
		With(`name`, d.Name).
		WithF(`node type`, `%T`, target).
		With(`pos`, c.pos(d.Pos())))
}
