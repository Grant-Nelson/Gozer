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
	fmt.Printf(">> FLAG: finish (%d)\n\n", r.tokenFile.Current()) // TODO: REMOVE

	r.tokenFile.Write(finalFileSet.FileSet())
	r.f.TempFileSet = finalFileSet
	finalFileSet.RegisterFile(r.f)

	//ast.Print(finalFileSet.FileSet(), r.f.File)   // TODO: REMOVE
	fmt.Printf("---------------------------\n")   // TODO: REMOVE
	for pt := range artifacts.WalkPos(r.f.File) { // TODO: REMOVE
		fmt.Printf(">> %s\n", pt.String()) // TODO: REMOVE
	}
	fmt.Printf("\n") // TODO: REMOVE
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
	fmt.Printf(">> FLAG: Decl (%d)\n", r.tokenFile.Current()) // TODO: REMOVE

	// Add a newline if one hasn't been added
	if !r.tokenFile.HadLine() {
		fmt.Printf(">> FLAG: Add new line (%d)\n", r.tokenFile.Current()) // TODO: REMOVE
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
	fmt.Printf(">> FLAG: GenDecl (%d)\n", r.tokenFile.Current()) // TODO: REMOVE
	r.remapCommentGroup(d.Doc)
	r.remapPos(&d.TokPos)
	r.remapPos(&d.Lparen)
	for _, s := range d.Specs {
		r.remapSpec(s)
	}
	r.remapPos(&d.Rparen)
}

func (r *fileRemapper) remapSpec(s ast.Spec) {
	fmt.Printf(">> FLAG: Spec (%d)\n", r.tokenFile.Current()) // TODO: REMOVE
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
	fmt.Printf(">> FLAG: ImportSpec (%d)\n", r.tokenFile.Current()) // TODO: REMOVE
	r.remapCommentGroup(s.Doc)
	cmt := r.getCommentMap(s.Comment)
	r.remapBranch(s, cmt)
	if s.EndPos.IsValid() {
		s.EndPos = r.tokenFile.Current()
	}
}

func (r *fileRemapper) remapTypeSpec(s *ast.TypeSpec) {
	fmt.Printf(">> FLAG: TypeSpec (%d)\n", r.tokenFile.Current()) // TODO: REMOVE
	r.remapCommentGroup(s.Doc)
	cmt := r.getCommentMap(s.Comment)
	r.remapBranch(s, cmt)
}

func (r *fileRemapper) remapValueSpec(s *ast.ValueSpec) {
	fmt.Printf(">> FLAG: ValueSpec (%d)\n", r.tokenFile.Current()) // TODO: REMOVE
	r.remapCommentGroup(s.Doc)
	cmt := r.getCommentMap(s.Comment)
	r.remapBranch(s, cmt)
}

func (r *fileRemapper) remapFuncDecl(d *ast.FuncDecl) {
	fmt.Printf(">> FLAG: FuncDecl (%d)\n", r.tokenFile.Current()) // TODO: REMOVE
	r.remapCommentGroup(d.Doc)
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
