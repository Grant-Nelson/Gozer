package converter

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/dictionary"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/readers/golang/packageSet"
)

type Converter interface {
	Convert(pkg *build.Package) (cPkg *constructs.CPackage, err error)
}

type converter struct {
	pkg      *build.Package
	pkgSet   packageSet.PackageSet
	fSet     *token.FileSet
	proj     *constructs.CProject
	cPkg     *constructs.CPackage
	fImports collections.Dictionary[string, *constructs.CImport]
}

func New(pkgSet packageSet.PackageSet, fSet *token.FileSet, proj *constructs.CProject) Converter {
	return &converter{
		pkg:      nil,
		pkgSet:   pkgSet,
		fSet:     fSet,
		proj:     proj,
		cPkg:     nil,
		fImports: dictionary.New[string, *constructs.CImport](),
	}
}

func (con *converter) Convert(pkg *build.Package) (cPkg *constructs.CPackage, err error) {
	defer func() {
		if r := recover(); r != nil {
			cPkg = nil
			err = terror.RecoveredPanic(r).
				With(`package name`, pkg.Name)
		}
	}()

	con.pkg = pkg
	con.createCPackage()
	con.createCImports()
	enumerator.Enumerate(pkg.GoFiles...).Foreach(con.addFile)
	con.pkg = nil
	return cPkg, nil
}

func (con *converter) createCPackage() {
	con.cPkg = constructs.NewPackage(con.pkg.Name, con.pkg.Dir)
}

func (con *converter) createCImports() {
	con.cPkg.Imports.AddFrom(enumerator.Select(
		enumerator.Enumerate(con.pkg.Imports...),
		constructs.NewImport))
}

func (con *converter) addFile(fileName string) {
	fileName = con.prepareFileName(fileName)
	f := con.parseFile(fileName)
	enumerator.Enumerate(f.Decls...).Foreach(con.addDecl)
	con.finishFile()
}

func (con *converter) prepareFileName(fileName string) string {
	if !filepath.IsAbs(fileName) {
		fileName = filepath.Join(con.pkg.Dir, fileName)
	}
	return fileName
}

func (con *converter) parseFile(fileName string) *ast.File {
	f, err := parser.ParseFile(con.pkgSet.FileSet(), fileName, nil, parser.ParseComments)
	if err != nil {
		panic(terror.New(`error parsing file`, err).
			With(`file name`, fileName))
	}
	return f
}

func (con *converter) finishFile() {
	con.fImports.Clear()
}

func (con *converter) pos(p token.Pos) string {
	return con.pkgSet.FileSet().Position(p).String()
}

func (con *converter) addDecl(decl ast.Decl) {
	switch tDecl := decl.(type) {
	case *ast.FuncDecl:
		con.addFunc(tDecl)
	case *ast.GenDecl:
		con.addGenDecl(tDecl)
	default:
		panic(terror.New(`unknown declaration in code`).
			With(`decl node`, decl).
			With(`from`, con.pos(decl.Pos())).
			With(`to`, con.pos(decl.End())))
	}
}

func (con *converter) addFunc(funcDecl *ast.FuncDecl) {
	// TODO: Implement
}

func (con *converter) addGenDecl(gDecl *ast.GenDecl) {
	enumerator.Enumerate(gDecl.Specs...).Foreach(con.addSpec)
}

func (con *converter) addSpec(spec ast.Spec) {
	switch tSpec := spec.(type) {
	case *ast.ImportSpec:
		con.addImportSpec(tSpec)
	case *ast.TypeSpec:
		con.addTypeSpec(tSpec)
	case *ast.ValueSpec:
		con.addValueSpec(tSpec)
	default:
		panic(terror.New(`unknown specification in code`).
			With(`spec node`, spec).
			With(`from`, con.pos(spec.Pos())).
			With(`to`, con.pos(spec.End())))
	}
}

func (con *converter) getImportPath(iSpec *ast.ImportSpec) string {
	path, err := strconv.Unquote(iSpec.Path.Value)
	if err != nil {
		panic(terror.New(`failed to unquote import path`, err).
			With(`from`, con.pos(iSpec.Pos())).
			With(`to`, con.pos(iSpec.End())))
	}
	return path
}

func (con *converter) getImportName(iSpec *ast.ImportSpec) string {
	if iSpec.Name != nil && len(iSpec.Name.Name) > 0 {
		return iSpec.Name.Name
	}

	path := con.getImportPath(iSpec)
	other, has := con.pkgSet.Packages().TryGet(path)
	if !has {
		panic(terror.New(`failed to find package for import path`).
			With(`path`, path).
			With(`from`, con.pos(iSpec.Pos())).
			With(`to`, con.pos(iSpec.End())))
	}
	return other.Name
}

func (con *converter) addImportSpec(iSpec *ast.ImportSpec) {
	name := con.getImportName(iSpec)
	if name == `_` {
		// Don't add reference to blanked import.
		return
	}
	// If name is `.` for an anomalous dot import, use it as is.

	path := con.getImportPath(iSpec)
	im := con.cPkg.ImportForPath(path)
	if older, exists := con.fImports.TryGet(name); exists && im != older {
		panic(terror.New(`import name already used`).
			With(`name`, name).
			With(`older`, older).
			With(`newer`, im))
	}
	con.fImports.Add(name, im)
}

func (con *converter) addTypeSpec(tSpec *ast.TypeSpec) {

	// TODO: Implement
}

func (con *converter) addValueSpec(vSpec *ast.ValueSpec) {
	// TODO: Implement
}
