package ssaSetup

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

type SsaSetup struct {
	prog     *ssa.Program
	errGroup *faults.ErrGroup
}

type ssaSetupMod struct {
	errGroup *faults.ErrGroup
	prog     *ssa.Program
	pkg      *project.Package
}

var (
	_ mods.ModFactory      = (*SsaSetup)(nil)
	_ mods.Modifier        = (*ssaSetupMod)(nil)
	_ mods.StartProjectExt = (*SsaSetup)(nil)
)

func New(cfg *Config) *SsaSetup {
	tc := &SsaSetup{
		errGroup: cfg.ErrGroup,
	}
	return tc
}

func (tc *SsaSetup) StartProject(f *project.Project) (err error) {
	defer tc.errGroup.Recover(&err)
	prog := ssa.NewProgram(f.FileSet, ssa.SanityCheckFunctions)
	tc.prog = prog
	f.Ssa = prog
	return nil
}

func (tc *SsaSetup) StartPackage(pkg *project.Package) (bool, mods.Modifier, error) {
	mod := &ssaSetupMod{
		errGroup: tc.errGroup,
		prog:     tc.prog,
		pkg:      pkg,
	}
	return true, mod, nil
}

func (tc *ssaSetupMod) PackageDone() (con bool, err error) {
	defer tc.errGroup.Recover(&err)
	p := tc.pkg.Ast
	tc.pkg.Ssa = tc.prog.CreatePackage(p.Types, p.Syntax, p.TypesInfo, true)
	return true, nil
}
