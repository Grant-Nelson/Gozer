package augmenter

import (
	"errors"
	"go/ast"
	"go/build/constraint"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/file"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/augmenter/directives"
)

func (a *Augmenter) AddPackage(path string, errGroup *faults.Group) (err error) {
	defer faults.Recover(&err)
	ar := &augReader{Augmenter: a, errGroup: errGroup}
	ar.addPackage(path)
	return nil
}

func (a *Augmenter) AddFile(filename string, src []byte, errGroup *faults.Group) (err error) {
	defer faults.Recover(&err)
	ar := &augReader{Augmenter: a, errGroup: errGroup}
	ar.addFile(filename, src)
	return nil
}

var (
	ErrParsingBuildConstraints = errors.New(`error parsing build constrains for augmentation file`)
	ErrParsingUnexpectedDecl   = errors.New(`unexpected declaration while parsing augmentation file`)
	ErrParsingUnexpectedSpec   = errors.New(`unexpected specification while parsing augmentation file`)
	ErrAugFuncNone             = errors.New(`a function must have a directive`)
	ErrAugSpecNone             = errors.New(`a specification must have a directive`)
	ErrAugGenWithFuncDirective = errors.New(`a general declaration may not have a directive for a function`)
	ErrAugRenameMultipleSpec   = errors.New(`names may not be applied to multiple constructs`)
	ErrAugRenameImport         = errors.New(`renames may not be applied to imports`)
	ErrAugDeleteAllImport      = errors.New(`delete all may not be applied to imports`)
	ErrAugDeleteAllValue       = errors.New(`delete all may not be applied to values`)
	ErrAugDeleteAllInterface   = errors.New(`delete all may not be applied to an interface`)
	ErrAugMethodNone           = errors.New(`an interface method must have a directive`)
	ErrAugMethodReplaceRecv    = errors.New(`an interface method may not have a replaceRecv`)
	ErrAugMethodDeleteAll      = errors.New(`an interface method may not have a deleteAll`)
	ErrAugMethodReplaceSig     = errors.New(`an interface method may not have a replaceSig, just use replace`)
	ErrAugFieldNone            = errors.New(`a struct field must have a directive`)
	ErrAugFieldReplaceRecv     = errors.New(`a struct field may not have a replaceRecv`)
	ErrAugFieldDeleteAll       = errors.New(`a struct field may not have a deleteAll`)
	ErrAugFieldReplaceSig      = errors.New(`a struct field may not have a replaceSig`)
)

type augReader struct {
	*Augmenter
	errGroup *faults.Group
	curFile  *file.File
	addSpecs []ast.Spec
}

func (ar *augReader) addPackage(path string) {
	dir := filepath.Join(ar.basePath, path)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}

	testDir := path == ar.testPkgPath
	for _, entry := range entries {
		if !entry.IsDir() {
			if filepath.Ext(entry.Name()) != `go` {
				continue
			}
			if !testDir || !strings.HasSuffix(entry.Name(), `_test.go`) {
				continue
			}
			ar.addFile(entry.Name(), nil)
		}
	}
}

func (ar *augReader) addFile(filename string, src []byte) {
	f, err := file.Load(ar.fileSet, filename, src)
	if err != nil {
		ar.errGroup.Panic(err)
		return
	}

	if !ar.shouldAdd(f) {
		return
	}

	ar.curFile = f
	for _, d := range f.File.Decls {
		ar.readDecl(d)
	}
	ar.curFile = nil
}

func (ar *augReader) shouldAdd(f *file.File) bool {
	if f.File.Doc == nil || len(f.File.Doc.List) <= 0 {
		return true
	}

	for _, com := range f.File.Doc.List {
		exp, err := constraint.Parse(com.Text)
		if err != nil {
			ar.errGroup.Panic(faults.From(ErrParsingBuildConstraints).
				With(`error`, err).
				With(`position`, ar.fileSet.Position(com.Pos())))
		}
		if !exp.Eval(func(tag string) bool {
			return slices.Contains(ar.build, tag)
		}) {
			return false
		}
	}
	return true
}

func (ar *augReader) pkgPath() string {
	return ar.curFile.PackagePath()
}

func (ar *augReader) pos(p token.Pos) token.Position {
	return ar.fileSet.Position(p)
}

func (ar *augReader) readDecl(decl ast.Decl) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		ar.readFuncDecl(d)
	case *ast.GenDecl:
		ar.readGenDecl(d)
	default:
		ar.errGroup.Panic(faults.From(ErrParsingUnexpectedDecl).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(d.Pos())))
	}
}

func (ar *augReader) readFuncDecl(fd *ast.FuncDecl) {
	pos := ar.pos(fd.Pos())
	dv, err := directives.Read(file.JoinComments(fd.Doc), ar.pkgPath(), pos, ar.errGroup)
	if err != nil {
		panic(err)
	}

	switch {
	case dv.Ignore():
		return
	case dv.None():
		ar.errGroup.Panic(faults.From(ErrAugFuncNone).
			With(`package path`, ar.pkgPath()).
			With(`function name`, fd.Name.Name).
			With(`position`, pos))
		return
	case dv.Add():
		ar.add.newFuncDecls = append(ar.add.newFuncDecls, fd)
		ar.add.beingAdded[fd.Name.Name] = fd.Pos()
	case dv.Delete():
		// TODO: Implement
		panic(errors.New(`unimplemented`))
	case dv.Replace():
		// TODO: Implement
		panic(errors.New(`unimplemented`))
	}
}

func (ar *augReader) readGenDecl(gd *ast.GenDecl) {
	declPos := ar.pos(gd.Pos())
	declDv, err := directives.Read(file.JoinComments(gd.Doc), ar.pkgPath(), declPos, ar.errGroup)
	if err != nil {
		panic(err)
	}

	if declDv.HasReplaceRecv() || declDv.ReplaceSig() {
		ar.errGroup.Panic(faults.From(ErrAugGenWithFuncDirective).
			With(`package path`, ar.pkgPath()).
			With(`position`, declPos))
	}

	if declDv.HasRename() && len(gd.Specs) != 1 {
		ar.errGroup.Panic(faults.From(ErrAugRenameMultipleSpec).
			With(`package path`, ar.pkgPath()).
			With(`position`, declPos).
			With(`name`, declDv.Rename()).
			With(`count`, len(gd.Specs)))
		return
	}

	for _, spec := range gd.Specs {
		specDv := ar.readSpecDirectives(declDv, spec)
		switch s := spec.(type) {
		case *ast.ImportSpec:
			ar.readImportSpec(specDv, gd, s)
		case *ast.TypeSpec:
			switch t := s.Type.(type) {
			case *ast.StructType:
				ar.readStructTypeSpec(specDv, gd, s, t)
			case *ast.InterfaceType:
				ar.readInterfaceTypeSpec(specDv, gd, s, t)
			default:
				ar.readOtherTypeSpec(specDv, gd, s)
			}
		case *ast.ValueSpec:
			ar.readValueSpec(specDv, gd, s)
		}
	}
	ar.finishGenDecl(gd)
}

func (ar *augReader) readSpecDirectives(declDv *directives.Directives, spec ast.Spec) *directives.Directives {
	specPos := ar.pos(spec.Pos())
	var comments []*ast.Comment
	switch s := spec.(type) {
	case *ast.ImportSpec:
		comments = file.JoinComments(s.Doc, s.Comment)
	case *ast.TypeSpec:
		comments = file.JoinComments(s.Doc, s.Comment)
	case *ast.ValueSpec:
		comments = file.JoinComments(s.Doc, s.Comment)
	default:
		ar.errGroup.Panic(faults.From(ErrParsingUnexpectedSpec).
			With(`package path`, ar.pkgPath()).
			With(`position`, specPos))
		return nil
	}

	specDv, err := directives.Read(comments, ar.pkgPath(), specPos, ar.errGroup)
	if err != nil {
		panic(err)
	}

	joinDv, err := declDv.Join(specDv, ar.pkgPath(), specPos, ar.errGroup)
	if err != nil {
		panic(err)
	}
	return joinDv
}

func (ar *augReader) readImportSpec(specDv *directives.Directives, gd *ast.GenDecl, spec *ast.ImportSpec) {
	switch {
	case specDv.Ignore(), specDv.None():
		// Imports default to ignore with none.
		return
	case specDv.HasRename():
		ar.errGroup.Panic(faults.From(ErrAugRenameImport).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`import path`, spec.Path.Value))
	case specDv.Add():
		ar.add.newImports = append(ar.add.newImports, spec)
		ar.add.beingAdded[spec.Path.Value] = spec.Pos()
	case specDv.DeleteAll():
		ar.errGroup.Panic(faults.From(ErrAugDeleteAllImport).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`import path`, spec.Path.Value))
	case specDv.Delete():
		// TODO: Implement
	case specDv.Replace():
		// TODO: Implement
	}
}

func (ar *augReader) readStructTypeSpec(specDv *directives.Directives, gd *ast.GenDecl, spec *ast.TypeSpec, ts *ast.StructType) {
	if specDv.Ignore() {
		return
	}

	for _, f := range ts.Fields.List {
		ar.readStructField(specDv, gd, spec, ts, f)
	}

	switch {
	case specDv.None():
		// TODO: Need to check if a type spec for field and method directives.
		ar.errGroup.Panic(faults.From(ErrAugSpecNone).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`structure`, spec.Name.Name))
	case specDv.Add():
		ar.addSpecs = append(ar.addSpecs, spec)
		ar.add.beingAdded[spec.Name.Name] = spec.Pos()
	case specDv.DeleteAll():
		// TODO: Implement
	case specDv.Delete():
		// TODO: Implement
	case specDv.HasRename():
		// TODO: Implement
	case specDv.Replace():
		// TODO: Implement
	}
}

func (ar *augReader) readStructField(specDv *directives.Directives, gd *ast.GenDecl, spec *ast.TypeSpec, ts *ast.StructType, m *ast.Field) {
	comments := file.JoinComments(m.Comment, m.Doc)
	mDv, err := directives.Read(comments, ar.pkgPath(), ar.pos(m.Pos()), ar.errGroup)
	if err != nil {
		panic(err)
	}
	joinDv, err := specDv.Join(mDv, ar.pkgPath(), ar.pos(m.Pos()), ar.errGroup)
	if err != nil {
		panic(err)
	}
	if joinDv.Ignore() {
		return
	}
	if joinDv.None() {
		ar.errGroup.Panic(faults.From(ErrAugFieldNone).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`struct`, spec.Name.Name).
			With(`field`, m.Names[0].Name))
		return
	}
	if mDv.HasReplaceRecv() {
		ar.errGroup.Panic(faults.From(ErrAugFieldReplaceRecv).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`struct`, spec.Name.Name).
			With(`field`, m.Names[0].Name))
	}
	if mDv.DeleteAll() {
		ar.errGroup.Panic(faults.From(ErrAugFieldDeleteAll).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`struct`, spec.Name.Name).
			With(`field`, m.Names[0].Name))
	}
	if mDv.ReplaceSig() {
		ar.errGroup.Panic(faults.From(ErrAugFieldReplaceSig).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`struct`, spec.Name.Name).
			With(`field`, m.Names[0].Name))
	}
	if mDv.HasRename() {
		// TODO: Implement
	}
	if !specDv.None() {
		// The specification will handle adding, deleting, or replacing as a whole.
		return
	}
	switch {
	case joinDv.Add():
		key := spec.Name.Name
		fieldsToAdd, has := ar.add.newFields[key]
		if !has {
			fieldsToAdd = &ast.StructType{
				Struct:     ts.Struct,
				Incomplete: true,
				Fields: &ast.FieldList{
					Opening: ts.Fields.Opening,
					Closing: ts.Fields.Closing,
				},
			}
			ar.add.newFields[key] = fieldsToAdd
		}
		fieldsToAdd.Fields.List = append(fieldsToAdd.Fields.List, m)
	case joinDv.Delete():
		// TODO: Implement
	case joinDv.Replace():
		// TODO: Implement
	}
}

func (ar *augReader) readInterfaceTypeSpec(specDv *directives.Directives, gd *ast.GenDecl, spec *ast.TypeSpec, ts *ast.InterfaceType) {
	switch {
	case specDv.Ignore():
		return
	case specDv.DeleteAll():
		ar.errGroup.Panic(faults.From(ErrAugDeleteAllInterface).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`interface`, spec.Name.Name))
	}

	for _, m := range ts.Methods.List {
		ar.readInterfaceMethod(specDv, gd, spec, ts, m)
	}

	switch {
	case specDv.Add():
		ar.addSpecs = append(ar.addSpecs, spec)
		ar.add.beingAdded[spec.Name.Name] = spec.Pos()
	case specDv.Delete():
		// TODO: Implement
	case specDv.HasRename():
		// TODO: Implement
	case specDv.Replace():
		// TODO: Implement
	}
}

func (ar *augReader) readInterfaceMethod(specDv *directives.Directives, gd *ast.GenDecl, spec *ast.TypeSpec, ts *ast.InterfaceType, m *ast.Field) {
	comments := file.JoinComments(m.Comment, m.Doc)
	mDv, err := directives.Read(comments, ar.pkgPath(), ar.pos(m.Pos()), ar.errGroup)
	if err != nil {
		panic(err)
	}
	joinDv, err := specDv.Join(mDv, ar.pkgPath(), ar.pos(m.Pos()), ar.errGroup)
	if err != nil {
		panic(err)
	}
	if joinDv.Ignore() {
		return
	}
	if joinDv.None() {
		ar.errGroup.Panic(faults.From(ErrAugMethodNone).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`interface`, spec.Name.Name).
			With(`method`, m.Names[0].Name))
		return
	}
	if mDv.HasReplaceRecv() {
		ar.errGroup.Panic(faults.From(ErrAugMethodReplaceRecv).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`interface`, spec.Name.Name).
			With(`method`, m.Names[0].Name))
	}
	if mDv.DeleteAll() {
		ar.errGroup.Panic(faults.From(ErrAugMethodDeleteAll).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`interface`, spec.Name.Name).
			With(`method`, m.Names[0].Name))
	}
	if mDv.ReplaceSig() {
		ar.errGroup.Panic(faults.From(ErrAugMethodReplaceSig).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`interface`, spec.Name.Name).
			With(`method`, m.Names[0].Name))
	}
	if mDv.HasRename() {
		// TODO: Implement
	}
	if !specDv.None() {
		// The specification will handle adding, deleting, or replacing as a whole.
		return
	}
	switch {
	case joinDv.Add():
		key := spec.Name.Name
		methodsToAdd, has := ar.add.newMethods[key]
		if !has {
			methodsToAdd = &ast.InterfaceType{
				Interface:  ts.Interface,
				Incomplete: true,
				Methods: &ast.FieldList{
					Opening: ts.Methods.Opening,
					Closing: ts.Methods.Closing,
				},
			}
			ar.add.newMethods[key] = methodsToAdd
		}
		methodsToAdd.Methods.List = append(methodsToAdd.Methods.List, m)
	case joinDv.Delete():
		// TODO: Implement
	case joinDv.Replace():
		// TODO: Implement
	}
}

func (ar *augReader) readOtherTypeSpec(specDv *directives.Directives, gd *ast.GenDecl, spec *ast.TypeSpec) {
	switch {
	case specDv.Ignore():
		return
	case specDv.None():
		ar.errGroup.Panic(faults.From(ErrAugSpecNone).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())))
	case specDv.Add():
		ar.addSpecs = append(ar.addSpecs, spec)
		ar.add.beingAdded[spec.Name.Name] = spec.Pos()
	case specDv.DeleteAll():
		// TODO: Implement
	case specDv.Delete():
		// TODO: Implement
	case specDv.HasRename():
		// TODO: Implement
	case specDv.Replace():
		// TODO: Implement
	}
}

func (ar *augReader) readValueSpec(specDv *directives.Directives, gd *ast.GenDecl, spec *ast.ValueSpec) {
	switch {
	case specDv.Ignore():
		return
	case specDv.None():
		ar.errGroup.Panic(faults.From(ErrAugSpecNone).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())))
	case specDv.Add():
		ar.addSpecs = append(ar.addSpecs, spec)
	case specDv.DeleteAll():
		ar.errGroup.Panic(faults.From(ErrAugDeleteAllValue).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())))
	case specDv.Delete():
		// TODO: Implement
	case specDv.HasRename():
		// TODO: Check Rename is only for one
	case specDv.Replace():
		// TODO: Implement
	}
}

func (ar *augReader) finishGenDecl(gd *ast.GenDecl) {
	if len(ar.addSpecs) > 0 {
		addGen := &ast.GenDecl{
			Doc:    gd.Doc,
			TokPos: gd.TokPos,
			Tok:    gd.Tok,
			Lparen: gd.Lparen,
			Rparen: gd.Rparen,
			Specs:  ar.addSpecs,
		}
		ar.add.newGenDecls = append(ar.add.newGenDecls, addGen)
	}
	// TODO: Implement
}
