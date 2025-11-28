package file

import (
	"fmt"
	"go/ast"
	"go/token"
)

type fileRemapper struct {
	f          *File
	offset     int
	frames     []token.Position
	priorPos   token.Position
	priorWidth int
}

// Remap will rewrite the file to a new file set to normalize the file information.
// This is required to be done prior to writing the file so that the file
// will output correctly.
func (f *File) Remap() {
	frm := &fileRemapper{f: f}
	frm.walk(f.File)
	frm.finish()
}

func (frm *fileRemapper) finish() {
	fileSet := token.NewFileSet()
	p := frm.f.FileSet.Position(frm.f.File.FileStart)
	f := fileSet.AddFile(p.Filename, 0, frm.offset)
	for _, frame := range frm.frames {
		f.AddLineColumnInfo(frame.Offset, frame.Filename, frame.Line, frame.Column)
	}
	frm.f.FileSet = fileSet
}

func (frm *fileRemapper) mapTokPos(n ast.Node, pos *token.Pos, tok token.Token) {
	frm.mapPos(n, pos, len(tok.String()))
}

func (frm *fileRemapper) mapPos(n ast.Node, pos *token.Pos, width int) {
	if !pos.IsValid() {
		return
	}

	p := frm.f.FileSet.Position(*pos)
	if p.Filename != frm.priorPos.Filename {
		frm.priorPos = p
		frm.priorWidth = width
		*pos = token.Pos(frm.offset)
		frm.frames = append(frm.frames, p)
		return
	}

	//p.Line
	//p.Column
	//p.Offset
	//p.Filename

	*pos = token.Pos(frm.offset)
	// TODO: Implement
}

func mapNodeList[N ast.Node](frm *fileRemapper, list []N) {
	for _, node := range list {
		frm.walk(node)
	}
}

func (frm *fileRemapper) walk(node ast.Node) {
	switch n := node.(type) {
	case nil:
		// do nothing

	// Comments and fields
	case *ast.Comment:
		frm.mapPos(n, &n.Slash, len(n.Text))

	case *ast.CommentGroup:
		mapNodeList(frm, n.List)

	case *ast.Field:
		frm.walk(n.Doc)
		mapNodeList(frm, n.Names)
		frm.walk(n.Type)
		frm.walk(n.Tag)
		frm.walk(n.Comment)

	case *ast.FieldList:
		frm.mapTokPos(n, &n.Opening, token.RBRACE)
		mapNodeList(frm, n.List)
		frm.mapTokPos(n, &n.Closing, token.LBRACE)

	// Expressions
	case *ast.BadExpr:
		frm.mapPos(n, &n.From, int(n.To)-int(n.From))
		frm.mapPos(n, &n.To, 1)

	case *ast.Ident:
		frm.mapPos(n, &n.NamePos, len(n.Name))

	case *ast.Ellipsis:
		frm.mapTokPos(n, &n.Ellipsis, token.ELLIPSIS)
		frm.walk(n.Elt)

	case *ast.BasicLit:
		frm.mapPos(n, &n.ValuePos, len(n.Value))

	case *ast.FuncLit:
		frm.walk(n.Type)
		frm.walk(n.Body)

	case *ast.CompositeLit:
		frm.walk(n.Type)
		frm.mapTokPos(n, &n.Lbrace, token.LBRACE)
		mapNodeList(frm, n.Elts)
		frm.mapTokPos(n, &n.Rbrace, token.RBRACE)

	case *ast.ParenExpr:
		frm.mapTokPos(n, &n.Lparen, token.LPAREN)
		frm.walk(n.X)
		frm.mapTokPos(n, &n.Rparen, token.RPAREN)

	case *ast.SelectorExpr:
		frm.walk(n.X)
		frm.walk(n.Sel)

	case *ast.IndexExpr:
		frm.walk(n.X)
		frm.mapTokPos(n, &n.Lbrack, token.LBRACK)
		frm.walk(n.Index)
		frm.mapTokPos(n, &n.Rbrack, token.RBRACK)

	case *ast.IndexListExpr:
		frm.walk(n.X)
		frm.mapTokPos(n, &n.Lbrack, token.LBRACK)
		mapNodeList(frm, n.Indices)
		frm.mapTokPos(n, &n.Rbrack, token.RBRACK)

	case *ast.SliceExpr:
		frm.walk(n.X)
		frm.mapTokPos(n, &n.Lbrack, token.LBRACK)
		frm.walk(n.Low)
		frm.walk(n.High)
		frm.walk(n.Max)
		frm.mapTokPos(n, &n.Rbrack, token.RBRACK)

	case *ast.TypeAssertExpr:
		frm.walk(n.X)
		frm.mapTokPos(n, &n.Lparen, token.LPAREN)
		frm.walk(n.Type)
		frm.mapTokPos(n, &n.Rparen, token.RPAREN)

	case *ast.CallExpr:
		frm.walk(n.Fun)
		frm.mapTokPos(n, &n.Lparen, token.LPAREN)
		mapNodeList(frm, n.Args)
		// TODO: Determine where the ellipsis needs to go
		frm.mapTokPos(n, &n.Ellipsis, token.ELLIPSIS)
		frm.mapTokPos(n, &n.Rparen, token.RPAREN)

	case *ast.StarExpr:
		frm.mapTokPos(n, &n.Star, token.MUL)
		frm.walk(n.X)

	case *ast.UnaryExpr:
		frm.mapTokPos(n, &n.OpPos, n.Op)
		frm.walk(n.X)

	case *ast.BinaryExpr:
		frm.walk(n.X)
		frm.mapTokPos(n, &n.OpPos, n.Op)
		frm.walk(n.Y)

	case *ast.KeyValueExpr:
		frm.walk(n.Key)
		frm.mapTokPos(n, &n.Colon, token.COLON)
		frm.walk(n.Value)

	// Types
	case *ast.ArrayType:
		frm.mapTokPos(n, &n.Lbrack, token.LBRACK)
		// TODO: Determine why no Rbrack?
		frm.walk(n.Len)
		frm.walk(n.Elt)

	case *ast.StructType:
		frm.mapTokPos(n, &n.Struct, token.STRUCT)
		frm.walk(n.Fields)

	case *ast.FuncType:
		frm.mapTokPos(n, &n.Func, token.FUNC)
		frm.walk(n.TypeParams)
		frm.walk(n.Params)
		frm.walk(n.Results)

	case *ast.InterfaceType:
		frm.mapTokPos(n, &n.Interface, token.INTERFACE)
		frm.walk(n.Methods)

	case *ast.MapType:
		frm.mapTokPos(n, &n.Map, token.MAP)
		frm.walk(n.Key)
		frm.walk(n.Value)

	case *ast.ChanType:
		// TODO: Check the order is correct for the Begin and Arrow
		frm.mapTokPos(n, &n.Begin, token.CHAN)
		frm.mapTokPos(n, &n.Arrow, token.ARROW)
		frm.walk(n.Value)

	// Statements
	case *ast.BadStmt:
		frm.mapPos(n, &n.From, int(n.To)-int(n.From))
		frm.mapPos(n, &n.To, 1)

	case *ast.DeclStmt:
		frm.walk(n.Decl)

	case *ast.EmptyStmt:
		frm.mapTokPos(n, &n.Semicolon, token.SEMICOLON)

	case *ast.LabeledStmt:
		frm.walk(n.Label)
		frm.mapTokPos(n, &n.Colon, token.COLON)
		frm.walk(n.Stmt)

	case *ast.ExprStmt:
		frm.walk(n.X)

	case *ast.SendStmt:
		frm.walk(n.Chan)
		frm.mapTokPos(n, &n.Arrow, token.ARROW)
		frm.walk(n.Value)

	case *ast.IncDecStmt:
		frm.walk(n.X)
		frm.mapTokPos(n, &n.TokPos, n.Tok)

	case *ast.AssignStmt:
		mapNodeList(frm, n.Lhs)
		frm.mapTokPos(n, &n.TokPos, n.Tok)
		mapNodeList(frm, n.Rhs)

	case *ast.GoStmt:
		frm.mapTokPos(n, &n.Go, token.GO)
		frm.walk(n.Call)

	case *ast.DeferStmt:
		frm.mapTokPos(n, &n.Defer, token.DEFER)
		frm.walk(n.Call)

	case *ast.ReturnStmt:
		frm.mapTokPos(n, &n.Return, token.RETURN)
		mapNodeList(frm, n.Results)

	case *ast.BranchStmt:
		frm.mapTokPos(n, &n.TokPos, n.Tok)
		frm.walk(n.Label)

	case *ast.BlockStmt:
		frm.mapTokPos(n, &n.Lbrace, token.LBRACE)
		mapNodeList(frm, n.List)
		frm.mapTokPos(n, &n.Rbrace, token.RBRACE)

	case *ast.IfStmt:
		frm.mapTokPos(n, &n.If, token.IF)
		frm.walk(n.Init)
		frm.walk(n.Cond)
		frm.walk(n.Body)
		frm.walk(n.Else)

	case *ast.CaseClause:
		frm.mapTokPos(n, &n.Case, token.CASE)
		mapNodeList(frm, n.List)
		frm.mapTokPos(n, &n.Colon, token.COLON)
		mapNodeList(frm, n.Body)

	case *ast.SwitchStmt:
		frm.mapTokPos(n, &n.Switch, token.SWITCH)
		frm.walk(n.Init)
		frm.walk(n.Tag)
		frm.walk(n.Body)

	case *ast.TypeSwitchStmt:
		frm.mapTokPos(n, &n.Switch, token.SWITCH)
		frm.walk(n.Init)
		frm.walk(n.Assign)
		frm.walk(n.Body)

	case *ast.CommClause:
		frm.mapTokPos(n, &n.Case, token.CASE)
		frm.walk(n.Comm)
		frm.mapTokPos(n, &n.Colon, token.COLON)
		mapNodeList(frm, n.Body)

	case *ast.SelectStmt:
		frm.mapTokPos(n, &n.Select, token.SELECT)
		frm.walk(n.Body)

	case *ast.ForStmt:
		frm.mapTokPos(n, &n.For, token.FOR)
		frm.walk(n.Init)
		frm.walk(n.Cond)
		frm.walk(n.Post)
		frm.walk(n.Body)

	case *ast.RangeStmt:
		frm.mapTokPos(n, &n.For, token.FOR)
		frm.walk(n.Key)
		frm.walk(n.Value)
		frm.mapTokPos(n, &n.TokPos, n.Tok)
		frm.mapTokPos(n, &n.Range, token.RANGE)
		frm.walk(n.X)
		frm.walk(n.Body)

	// Declarations
	case *ast.ImportSpec:
		frm.walk(n.Doc)
		frm.walk(n.Name)
		frm.walk(n.Path)
		frm.walk(n.Comment)
		frm.mapTokPos(n, &n.EndPos, 0)

	case *ast.ValueSpec:
		frm.walk(n.Doc)
		mapNodeList(frm, n.Names)
		frm.walk(n.Type)
		mapNodeList(frm, n.Values)
		frm.walk(n.Comment)

	case *ast.TypeSpec:
		frm.walk(n.Doc)
		frm.walk(n.Name)
		frm.walk(n.TypeParams)
		frm.mapTokPos(n, &n.Assign, token.EQL)
		frm.walk(n.Type)
		frm.walk(n.Comment)

	case *ast.BadDecl:
		frm.mapPos(n, &n.From, int(n.To)-int(n.From))
		frm.mapPos(n, &n.To, 0)

	case *ast.GenDecl:
		frm.walk(n.Doc)
		frm.mapTokPos(n, &n.TokPos, n.Tok)
		frm.mapTokPos(n, &n.Lparen, token.LPAREN)
		mapNodeList(frm, n.Specs)
		frm.mapTokPos(n, &n.Rparen, token.RPAREN)

	case *ast.FuncDecl:
		frm.walk(n.Doc)
		frm.walk(n.Recv)
		frm.walk(n.Name)
		frm.walk(n.Type)
		frm.walk(n.Body)

	// Files
	case *ast.File:
		frm.mapPos(n, &n.FileStart, 0)
		frm.walk(n.Doc)
		frm.mapTokPos(n, &n.Package, token.PACKAGE)
		frm.walk(n.Name)
		mapNodeList(frm, n.Decls)
		frm.mapPos(n, &n.FileEnd, 0)

	default:
		panic(fmt.Errorf(`unexpected node in mapPos: (%[1]T) %[1]v`, n))
	}
}
