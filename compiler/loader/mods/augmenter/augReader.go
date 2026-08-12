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

	"github.com/Grant-Nelson/Gozer/avail/astTools"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/augmenter/directives"
	"github.com/Grant-Nelson/Gozer/project/loader/parser"
)

var (
	ErrParsingBuildConstraints  = errors.New(`error parsing build constrains for augmentation file`)
	ErrParsingUnexpectedDecl    = errors.New(`unexpected declaration while parsing augmentation file`)
	ErrParsingUnexpectedGenDecl = errors.New(`unexpected general declaration while parsing augmentation file`)
	ErrParsingUnexpectedSpec    = errors.New(`unexpected specification while parsing augmentation file`)
	ErrAugFuncNone              = errors.New(`a function must have a directive`)
	ErrAugSpecNone              = errors.New(`a specification must have a directive`)
	ErrAugGenWithFuncDirective  = errors.New(`a general declaration may not have a directive for a function`)
	ErrAugRenameMultipleSpec    = errors.New(`names may not be applied to multiple constructs`)
	ErrAugRenameImport          = errors.New(`renames may not be applied to imports`)
	ErrAugDeleteAllImport       = errors.New(`delete all may not be applied to imports`)
	ErrAugDeleteAllValue        = errors.New(`delete all may not be applied to values`)
	ErrAugDeleteAllInterface    = errors.New(`delete all may not be applied to an interface`)
	ErrAugMethodNone            = errors.New(`an interface method must have a directive`)
	ErrAugMethodReplaceRecv     = errors.New(`an interface method may not have a replaceRecv`)
	ErrAugMethodDeleteAll       = errors.New(`an interface method may not have a deleteAll`)
	ErrAugMethodReplaceSig      = errors.New(`an interface method may not have a replaceSig, just use replace`)
	ErrAugFieldNone             = errors.New(`a struct field must have a directive`)
	ErrAugFieldReplaceRecv      = errors.New(`a struct field may not have a replaceRecv`)
	ErrAugFieldDeleteAll        = errors.New(`a struct field may not have a deleteAll`)
	ErrAugFieldReplaceSig       = errors.New(`a struct field may not have a replaceSig`)
)

type augReader struct {
	*augPackage
	parser   parser.Parser
	errGroup *faults.ErrGroup
	build    []string

	curFile *ast.File

	addSpecs        []ast.Spec
	addSpecComments []*ast.CommentGroup
}

func newReader(pkg *augPackage, build []string, parser parser.Parser, errGroup *faults.ErrGroup) *augReader {
	return &augReader{
		augPackage: pkg,
		parser:     parser,
		errGroup:   errGroup,
		build:      build,
	}
}

func (ar *augReader) readPackage(dir string, data any) {
	// TODO: Handle the data if it is not nil, e.g. an embed.FS.

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		ar.errGroup.Add(err)
		panic(ar.errGroup)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			ar.addFile(entry.Name(), nil)
		}
	}
}

func (ar *augReader) addFile(filename string, src []byte) {
	if filepath.Ext(filename) != `.go` {
		return
	}

	if strings.HasSuffix(filename, `_test.go`) != ar.pkg.IsTest() {
		return
	}

	f, err := ar.parser(ar.pkg.Ast.Fset, filename, src)
	if err != nil {
		ar.errGroup.Add(err)
		panic(ar.errGroup)
	}

	ar.curFile = f
	if !ar.shouldAdd() {
		return
	}
	for _, d := range f.Decls {
		ar.readDecl(d)
	}
	ar.curFile = nil
}

func (ar *augReader) pkgPath() string {
	return ar.pkg.PkgPath()
}

func (ar *augReader) pos(p token.Pos) token.Position {
	return ar.pkg.Position(p)
}

func (ar *augReader) shouldAdd() bool {
	if astTools.IsTest(ar.pkg.Ast.Fset, ar.curFile) != ar.pkg.IsTest() ||
		astTools.IsXTest(ar.curFile) != ar.pkg.IsXTest() {
		return false
	}

	// Check build constraints
	if ar.curFile.Doc == nil || len(ar.curFile.Doc.List) <= 0 {
		return true // No build constraints
	}

	for _, com := range ar.curFile.Doc.List {
		exp, err := constraint.Parse(com.Text)
		if err != nil {
			ar.errGroup.Add(faults.From(ErrParsingBuildConstraints).
				With(`error`, err).
				With(`position`, ar.pos(com.Pos())))
			panic(ar.errGroup)
		}

		if !exp.Eval(func(tag string) bool {
			return slices.Contains(ar.build, tag)
		}) {
			return false
		}
	}
	return true
}

func (ar *augReader) readDecl(decl ast.Decl) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		ar.readFuncDecl(d)
	case *ast.GenDecl:
		ar.readGenDecl(d)
	default:
		ar.errGroup.Add(faults.From(ErrParsingUnexpectedDecl).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(d.Pos())))
		panic(ar.errGroup)
	}
}

func (ar *augReader) readFuncDecl(fd *ast.FuncDecl) {
	pos := ar.pos(fd.Pos())
	dv, err := directives.Read(astTools.JoinComments(fd.Doc), ar.pkgPath(), pos)
	if err != nil {
		ar.errGroup.Add(err)
		panic(ar.errGroup)
	}
	directives.RemoveDirectives(fd.Doc)

	switch {
	case dv.Ignore():
		return
	case dv.None():
		ar.errGroup.Add(faults.From(ErrAugFuncNone).
			With(`package path`, ar.pkgPath()).
			With(`function name`, fd.Name.Name).
			With(`position`, pos))
		panic(ar.errGroup)
	case dv.Add():
		ar.add.newDecls = append(ar.add.newDecls, fd)
		localComments := astTools.CommentsAttachedToNode(fd)
		ar.add.newDeclsComments = append(ar.add.newDeclsComments, localComments)
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
	declDv, err := directives.Read(astTools.JoinComments(gd.Doc), ar.pkgPath(), declPos)
	if err != nil {
		ar.errGroup.Add(err)
		panic(ar.errGroup)
	}
	directives.RemoveDirectives(gd.Doc)

	if declDv.HasReplaceRecv() || declDv.ReplaceSig() {
		ar.errGroup.Add(faults.From(ErrAugGenWithFuncDirective).
			With(`package path`, ar.pkgPath()).
			With(`position`, declPos))
		panic(ar.errGroup)
	}

	if declDv.HasRename() && len(gd.Specs) != 1 {
		ar.errGroup.Add(faults.From(ErrAugRenameMultipleSpec).
			With(`package path`, ar.pkgPath()).
			With(`position`, declPos).
			With(`name`, declDv.Rename()).
			With(`count`, len(gd.Specs)))
		panic(ar.errGroup)
	}

	switch gd.Tok {
	case token.IMPORT:
		ar.readImportDecl(declDv, gd)
	case token.TYPE:
		ar.readTypeDecl(declDv, gd)
	case token.VAR, token.CONST:
		ar.readValueDecl(declDv, gd)
	default:
		ar.errGroup.Add(faults.From(ErrParsingUnexpectedGenDecl).
			With(`package path`, ar.pkgPath()).
			With(`token`, gd.Tok.String()).
			With(`position`, ar.pos(gd.Pos())))
		panic(ar.errGroup)
	}
}

func (ar *augReader) readSpecDirectives(declDv *directives.Directives, pos token.Pos, doc, comment *ast.CommentGroup) *directives.Directives {
	specPos := ar.pos(pos)
	comments := astTools.JoinComments(doc, comment)
	directives.RemoveDirectives(doc)
	directives.RemoveDirectives(comment)

	specDv, err := directives.Read(comments, ar.pkgPath(), specPos)
	if err != nil {
		ar.errGroup.Add(err)
		panic(ar.errGroup)
	}

	joinDv, err := declDv.Join(specDv, ar.pkgPath(), specPos)
	if err != nil {
		ar.errGroup.Add(err)
		panic(ar.errGroup)
	}
	return joinDv
}

func (ar *augReader) readImportDecl(declDv *directives.Directives, gd *ast.GenDecl) {
	for _, spec := range gd.Specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			specDv := ar.readSpecDirectives(declDv, s.Pos(), s.Doc, s.Comment)
			ar.readImportSpec(specDv, gd, s)
		default:
			ar.errGroup.Add(faults.From(ErrParsingUnexpectedSpec).
				With(`package path`, ar.pkgPath()).
				With(`position`, spec.Pos()))
			panic(ar.errGroup)
		}
	}
	ar.finishGenDecl(gd)
}

func (ar *augReader) readImportSpec(specDv *directives.Directives, gd *ast.GenDecl, spec *ast.ImportSpec) {
	switch {
	case specDv.Ignore(), specDv.None():
		// Imports default to ignore with none.
		return
	case specDv.HasRename():
		ar.errGroup.Add(faults.From(ErrAugRenameImport).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`import path`, spec.Path.Value))
		panic(ar.errGroup)
	case specDv.Add():
		ar.addSpecs = append(ar.addSpecs, spec)
		ar.add.newImportSpecs = append(ar.add.newImportSpecs, spec)
		ar.add.beingAdded[spec.Path.Value] = spec.Pos()
	case specDv.DeleteAll():
		ar.errGroup.Add(faults.From(ErrAugDeleteAllImport).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`import path`, spec.Path.Value))
		panic(ar.errGroup)
	case specDv.Delete():
		// TODO: Implement
	case specDv.Replace():
		// TODO: Implement
	}
}

func (ar *augReader) readTypeDecl(declDv *directives.Directives, gd *ast.GenDecl) {
	for _, spec := range gd.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			specDv := ar.readSpecDirectives(declDv, s.Pos(), s.Doc, s.Comment)
			switch t := s.Type.(type) {
			case *ast.StructType:
				ar.readStructTypeSpec(specDv, gd, s, t)
			case *ast.InterfaceType:
				ar.readInterfaceTypeSpec(specDv, gd, s, t)
			default:
				ar.readOtherTypeSpec(specDv, gd, s)
			}
		default:
			ar.errGroup.Add(faults.From(ErrParsingUnexpectedSpec).
				With(`package path`, ar.pkgPath()).
				With(`position`, spec.Pos()))
			panic(ar.errGroup)
		}
	}
	ar.finishGenDecl(gd)
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
		ar.errGroup.Add(faults.From(ErrAugSpecNone).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`structure`, spec.Name.Name))
		panic(ar.errGroup)
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
	comments := astTools.JoinComments(m.Comment, m.Doc)
	directives.RemoveDirectives(m.Comment)
	directives.RemoveDirectives(m.Doc)
	mDv, err := directives.Read(comments, ar.pkgPath(), ar.pos(m.Pos()))
	if err != nil {
		ar.errGroup.Add(err)
		panic(ar.errGroup)
	}
	joinDv, err := specDv.Join(mDv, ar.pkgPath(), ar.pos(m.Pos()))
	if err != nil {
		ar.errGroup.Add(err)
		panic(ar.errGroup)
	}
	if joinDv.Ignore() {
		return
	}
	if joinDv.None() {
		ar.errGroup.Add(faults.From(ErrAugFieldNone).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`struct`, spec.Name.Name).
			With(`field`, m.Names[0].Name))
		panic(ar.errGroup)
	}
	if mDv.HasReplaceRecv() {
		ar.errGroup.Add(faults.From(ErrAugFieldReplaceRecv).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`struct`, spec.Name.Name).
			With(`field`, m.Names[0].Name))
		panic(ar.errGroup)
	}
	if mDv.DeleteAll() {
		ar.errGroup.Add(faults.From(ErrAugFieldDeleteAll).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`struct`, spec.Name.Name).
			With(`field`, m.Names[0].Name))
		panic(ar.errGroup)
	}
	if mDv.ReplaceSig() {
		ar.errGroup.Add(faults.From(ErrAugFieldReplaceSig).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`struct`, spec.Name.Name).
			With(`field`, m.Names[0].Name))
		panic(ar.errGroup)
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
		ar.errGroup.Add(faults.From(ErrAugDeleteAllInterface).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`interface`, spec.Name.Name))
		panic(ar.errGroup)
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
	comments := astTools.JoinComments(m.Comment, m.Doc)
	directives.RemoveDirectives(m.Comment)
	directives.RemoveDirectives(m.Doc)
	mDv, err := directives.Read(comments, ar.pkgPath(), ar.pos(m.Pos()))
	if err != nil {
		ar.errGroup.Add(err)
		panic(ar.errGroup)
	}
	joinDv, err := specDv.Join(mDv, ar.pkgPath(), ar.pos(m.Pos()))
	if err != nil {
		ar.errGroup.Add(err)
		panic(ar.errGroup)
	}
	if joinDv.Ignore() {
		return
	}
	if joinDv.None() {
		ar.errGroup.Add(faults.From(ErrAugMethodNone).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`interface`, spec.Name.Name).
			With(`method`, m.Names[0].Name))
		panic(ar.errGroup)
	}
	if mDv.HasReplaceRecv() {
		ar.errGroup.Add(faults.From(ErrAugMethodReplaceRecv).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`interface`, spec.Name.Name).
			With(`method`, m.Names[0].Name))
		panic(ar.errGroup)
	}
	if mDv.DeleteAll() {
		ar.errGroup.Add(faults.From(ErrAugMethodDeleteAll).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`interface`, spec.Name.Name).
			With(`method`, m.Names[0].Name))
		panic(ar.errGroup)
	}
	if mDv.ReplaceSig() {
		ar.errGroup.Add(faults.From(ErrAugMethodReplaceSig).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())).
			With(`interface`, spec.Name.Name).
			With(`method`, m.Names[0].Name))
		panic(ar.errGroup)
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
		ar.errGroup.Add(faults.From(ErrAugSpecNone).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())))
		panic(ar.errGroup)
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

func (ar *augReader) readValueDecl(declDv *directives.Directives, gd *ast.GenDecl) {
	for _, spec := range gd.Specs {
		switch s := spec.(type) {
		case *ast.ValueSpec:
			specDv := ar.readSpecDirectives(declDv, s.Pos(), s.Doc, s.Comment)
			ar.readValueSpec(specDv, gd, s)
		default:
			ar.errGroup.Add(faults.From(ErrParsingUnexpectedSpec).
				With(`package path`, ar.pkgPath()).
				With(`position`, spec.Pos()))
			panic(ar.errGroup)
		}
	}
	ar.finishGenDecl(gd)
}

func (ar *augReader) readValueSpec(specDv *directives.Directives, gd *ast.GenDecl, spec *ast.ValueSpec) {
	switch {
	case specDv.Ignore():
		return
	case specDv.None():
		ar.errGroup.Add(faults.From(ErrAugSpecNone).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())))
		panic(ar.errGroup)
	case specDv.Add():
		ar.addSpecs = append(ar.addSpecs, spec)
	case specDv.DeleteAll():
		ar.errGroup.Add(faults.From(ErrAugDeleteAllValue).
			With(`package path`, ar.pkgPath()).
			With(`position`, ar.pos(spec.Pos())))
		panic(ar.errGroup)
	case specDv.Delete():
		// TODO: Implement
	case specDv.HasRename():
		// TODO: Check Rename is only for one
	case specDv.Replace():
		// TODO: Implement
	}
}

func (ar *augReader) finishGenDecl(gd *ast.GenDecl) {
	if len(ar.addSpecs) <= 0 {
		return
	}

	addGen := &ast.GenDecl{
		Doc:    gd.Doc,
		TokPos: gd.TokPos,
		Tok:    gd.Tok,
		Lparen: gd.Lparen,
		Rparen: gd.Rparen,
		Specs:  ar.addSpecs,
	}

	switch gd.Tok {
	case token.IMPORT:
		ar.add.newImports = append(ar.add.newImports, addGen)
	default:
		ar.add.newDecls = append(ar.add.newDecls, addGen)
		localComments := astTools.CommentsAttachedToNode(addGen)
		ar.add.newDeclsComments = append(ar.add.newDeclsComments, localComments)
	}

	ar.addSpecs = []ast.Spec{}
}
