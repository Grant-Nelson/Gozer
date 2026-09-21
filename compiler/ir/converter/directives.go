package converter

import (
	"go/ast"
	"strings"

	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
)

func (c *converter) FromCommentGroup(cg *ast.CommentGroup) []ir.Directive {
	ds := []ir.Directive{}
	for d := range astTools.DirectivesFromGroup(cg) {
		if dt := c.FromDirective(d); dt != nil {
			ds = append(ds, dt)
		}
	}
	return ds
}

func (c *converter) FromDirective(d *ast.Directive) ir.Directive {
	switch d.Tool + `:` + d.Name {
	case `go:linkname`:
		return c.fromLinkName(d)
	}
	return nil
}

func (c *converter) fromLinkName(d *ast.Directive) *ir.LinkName {
	args, err := d.ParseArgs()
	if err != nil {
		c.addFault(faults.New(`unexpected go:linkname directive arguments`, err).
			With(`args`, d.Args))
	}
	if len(args) != 1 && len(args) != 2 {
		c.addFault(faults.New(`expected go:linkname directive to have one or two arguments`, err).
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
