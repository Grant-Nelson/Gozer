package golang

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/dictionary"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/Snow-Gremlin/Gozer/constructs"
)

type converter struct {
	p        constructs.CPackage
	pkg      *build.Package
	fSet     *token.FileSet
	proj     constructs.CProject
	fImports collections.Dictionary[string, constructs.CPackage]
}

func newConverter(p constructs.CPackage, pkg *build.Package, fSet *token.FileSet, proj constructs.CProject) *converter {
	return &converter{
		p:        p,
		pkg:      pkg,
		fSet:     fSet,
		proj:     proj,
		fImports: dictionary.New[string, constructs.CPackage](),
	}
}

func (con *converter) convertPackage() {
	enumerator.Enumerate(con.pkg.GoFiles...).Foreach(con.addFile)
}

func (con *converter) addFile(fileName string) {
	fileName = con.prepareFileName(fileName)
	f := con.parseFile(fileName)
	con.addFileNode(f)
}

func (con *converter) prepareFileName(fileName string) string {
	if !filepath.IsAbs(fileName) {
		fileName = filepath.Join(con.pkg.Dir, fileName)
	}
	return fileName
}

func (con *converter) parseFile(fileName string) *ast.File {
	f, err := parser.ParseFile(con.fSet, fileName, nil, parser.ParseComments)
	if err != nil {
		panic(terror.New(`error parsing file`, err).
			With(`file name`, fileName))
	}
	return f
}

func (con *converter) addFileNode(f *ast.File) {
	con.fImports.Clear()
	enumerator.Enumerate(f.Decls...).Foreach(con.addDecl)
}

func (con *converter) pos(p token.Pos) string {
	return con.fSet.Position(p).String()
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

var directiveMatch = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^\/(?:\/|\*)\w+:\w+.*$`)
})

func (con *converter) readDirectives(comments ...*ast.CommentGroup) []string {
	found := []string{}
	for _, cs := range comments {
		if cs != nil {
			for _, c := range cs.List {
				if directiveMatch().MatchString(c.Text) {
					found = append(found, c.Text)
				}
			}
		}
	}
	return found
}

func funcReceiverIdent(funcDecl *ast.FuncDecl) string {
	if funcDecl == nil || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		return ``
	}
	recv := funcDecl.Recv.List[0].Type
	for {
		switch r := recv.(type) {
		case *ast.IndexListExpr:
			recv = r.X
		case *ast.IndexExpr:
			recv = r.X
		case *ast.StarExpr:
			recv = r.X
		case *ast.Ident:
			return r.Name
		default:
			panic(terror.New(`unexpected type in receiver of function`).
				With(`type`, recv).
				With(`function`, funcDecl))
		}
	}
}

func (con *converter) addFunc(funcDecl *ast.FuncDecl) {
	/*
		m := cMethod.New()
		m.SetName(funcDecl.Name.Name)
		m.Directives = con.readDirectives(funcDecl.Doc)
		recvId := funcReceiverIdent(funcDecl)
	*/

	// TODO: Implement

}

func (con *converter) addGenDecl(gDecl *ast.GenDecl) {
	directives := con.readDirectives(gDecl.Doc)
	for _, spec := range gDecl.Specs {
		con.addSpec(spec, directives)
	}
}

func (con *converter) addSpec(spec ast.Spec, declDirectives []string) {
	switch tSpec := spec.(type) {
	case *ast.ImportSpec:
		con.addImportSpec(tSpec, declDirectives)
	case *ast.TypeSpec:
		con.addTypeSpec(tSpec, declDirectives)
	case *ast.ValueSpec:
		con.addValueSpec(tSpec, declDirectives)
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
	other, has := con.proj.Packages().TryGetByPath(path)
	if !has {
		panic(terror.New(`failed to find package for package name`).
			With(`path`, path).
			With(`from`, con.pos(iSpec.Pos())).
			With(`to`, con.pos(iSpec.End())))
	}
	return other.Name()
}

func (con *converter) addImportSpec(iSpec *ast.ImportSpec, declDirectives []string) {
	//directives := con.readDirectives(iSpec.Doc, iSpec.Comment)
	name := con.getImportName(iSpec)
	if name == `_` {
		// Don't add reference to blanked import.
		return
	}
	// If name is `.` for an anomalous dot import, use it as is.

	path := con.getImportPath(iSpec)
	other, has := con.proj.Packages().TryGetByPath(path)
	if !has {
		panic(terror.New(`failed to find package for import path`).
			With(`path`, path).
			With(`from`, con.pos(iSpec.Pos())).
			With(`to`, con.pos(iSpec.End())))
	}
	if older, exists := con.fImports.TryGet(name); exists && other != older {
		panic(terror.New(`import name already used`).
			With(`name`, name).
			With(`older`, older).
			With(`newer`, other))
	}
	con.fImports.Add(name, other)
}

func (con *converter) addTypeSpec(tSpec *ast.TypeSpec, declDirectives []string) {
	//directives := con.readDirectives(tSpec.Doc, tSpec.Comment)
	name := tSpec.Name.Name
	//aliased := tSpec.Assign == token.NoPos

	fmt.Printf("%s: %+v\n", name, tSpec)

	// TODO: Implement
}

func (con *converter) addValueSpec(vSpec *ast.ValueSpec, declDirectives []string) {
	//directives := con.readDirectives(vSpec.Doc, vSpec.Comment)

	// TODO: Implement

}
