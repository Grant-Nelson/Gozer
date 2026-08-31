package converter

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/converter"
	"github.com/Grant-Nelson/Gozer/compiler/loader/mods"
	"github.com/Grant-Nelson/Gozer/compiler/project"
)

type Config struct {

	// KeepAst indicates that the AST on the package should NOT be
	// set to nil when done converting into the IR.
	KeepAst bool

	// ErrGroup is used to collect multiple errors.
	ErrGroup *faults.ErrGroup
}

type Converter struct {
	keepAst  bool
	errGroup *faults.ErrGroup
}

type converterMod struct {
	keepAst  bool
	errGroup *faults.ErrGroup
	pkg      *project.Package
}

var (
	_ mods.ModFactory = (*Converter)(nil)
	_ mods.Modifier   = (*converterMod)(nil)
)

func New(cfg *Config) *Converter {
	return &Converter{
		keepAst:  cfg.KeepAst,
		errGroup: cfg.ErrGroup,
	}
}

func (cv *Converter) StartPackage(pkg *project.Package) (bool, mods.Modifier, error) {
	mod := &converterMod{
		keepAst:  cv.keepAst,
		errGroup: cv.errGroup,
		pkg:      pkg,
	}
	return true, mod, nil
}

func (cv *converterMod) PackageDone() (bool, error) {
	ir, err := converter.ConvertPackage(cv.pkg.Ast, cv.errGroup)
	if err != nil {
		return false, err
	}
	cv.pkg.Ir = ir
	if !cv.keepAst {
		cv.pkg.Ast = nil
	}
	return true, nil
}
