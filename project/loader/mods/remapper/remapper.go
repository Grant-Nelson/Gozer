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

	// TODO: Add a check if the file needs remapping before remapping.

	p := f.TempFileSet.Position(f.File.FileStart)
	tf := f.TempFileSet.FileSet().File(f.File.FileStart)
	base := finalFileSet.FileSet().Base()
	r := &fileRemapper{
		f:         f,
		tokenFile: tokenFileGen.New(p.Filename, base),
		baseShift: int(tf.Base()) - base,
	}
	r.remapFile()
	r.finished(finalFileSet)
	return nil
}

type fileRemapper struct {
	f         *artifacts.File
	tokenFile *tokenFileGen.TokenFileGen
	baseShift int
	needsInfo bool
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

func (r *fileRemapper) addInfo(pos token.Pos, doc **ast.CommentGroup) {
	if !r.needsInfo && int(pos)-r.baseShift != int(r.tokenFile.Current()) {
		r.needsInfo = true
	}

	if r.needsInfo {
		p := r.f.TempFileSet.Position(pos)
		text := fmt.Sprintf(`//line %s:%d:%d`, p.Filename, p.Line, p.Column)
		lineCmt := &ast.Comment{
			Slash: r.tokenFile.Current(),
			Text:  text,
		}

		cg := *doc
		if cg == nil {
			cg = &ast.CommentGroup{}
		}
		cg.List = append(cg.List, lineCmt)
		*doc = cg

		r.tokenFile.Add(len(text) + 1)
		r.tokenFile.AddInfo(p.Filename, p.Line, p.Column)
	}
}

func (r *fileRemapper) getCommentMap(n ast.Node) map[int]*token.Pos {
	commentNodes := map[int]*token.Pos{}
	for pt := range artifacts.WalkPos(n) {
		switch pt.Node.(type) {
		case *ast.Comment:
			prev := r.f.TempFileSet.FindPrevious(*pt.Pos)
			commentNodes[int(prev)] = pt.Pos
		}
	}
	return commentNodes
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
	// Add a newline if one hasn't been added
	if !r.tokenFile.HadLine() {
		r.tokenFile.Add(1)
		r.tokenFile.AddLine()
		r.tokenFile.Add(1)
		r.tokenFile.AddLine()
	}

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
	r.addInfo(d.Pos(), &d.Doc)
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
	r.addInfo(s.Pos(), &s.Doc)
	cmt := r.getCommentMap(s.Comment)
	r.remapBranch(s, cmt)
	if s.EndPos.IsValid() {
		s.EndPos = r.tokenFile.Current()
	}
}

func (r *fileRemapper) remapTypeSpec(s *ast.TypeSpec) {
	r.remapCommentGroup(s.Doc)
	r.addInfo(s.Pos(), &s.Doc)
	cmt := r.getCommentMap(s.Comment)
	r.remapBranch(s, cmt)
}

func (r *fileRemapper) remapValueSpec(s *ast.ValueSpec) {
	r.remapCommentGroup(s.Doc)
	r.addInfo(s.Pos(), &s.Doc)
	cmt := r.getCommentMap(s.Comment)
	r.remapBranch(s, cmt)
}

func (r *fileRemapper) remapFuncDecl(d *ast.FuncDecl) {
	r.remapCommentGroup(d.Doc)
	r.addInfo(d.Pos(), &d.Doc)
	cmt := r.getCommentMap(d)
	r.remapBranch(d, cmt)
}

func (r *fileRemapper) remapBranch(n ast.Node, cmt map[int]*token.Pos) {
	for pt := range artifacts.WalkPos(n) {
		// Skip comments directly
		if _, ok := pt.Node.(*ast.Comment); ok {
			continue
		}

		p := int(*pt.Pos)
		r.remapPos(pt.Pos)

		// Add any comments needed for the current node.
		for {
			c, ok := cmt[p]
			if !ok {
				break
			}
			p = int(*c)
			r.remapPos(c)
		}
	}
}
