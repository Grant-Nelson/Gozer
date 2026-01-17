package artifacts

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_WalkPos_Package(t *testing.T) {
	f, _ := loadTest(t,
		`// comment 1`,
		``,
		`// comment 2`,
		`package foo`,
		`// comment 3`,
		``,
		`// comment 4`)
	checkWalkPos(t, f)(
		`1:File.Start:0`,
		`1:File.Comment:12"// comment 1"`,
		`15:File.Doc:12"// comment 2"`,
		`28:File.Package:7"package"`,
		`36:Ident.Name:3"foo"`,
		`40:File.Comment:12"// comment 3"`,
		`54:File.Comment:12"// comment 4"`,
		`66:File.End:0`)
}

func Test_WalkPos_MultilineFunc(t *testing.T) {
	f, _ := loadTest(t,
		`package foo`,
		`// comment 1`,
		``,
		`// comment 2`,
		`func Foo( // comment 3`,
		`	a []int, // comment 4`,
		`	b ...string, // comment 5`,
		`)( // comment 6`,
		`	c bool, // comment 7`,
		``,
		`// comment 8`,
		`){`,
		`	// comment 9`,
		`	return false // comment 10`,
		`	// comment 11`,
		`} // comment 12`)
	checkWalkPos(t, f)(
		`1:File.Start:0`,
		`1:File.Package:7"package"`,
		`9:Ident.Name:3"foo"`,
		`13:File.Comment:12"// comment 1"`,
		`27:FuncDecl.Doc:12"// comment 2"`,
		`40:FuncType.Func:4"func"`,
		`45:Ident.Name:3"Foo"`,
		`48:FieldList.Opening:1"("`,
		`50:File.Comment:12"// comment 3"`,
		`64:Ident.Name:1"a"`,
		`66:ArrayType.Lbrack:1"["`,
		`68:Ident.Name:3"int"`,
		`73:File.Comment:12"// comment 4"`,
		`87:Ident.Name:1"b"`,
		`89:X.Ellipsis:3"..."`,
		`92:Ident.Name:6"string"`,
		`100:File.Comment:12"// comment 5"`,
		`113:FieldList.Closing:1")"`,
		`114:FieldList.Opening:1"("`,
		`116:File.Comment:12"// comment 6"`,
		`130:Ident.Name:1"c"`,
		`132:Ident.Name:4"bool"`,
		`138:File.Comment:12"// comment 7"`,
		`152:File.Comment:12"// comment 8"`,
		`165:FieldList.Closing:1")"`,
		`166:BlockStmt.Lbrace:1"{"`,
		`169:File.Comment:12"// comment 9"`,
		`183:ReturnStmt.Return:6"return"`,
		`190:Ident.Name:5"false"`,
		`196:File.Comment:13"// comment 10"`,
		`211:File.Comment:13"// comment 11"`,
		`225:BlockStmt.Rbrace:1"}"`,
		`227:File.Comment:13"// comment 12"`,
		`240:File.End:0`)
	checkWalkPos(t, f, SkipFileComments)(
		`1:File.Start:0`,
		`1:File.Package:7"package"`,
		`9:Ident.Name:3"foo"`,
		`27:FuncDecl.Doc:12"// comment 2"`,
		`40:FuncType.Func:4"func"`,
		`45:Ident.Name:3"Foo"`,
		`48:FieldList.Opening:1"("`,
		`64:Ident.Name:1"a"`,
		`66:ArrayType.Lbrack:1"["`,
		`68:Ident.Name:3"int"`,
		`87:Ident.Name:1"b"`,
		`89:X.Ellipsis:3"..."`,
		`92:Ident.Name:6"string"`,
		`113:FieldList.Closing:1")"`,
		`114:FieldList.Opening:1"("`,
		`130:Ident.Name:1"c"`,
		`132:Ident.Name:4"bool"`,
		`165:FieldList.Closing:1")"`,
		`166:BlockStmt.Lbrace:1"{"`,
		`183:ReturnStmt.Return:6"return"`,
		`190:Ident.Name:5"false"`,
		`225:BlockStmt.Rbrace:1"}"`,
		`240:File.End:0`)
}

func Test_WalkPos_Values_Arrays(t *testing.T) {
	f, _ := loadTest(t,
		`package foo`,
		``,
		`var (`,
		`	a []int`,
		`	b [   ]  int`,
		`	c [...]int`,
		`	d [42]int`,
		`)`)
	checkWalkPos(t, f)(
		`1:File.Start:0`,
		`1:File.Package:7"package"`,
		`9:Ident.Name:3"foo"`,
		`14:GenDecl.Tok:3"var"`,
		`18:GenDecl.Lparen:1"("`,
		// a []int
		`21:Ident.Name:1"a"`,
		`23:ArrayType.Lbrack:1"["`,
		`25:Ident.Name:3"int"`,
		// b [   ]  int
		`30:Ident.Name:1"b"`,
		`32:ArrayType.Lbrack:1"["`,
		`39:Ident.Name:3"int"`,
		// c [...]int
		`44:Ident.Name:1"c"`,
		`46:ArrayType.Lbrack:1"["`,
		`47:X.Ellipsis:3"..."`,
		`51:Ident.Name:3"int"`,
		// d [42]int
		`56:Ident.Name:1"d"`,
		`58:ArrayType.Lbrack:1"["`,
		`59:BasicLit.Value:2"42"`,
		`62:Ident.Name:3"int"`,
		`66:GenDecl.Rparen:1")"`,
		`67:File.End:0`)
	checkWalkPos(t, f, AddPseudoNodes)(
		`1:File.Start:0`,
		`1:File.Package:7"package"`,
		`9:Ident.Name:3"foo"`,
		`14:GenDecl.Tok:3"var"`,
		`18:GenDecl.Lparen:1"("`,
		// a []int
		`21:Ident.Name:1"a"`,
		`23:ArrayType.Lbrack:1"["`,
		`24:ArrayType.Rbrack:1"]"(P)`,
		`25:Ident.Name:3"int"`,
		// b [   ]  int
		`30:Ident.Name:1"b"`,
		`32:ArrayType.Lbrack:1"["`,
		`33:ArrayType.Rbrack:1"]"(P)`,
		`39:Ident.Name:3"int"`,
		// c [...]int
		`44:Ident.Name:1"c"`,
		`46:ArrayType.Lbrack:1"["`,
		`47:X.Ellipsis:3"..."`,
		`50:ArrayType.Rbrack:1"]"(P)`,
		`51:Ident.Name:3"int"`,
		// d [42]int
		`56:Ident.Name:1"d"`,
		`58:ArrayType.Lbrack:1"["`,
		`59:BasicLit.Value:2"42"`,
		`61:ArrayType.Rbrack:1"]"(P)`,
		`62:Ident.Name:3"int"`,
		`66:GenDecl.Rparen:1")"`,
		`67:File.End:0`)
}

func Test_WalkPos_Values_Comments(t *testing.T) {
	f, _ := loadTest(t,
		`package foo`,
		``,
		`// comment 1`,
		`var a int // comment 2`,
		``,
		`// comment 3`,
		`var b, // comment 4`,
		`	c int`,
		``,
		`// comment 5`,
		`var (`,
		`	// comment 6`,
		`	d int // comment 7`,
		`	`,
		`	// comment 8`,
		`	e int // comment 9`,
		`)`)
	checkWalkPos(t, f)(
		`1:File.Start:0`,
		`1:File.Package:7"package"`,
		`9:Ident.Name:3"foo"`,
		`14:GenDecl.Doc:12"// comment 1"`,
		`27:GenDecl.Tok:3"var"`,
		`31:Ident.Name:1"a"`,
		`33:Ident.Name:3"int"`,
		`37:ValueSpec.Comment:12"// comment 2"`,
		`51:GenDecl.Doc:12"// comment 3"`,
		`64:GenDecl.Tok:3"var"`,
		`68:Ident.Name:1"b"`,
		`71:File.Comment:12"// comment 4"`,
		`85:Ident.Name:1"c"`,
		`87:Ident.Name:3"int"`,
		`92:GenDecl.Doc:12"// comment 5"`,
		`105:GenDecl.Tok:3"var"`,
		`109:GenDecl.Lparen:1"("`,
		`112:ValueSpec.Doc:12"// comment 6"`,
		`126:Ident.Name:1"d"`,
		`128:Ident.Name:3"int"`,
		`132:ValueSpec.Comment:12"// comment 7"`,
		`148:ValueSpec.Doc:12"// comment 8"`,
		`162:Ident.Name:1"e"`,
		`164:Ident.Name:3"int"`,
		`168:ValueSpec.Comment:12"// comment 9"`,
		`181:GenDecl.Rparen:1")"`,
		`182:File.End:0`)
}

func Test_WalkPos_Struct(t *testing.T) {
	f, _ := loadTest(t,
		`package foo`,
		`// comment 1`,
		`type Foo struct { // comment 2`,
		``,
		`	// comment 3`,
		"	x int `json:\"-\"` // comment 4",
		``,
		`	// comment 5`,
		``,
		`	// comment 6`,
		`	y, // comment 7`,
		`	z int // comment 8`,
		``,
		`	// comment 9`,
		`}`)
	checkWalkPos(t, f)(
		`1:File.Start:0`,
		`1:File.Package:7"package"`,
		`9:Ident.Name:3"foo"`,
		`13:GenDecl.Doc:12"// comment 1"`,
		// type Foo struct {
		`26:GenDecl.Tok:4"type"`,
		`31:Ident.Name:3"Foo"`,
		`35:StructType.Struct:6"struct"`,
		`42:FieldList.Opening:1"{"`,
		`44:File.Comment:12"// comment 2"`,
		`59:Field.Doc:12"// comment 3"`,
		// x int `json:\"-\"`
		`73:Ident.Name:1"x"`,
		`75:Ident.Name:3"int"`,
		"79:BasicLit.Value:10\"`json:\\\"-\\\"`\"",
		`90:Field.Comment:12"// comment 4"`,
		`105:File.Comment:12"// comment 5"`,
		`120:Field.Doc:12"// comment 6"`,
		// y,
		`134:Ident.Name:1"y"`,
		`137:File.Comment:12"// comment 7"`,
		// z int
		`151:Ident.Name:1"z"`,
		`153:Ident.Name:3"int"`,
		`157:Field.Comment:12"// comment 8"`,
		`172:File.Comment:12"// comment 9"`,
		`185:FieldList.Closing:1"}"`,
		`186:File.End:0`)
}

func Test_WalkPos_Channels(t *testing.T) {
	f, _ := loadTest(t,
		`package foo`,
		`func Foo(src <-chan int, dst chan<- int, notUsed chan int) {`,
		`	dst <- src`,
		`}`)
	checkWalkPos(t, f)(
		`1:File.Start:0`,
		`1:File.Package:7"package"`,
		`9:Ident.Name:3"foo"`,
		`13:FuncType.Func:4"func"`,
		`18:Ident.Name:3"Foo"`,
		`21:FieldList.Opening:1"("`,
		// src <-chan int
		`22:Ident.Name:3"src"`,
		`26:ArrowChan.Arrow:0`,
		`26:ArrowChan.Chan:6"<-chan"`,
		`33:Ident.Name:3"int"`,
		// dst chan<- int
		`38:Ident.Name:3"dst"`,
		`42:ChanArrow.Chan:4"chan"`,
		`46:ChanArrow.Arrow:2"<-"`,
		`49:Ident.Name:3"int"`,
		// notUsed chan int
		`54:Ident.Name:7"notUsed"`,
		`62:Chan.Chan:4"chan"`,
		`67:Ident.Name:3"int"`,
		`70:FieldList.Closing:1")"`,
		`72:BlockStmt.Lbrace:1"{"`,
		// dst <- src
		`75:Ident.Name:3"dst"`,
		`79:SendStmt.Arrow:2"<-"`,
		`82:Ident.Name:3"src"`,
		`86:BlockStmt.Rbrace:1"}"`,
		`87:File.End:0`)
}

func loadTest(t testing.TB, code ...string) (*ast.File, *token.FileSet) {
	t.Helper()
	const mode = parser.AllErrors |
		parser.ParseComments |
		parser.DeclarationErrors |
		parser.SkipObjectResolution
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, `test.go`, strings.Join(code, "\n"), mode)
	if err != nil {
		t.Fatalf(`failed to load test file: %v`, err)
	}
	return f, fs
}

func checkWalkPos(t testing.TB, f *ast.File, options ...WalkPosOption) func(expLines ...string) {
	t.Helper()
	lines := []string{}
	var prior int
	for pt := range WalkPos(f, options...) {
		if !pt.Pos.IsValid() {
			t.Errorf("invalid position returned for %s", pt.String())
		}
		if prior > int(*pt.Pos) {
			t.Errorf("out-of-order positions returned for %s", pt.String())
		}
		lines = append(lines, pt.String())
		prior = int(*pt.Pos)
	}
	return func(expLines ...string) {
		if diff := cmp.Diff(expLines, lines); len(diff) > 0 {
			t.Errorf("the line for WalkPos didn't match expected lines:\n%s", diff)
		}
	}
}
