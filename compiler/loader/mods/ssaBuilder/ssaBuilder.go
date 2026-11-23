package ssaBuilder

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/loader/mods"
	"github.com/Grant-Nelson/Gozer/compiler/project"
	"golang.org/x/tools/go/ssa"
)

type Config struct {
	// ErrGroup is used to collect multiple errors.
	ErrGroup *faults.ErrGroup
}

type SsaBuilder struct {
	prog     *ssa.Program
	errGroup *faults.ErrGroup
}

type ssaBuilderMod struct {
	errGroup *faults.ErrGroup
	prog     *ssa.Program
	pkg      *project.Package
}

var (
	_ mods.ModFactory      = (*SsaBuilder)(nil)
	_ mods.Modifier        = (*ssaBuilderMod)(nil)
	_ mods.StartProjectExt = (*SsaBuilder)(nil)
)

func New(cfg *Config) *SsaBuilder {
	tc := &SsaBuilder{
		errGroup: cfg.ErrGroup,
	}
	return tc
}

func (tc *SsaBuilder) StartProject(f *project.Project) (err error) {
	tc.errGroup.Recover(&err)
	prog := ssa.NewProgram(f.FileSet, ssa.SanityCheckFunctions)
	tc.prog = prog
	f.Ssa = prog
	return nil
}

func (tc *SsaBuilder) StartPackage(pkg *project.Package) (bool, mods.Modifier, error) {
	mod := &ssaBuilderMod{
		errGroup: tc.errGroup,
		prog:     tc.prog,
		pkg:      pkg,
	}
	return true, mod, nil
}

func (tc *ssaBuilderMod) PackageDone() (con bool, err error) {
	tc.errGroup.Recover(&err)
	p := tc.pkg.Ast
	p2 := tc.prog.CreatePackage(p.Types, p.Syntax, p.TypesInfo, true)
	p2.Build()
	tc.pkg.Ssa = p2
	return true, nil
}
