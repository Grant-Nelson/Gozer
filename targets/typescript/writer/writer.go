package writer

import (
	"path/filepath"

	"github.com/Grant-Nelson/Gozer/avail/coder"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
	"github.com/Grant-Nelson/Gozer/compiler/modeler/remodel"
	"github.com/Grant-Nelson/Gozer/compiler/project"
)

type Config struct {
	ErrGroup *faults.ErrGroup
}

type writerFactory struct {
	errGroup *faults.ErrGroup
}

func New(cfg *Config) remodel.RemodelFactory {
	return &writerFactory{
		errGroup: cfg.ErrGroup,
	}
}

func (f *writerFactory) StartPackage(pkg *project.Package) (bool, remodel.Remodeler, error) {
	w := &writer{
		errGroup: f.errGroup,
		pkg:      pkg,
	}
	return true, w, nil
}

type writer struct {
	errGroup *faults.ErrGroup
	pkg      *project.Package
	c        *coder.Coder
}

func (w *writer) getModuleFilePath(pkgPath string) string {
	// TODO: Need to properly create the file path to output.
	// The name probably should be part of the IR so that caching and writing
	// can always check the same path.
	return filepath.Join(pkgPath, `temp.ts`)
}

func (w *writer) Remodel() (con bool, err error) {
	defer w.errGroup.Recover(&err)

	modulePath := w.getModuleFilePath(w.pkg.PkgPath())
	w.c = coder.New(modulePath, `  `, `// `)

	w.writePackageComment()
	w.writeImports()
	w.writeConsts()
	w.writeTypeDefs()
	w.writeVars()
	w.writeFuncs()

	return true, w.errGroup.FullOrNil()
}

func (w *writer) PackageDone() (con bool, err error) {
	defer w.errGroup.Recover(&err)

	// TODO: Implement saving module file

	return true, nil
}

func (w *writer) writePackageComment() {
	if comment := w.pkg.Ir.Comment; len(comment) > 0 {
		w.c.WriteComment(comment)
		w.c.WriteLn()
	}
}

func (w *writer) writeImports() {
	if imports := w.pkg.Ir.Imports; len(imports) > 0 {
		for _, imp := range imports {
			w.writeImport(imp)
		}
		w.c.WriteLn()
	}
}

func (w *writer) writeImport(imp *ir.ImportDef) {
	modulePath := w.getModuleFilePath(imp.PkgObj.Pkg().Path())
	name := imp.PkgObj.Name()
	w.c.WriteF(`import { %s } from '%s'`, name, modulePath)
	w.c.WriteLn()
}

func (w *writer) writeConsts() {
	for _, ct := range w.pkg.Ir.Consts {
		w.writeConst(ct)
	}
}

func (w *writer) writeConst(ct *ir.ConstDef) {
	w.c.WriteComment(ct.Comment)
	if ct.ConstObj.Exported() {
		w.c.Write(`export `)
	}
	w.c.Write(`const `, ct.ConstObj.Name(), ` = `)
	w.c.Write(ct.ConstObj.Val().ExactString())
	w.c.WriteLn(`;`)
}

func (w *writer) writeTypeDefs() {
	for _, t := range w.pkg.Ir.Types {
		w.writeTypeDef(t)
	}
}

func (w *writer) writeTypeDef(t *ir.TypeDef) {
	// TODO: Implement
}

func (w *writer) writeVars() {
	for _, vr := range w.pkg.Ir.Vars {
		w.writeVar(vr)
	}
}

func (w *writer) writeVar(vr *ir.VarDef) {
	w.c.WriteComment(vr.Comment)
	if vr.VarObj.Exported() {
		w.c.Write(`export `)
	}
	w.c.WriteLn(`var `, vr.VarObj.Name(), ` = `)
	w.writeExpr(vr.Value)
	w.c.WriteLn(`;`)
}

func (w *writer) writeExpr(expr ir.Expr) {
	switch expr := expr.(type) {

	// TODO: Implement

	default:
		w.errGroup.Add(faults.New(`Unexpected expression during write`).
			With(`pos`, w.pkg.Position(expr.Pos())).
			With(`expr`, expr.String()).
			WithTypeOf(`type`, expr))
	}
}

func (w *writer) writeFuncs() {
	for _, fn := range w.pkg.Ir.Funcs {
		w.writeFunc(fn)
	}
}

func (w *writer) writeFunc(fn *ir.FuncDef) {
	// TODO: Implement
	w.writeBlocks(fn.Func.Blocks)
}

func (w *writer) writeBlocks(blocks []*ir.Block) {
	for _, blk := range blocks {
		w.writeBlock(blk)
	}
}

func (w *writer) writeBlock(blk *ir.Block) {
	w.c.WriteComment(blk.Hint)

	// TODO: Implement
}
