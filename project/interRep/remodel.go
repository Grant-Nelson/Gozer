package interRep

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/interRep/irc"
)

type Config struct {
	// Logger to log verbose messages with. Has no affect if verbose was false.
	Logger *logger.Logger

	// ErrGroup is the collector to handle multiple errors.
	ErrGroup *faults.ErrGroup

	// Package is the package to create a IRC for.
	Package *project.Package
}

// Remodel converts the AST in the packages in the given configs
// into the IRC form and sets that to the package.
func Remodel(cfg *Config) (err error) {
	defer faults.Recover(&err)
	cfg.Logger.LogGroup(`Remodelling %q`, cfg.Package.PkgPath())

	pkg := cfg.Package
	pkg.Irc = &irc.Package{}

	rm := &modeler{
		logger:   cfg.Logger,
		errGroup: cfg.ErrGroup,
		pkg:      pkg,
	}
	for _, f := range cfg.Package.Ast.Syntax {
		if err := rm.addFile(f); err != nil {
			return err
		}
	}

	// TODO: run any additional processing

	return rm.errGroup.AnyOrNil()
}

type modeler struct {
	logger   *logger.Logger
	errGroup *faults.ErrGroup
	pkg      *project.Package
}

func (rm *modeler) pos(p token.Pos) token.Position {
	return rm.pkg.Ast.Fset.Position(p)
}

func (rm *modeler) addFile(f *ast.File) error {
	for it := range astTools.DeclSpecs(rm.pkg.Ast.Fset, f) {
		switch n := it.Node.(type) {
		case *ast.ImportSpec:
			// ignore import specs
		case *ast.FuncDecl:
			if err := rm.addFuncDecl(n); err != nil {
				return err
			}
		case *ast.TypeSpec:
			if err := rm.addTypeSpec(it.GenDecl, n); err != nil {
				return err
			}
		case *ast.ValueSpec:
			if err := rm.addValueSpec(it.GenDecl, n); err != nil {
				return err
			}
		default:
			err := faults.New(`unexpected node type`).
				With(`pos`, rm.pos(n.Pos())).
				WithF(`type`, `%T`, n)
			if err := rm.errGroup.Add(err); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rm *modeler) addFuncDecl(astFunc *ast.FuncDecl) error {
	fn := &irc.Func{
		Ast: astFunc,
	}
	rm.pkg.Irc.Funcs = append(rm.pkg.Irc.Funcs, fn)
	cv := &converter{
		logger:   rm.logger,
		errGroup: rm.errGroup,
		pkg:      rm.pkg,
		fn:       fn,
	}

	// Create entry block for this function.
	first := fn.NewBlock()
	last, err := cv.blockStmt(first, astFunc.Body)
	if err != nil {
		return err
	}

	// If last doesn't have an breaking block control, add the implicit function return.

	// TODO: Finish
	_ = last

	return nil
}

func (rm *modeler) addTypeSpec(astGen *ast.GenDecl, astTypr *ast.TypeSpec) error {
	// TODO: Implement
	return nil
}

func (rm *modeler) addValueSpec(astGen *ast.GenDecl, astValue *ast.ValueSpec) error {
	// TODO: Implement
	return nil
}
