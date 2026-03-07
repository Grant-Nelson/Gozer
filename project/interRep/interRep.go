package interRep

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/enums/buildState"
	"github.com/Grant-Nelson/Gozer/project/interRep/irc"
	remodel "github.com/Grant-Nelson/Gozer/project/interRep/remods"
)

type Config struct {
	// Logger to log verbose messages with. Has no affect if verbose was false.
	Logger *logger.Logger

	// ErrGroup is the collector to handle multiple errors.
	ErrGroup *faults.ErrGroup

	// Package is the package to create a IRC for.
	Package *project.Package

	// Remodelers are the tools to apply to the IRC to transform from
	// very similar to AST over to a form needed by the target.
	Remodelers remodel.Group
}

// Remodel converts the AST in the packages in the given configs
// into the IRC form and sets that to the package.
func Remodel(cfg *Config) (err error) {
	defer faults.Recover(&err)
	cfg.Logger.LogGroup(`Remodelling %q`, cfg.Package.PkgPath())

	pkg := cfg.Package
	pkg.State = buildState.Remodelling
	pkg.Irc = &irc.Package{}

	rm := &modeler{
		logger:   cfg.Logger,
		errGroup: cfg.ErrGroup,
		pkg:      pkg,
	}

	if err := rm.remodelPackageStart(cfg.Remodelers); err != nil {
		return err
	}

	for _, f := range cfg.Package.Ast.Syntax {
		if err := rm.addFile(f); err != nil {
			return err
		}
	}

	if err := rm.remodelPackageDone(); err != nil {
		return err
	}

	pkg.State = buildState.Remodelled
	return rm.errGroup.FullOrNil()
}

type modeler struct {
	logger   *logger.Logger
	errGroup *faults.ErrGroup
	pkg      *project.Package
	group    remodel.Remodeler
}

func (rm *modeler) pos(p token.Pos) token.Position {
	return rm.pkg.Ast.Fset.Position(p)
}

func (rm *modeler) remodelPackageStart(group remodel.Group) error {
	rm.logger.Printf(`Package Starting`)
	defer rm.logger.Indent()()

	con, rg, err := group.StartPackage(rm.pkg)
	if err != nil || !con {
		return rm.errGroup.Add(err)
	}
	rm.group = rg
	return rm.errGroup.FullOrNil()
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
		Ast:  astFunc,
		Name: astFunc.Name.Name,
	}
	rm.pkg.Irc.Funcs = append(rm.pkg.Irc.Funcs, fn)

	// Create initial block and populate it with current statements.
	block := fn.NewBlock()
	if astFunc.Body != nil {
		for _, s := range astFunc.Body.List {
			block.Body = append(block.Body, &irc.BaseStmt{Stmt: s})
		}
	}

	if rm.group != nil {
		if m, ok := rm.group.(remodel.RemodelFuncExt); ok {
			if _, err := m.RemodelFunc(fn); err != nil {
				rm.logger.Printf(`Error Remodeling function`)
				return rm.errGroup.Add(err)
			}
		}
	}
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

func (rm *modeler) remodelPackageDone() error {
	if rm.group == nil {
		return nil
	}

	rm.logger.Printf(`Package Finishing`)
	defer rm.logger.Indent()()

	if _, err := rm.group.PackageDone(); err != nil {
		return rm.errGroup.Add(err)
	}
	return rm.errGroup.FullOrNil()
}
