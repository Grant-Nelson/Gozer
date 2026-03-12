// This package is for the trimmer remodeler will remove empty statements,
// code block statements, and parentheses in expressions.
//
// The code block statements are used for scoping in a function but since the
// [types.Info] has already been determined, the scoping will still work.
// The target language just needs to ensure the correct variables reference
// each other and renamed if needed to preserve the code.
package trimmer

import (
	"go/ast"
	"slices"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/modeler/irc"
	"github.com/Grant-Nelson/Gozer/project/modeler/remodel"
)

type Config struct {

	// ErrGroup is used to collect multiple errors.
	ErrGroup *faults.ErrGroup

	// KeepCodeBlocks will make the trimmer not remove code blocks, [ast.BlockStmt].
	KeepCodeBlocks bool

	// KeepParens will make the trimmer not remove parens, [ast.ParenExpr].
	KeepParens bool

	// KeepEmpty will make the trimmer not remove empty statements, [ast.EmptyStmt].
	KeepEmpty bool
}

type trimmer struct {
	errGroup   *faults.ErrGroup
	keepBlocks bool
	keepParens bool
	keepEmpty  bool
}

func New(cfg *Config) remodel.RemodelFactory {
	return &trimmer{
		errGroup:   cfg.ErrGroup,
		keepBlocks: cfg.KeepCodeBlocks,
		keepParens: cfg.KeepParens,
		keepEmpty:  cfg.KeepEmpty,
	}
}

func (t *trimmer) StartPackage(pkg *project.Package) (bool, remodel.Remodeler, error) {
	if t.keepBlocks && t.keepParens && t.keepEmpty {
		return true, nil, nil
	}
	tr := &trimmerRemodel{
		errGroup:   t.errGroup,
		pkg:        pkg,
		keepBlocks: t.keepBlocks,
		keepParens: t.keepParens,
		keepEmpty:  t.keepEmpty,
	}
	return true, tr, nil
}

type trimmerRemodel struct {
	errGroup   *faults.ErrGroup
	pkg        *project.Package
	keepBlocks bool
	keepParens bool
	keepEmpty  bool
}

func (t *trimmerRemodel) PackageDone() (bool, error) { return true, nil }

func (t *trimmerRemodel) RemodelFunc(f *irc.Func) (con bool, err error) {
	t.errGroup.Recover(&err)
	for _, b := range f.Blocks {
		t.trimAstBlocksFromIrcBlock(b)
		t.trimInAstStmtFromIrcBlock(b)
	}
	return true, t.errGroup.FullOrNil()
}

func (t *trimmerRemodel) trimAstBlocksFromIrcBlock(b *irc.Block) {
	s := b.Body
	for i := 0; i < len(s); i++ {
		if replace, ok := t.unwrapAstBlockFromIrcStmt(s[i]); ok {
			s = slices.Replace(s, i, i+1, replace...)
			i-- // step back so we can recheck new nodes
			continue
		}
	}
	b.Body = s
}

func (t *trimmerRemodel) unwrapAstBlockFromIrcStmt(s irc.Stmt) ([]irc.Stmt, bool) {
	if sg, ok := s.(*irc.BaseStmt); ok {
		switch st := sg.Stmt.(type) {
		case *ast.BlockStmt:
			if !t.keepBlocks {
				replacement := make([]irc.Stmt, len(st.List))
				for i, inner := range st.List {
					replacement[i] = &irc.BaseStmt{Stmt: inner}
				}
				return replacement, true
			}
		case *ast.EmptyStmt:
			if !t.keepEmpty {
				return []irc.Stmt{}, true
			}
		}
	}
	return nil, false
}

func (t *trimmerRemodel) trimInAstStmtFromIrcBlock(s *irc.Block) {
	for _, s := range s.Body {
		if sb, ok := s.(*irc.BaseStmt); ok {
			sb.Stmt = t.trimNode(sb.Stmt).(ast.Stmt)
		}
	}
}

func (t *trimmerRemodel) trimNode(n ast.Node) ast.Node {
	return astutil.Apply(n, func(c *astutil.Cursor) bool {
		switch p := c.Node().(type) {
		case *ast.ParenExpr:
			if !t.keepParens {
				c.Replace(p.X)
			}

		case *ast.EmptyStmt:
			if !t.keepEmpty {
				c.Delete()
			}

		case *ast.BlockStmt:
			if t.keepBlocks {
				return true
			}
			if c.Name() == `Body` {
				switch c.Parent().(type) {
				case *ast.IfStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt,
					*ast.SelectStmt, *ast.ForStmt, *ast.RangeStmt, *ast.FuncLit:
					// If the block is the Body of any of these nodes, then
					// it can't be replaces since the field is [*ast.BlockStmt].
					return true
				}
			}
			for i := len(p.List) - 1; i >= 0; i-- {
				c.InsertAfter(p.List[i])
			}
			c.Delete()
		}
		return true
	}, nil)
}
