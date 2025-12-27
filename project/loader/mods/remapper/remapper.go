package remapper

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/remapper/tokenFileGen"
)

// Remap will rewrite the file to a new file set to normalize the file information.
// This is required to be done prior to writing the file so that the file
// will output correctly.
func Remap(f *artifacts.File, finalFileSet *artifacts.FileSet, errGroup *faults.Group) error {
	p := f.TempFileSet.Position(f.File.FileStart)
	base := finalFileSet.FileSet().Base()
	r := &fileRemapper{
		f:         f,
		tokenFile: tokenFileGen.New(p.Filename, base),
	}
	r.remapFile()
	r.finished(finalFileSet)
	return nil
}

type fileRemapper struct {
	f         *artifacts.File
	tokenFile *tokenFileGen.TokenFileGen
}

func (r *fileRemapper) finished(finalFileSet *artifacts.FileSet) {
	r.tokenFile.Write(finalFileSet.FileSet())
	r.f.TempFileSet = finalFileSet
	finalFileSet.RegisterFile(r.f)
}

func (r *fileRemapper) remapFile() {
	file := r.f.File
	file.FileStart = r.tokenFile.Current()
	r.remapCommentGroup(file.Doc)
	r.remapPos(&file.Package)
	r.remapPos(&file.Name.NamePos)
	for _, decl := range file.Decls {
		r.remapDecl(decl)
	}
	file.FileEnd = r.tokenFile.Current()
}

func (r *fileRemapper) remapCommentGroup(cg *ast.CommentGroup) {
	if cg != nil {
		for _, c := range cg.List {
			r.remapPos(&c.Slash)
		}
	}
}

func (r *fileRemapper) addInfo(pos token.Pos) {
	//TODO: Implement
}

func (r *fileRemapper) remapPos(pos *token.Pos) {
	if !pos.IsValid() {
		return
	}
	start := r.tokenFile.Current()
	_, lines := r.f.TempFileSet.Widths(*pos)
	for i, line := range lines {
		if i > 0 {
			r.tokenFile.AddLine()
		}
		r.tokenFile.Add(line)
	}
	*pos = start
}

func (r *fileRemapper) remapDecl(d ast.Decl) {
	switch d := d.(type) {
	case *ast.BadDecl:
		r.remapBadDecl(d)
	case *ast.GenDecl:
		r.remapGenDecl(d)
	case *ast.FuncDecl:
		r.remapFuncDecl(d)
	default:
		panic(fmt.Errorf(`unexpected decl type: %T`, d))
	}
}

func (r *fileRemapper) remapBadDecl(d *ast.BadDecl) {
	r.remapPos(&d.From)
	d.To = r.tokenFile.Current()
}

func (r *fileRemapper) remapGenDecl(d *ast.GenDecl) {
	r.remapCommentGroup(d.Doc)
	r.remapPos(&d.TokPos)
	r.remapPos(&d.Lparen)
	for _, s := range d.Specs {
		r.remapSpec(s)
	}
	r.remapPos(&d.Rparen)
}

func (r *fileRemapper) remapSpec(s ast.Spec) {
	switch s := s.(type) {
	case *ast.ImportSpec:
		r.remapImportSpec(s)
	case *ast.TypeSpec:
		r.remapTypeSpec(s)
	case *ast.ValueSpec:
		r.remapValueSpec(s)
	default:
		panic(fmt.Errorf(`unexpected spec type: %T`, s))
	}
}

func (r *fileRemapper) remapImportSpec(s *ast.ImportSpec) {
	r.remapCommentGroup(s.Doc)
	if s.Name != nil {
		r.remapPos(&s.Name.NamePos)
	}
	r.remapPos(&s.Path.ValuePos)
	// assume comments are after path
	r.remapCommentGroup(s.Comment)
	if s.EndPos.IsValid() {
		s.EndPos = r.tokenFile.Current()
	}
}

func (r *fileRemapper) remapTypeSpec(s *ast.TypeSpec) {
	r.remapCommentGroup(s.Doc)

	// TODO: Implement
	// Name       *Ident        // type name
	// TypeParams *FieldList    // type parameters; or nil
	// Assign     token.Pos     // position of '=', if any
	// Type       Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
	// Comment    *CommentGroup // line comments; or nil

}

func (r *fileRemapper) remapValueSpec(s *ast.ValueSpec) {
	r.remapCommentGroup(s.Doc)

	// TODO: Implement
	// Names   []*Ident      // value names (len(Names) > 0)
	// Type    Expr          // value type; or nil
	// Values  []Expr        // initial values; or nil
	// Comment *CommentGroup // line comments; or nil

}

func (r *fileRemapper) remapFuncDecl(d *ast.FuncDecl) {
	r.remapCommentGroup(d.Doc)

	// TODO: Implement
	// Recv *FieldList    // receiver (methods); or nil (functions)
	// Name *Ident        // function/method name
	// Type *FuncType     // function signature: type and value parameters, results, and position of "func" keyword
	// Body *BlockStmt    // function body; or nil for external (non-Go) function

}

func (r *fileRemapper) remapBranch(n ast.Node) {

	// TODO: Implement

}
