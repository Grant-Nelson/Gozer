package artifacts

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_WalkPos_Package(t *testing.T) {
	f := loadTest(t,
		`// comment 1`,
		``,
		`// comment 2`,
		`package foo`,
		`// comment 3`,
		``,
		`// comment 4`)
	checkWalkPos(t, f,
		`1:File:Start`,
		`1:Comment:File.Comment`,
		`15:Comment:File.Doc`,
		`28:File:Package`,
		`36:Ident:foo`,
		`40:Comment:File.Comment`,
		`54:Comment:File.Comment`,
		`66:File:End`)
}

func Test_WalkPos_MultilineFunc(t *testing.T) {
	f := loadTest(t,
		`package foo`,
		`// comment 1`,
		``,
		`// comment 2`,
		`func Foo( // comment 3`,
		`	a []int, // comment 4`,
		`	b ...string, // comment 5`,
		`)( // comment 6`,
		`	c bool, // comment7`,
		``,
		`// comment 8`,
		`){`,
		`	// comment 9`,
		`	return false // comment 10`,
		`	// comment 11`,
		`} // comment 12`)
	checkWalkPos(t, f,
		`1:File:Start`,
		`1:File:Package`,
		`9:Ident:foo`,
		`13:Comment:File.Comment`,
		`27:Comment:FuncDecl.Doc`,
		`40:FuncDecl:Func`,
		`45:Ident:Foo`,
		`48:FieldList:Opening`,
		`50:Comment:File.Comment`,
		`64:Ident:a`,
		`66:ArrayType:LeftBracket`,
		`68:Ident:int`,
		`73:Comment:File.Comment`,
		`87:Ident:b`,
		`89:Ellipsis:Ellipsis`,
		`92:Ident:string`,
		`100:Comment:File.Comment`,
		`113:FieldList:Closing`,
		`114:FieldList:Opening`,
		`116:Comment:File.Comment`,
		`130:Ident:c`,
		`132:Ident:bool`,
		`138:Comment:File.Comment`,
		`151:Comment:File.Comment`,
		`164:FieldList:Closing`,
		`165:BlockStmt:LeftBrace`,
		`168:Comment:File.Comment`,
		`182:ReturnStmt:Return`,
		`189:Ident:false`,
		`195:Comment:File.Comment`,
		`210:Comment:File.Comment`,
		`224:BlockStmt:RightBrace`,
		`226:Comment:File.Comment`,
		`239:File:End`)
}

func Test_WalkPos_Values_Arrays(t *testing.T) {
	f := loadTest(t,
		`package foo`,
		``,
		`var (`,
		`	a []int`,
		`	b [   ]  int`,
		`	c [...]int`,
		`	d [42]int`,
		`)`)
	checkWalkPos(t, f,
		`1:File:Start`,
		`1:File:Package`,
		`9:Ident:foo`,
		`14:GenDecl:var.Pos`,
		`18:GenDecl:var.LeftParen`,
		// a []int
		`21:Ident:a`,
		`23:ArrayType:LeftBracket`,
		`25:Ident:int`,
		// b [   ]  int
		`30:Ident:b`,
		`32:ArrayType:LeftBracket`,
		`39:Ident:int`,
		// c [...]int
		`44:Ident:c`,
		`46:ArrayType:LeftBracket`,
		`47:Ellipsis:Ellipsis`,
		`51:Ident:int`,
		// d [42]int
		`56:Ident:d`,
		`58:ArrayType:LeftBracket`,
		`59:BasicLit:ValuePos`,
		`62:Ident:int`,
		`66:GenDecl:var.RightParen`,
		`67:File:End`)
}

func Test_WalkPos_Values_Comments(t *testing.T) {
	f := loadTest(t,
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
	checkWalkPos(t, f,
		`1:File:Start`,
		`1:File:Package`,
		`9:Ident:foo`,
		`14:Comment:GenDecl.Doc`, // comment 1
		`27:GenDecl:var.Pos`,
		`31:Ident:a`,
		`33:Ident:int`,
		`37:Comment:ValueSpec.Comment`, // comment 2
		`51:Comment:GenDecl.Doc`,       // comment 3
		`64:GenDecl:var.Pos`,
		`68:Ident:b`,
		`71:Comment:File.Comment`, // comment 4
		`85:Ident:c`,
		`87:Ident:int`,
		`92:Comment:GenDecl.Doc`, // comment 5
		`105:GenDecl:var.Pos`,
		`109:GenDecl:var.LeftParen`,
		`112:Comment:ValueSpec.Doc`, // comment 6
		`126:Ident:d`,
		`128:Ident:int`,
		`132:Comment:ValueSpec.Comment`, // comment 7
		`148:Comment:ValueSpec.Doc`,     // comment 8
		`162:Ident:e`,
		`164:Ident:int`,
		`168:Comment:ValueSpec.Comment`, // comment 9
		`181:GenDecl:var.RightParen`,
		`182:File:End`)
}

func Test_WalkPos_Struct(t *testing.T) {
	f := loadTest(t,
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
	checkWalkPos(t, f,
		`1:File:Start`,
		`1:File:Package`,
		`9:Ident:foo`,
		`13:Comment:GenDecl.Doc`, // comment 1
		// type Foo struct {
		`26:GenDecl:type.Pos`,
		`31:Ident:Foo`,
		`35:StructType:Struct`,
		`42:FieldList:Opening`,
		`44:Comment:File.Comment`, // comment 2
		`59:Comment:Field.Doc`,    // comment 3
		// x int `json:\"-\"`
		`73:Ident:x`,
		`75:Ident:int`,
		`79:BasicLit:ValuePos`,
		`90:Comment:Field.Comment`, // comment 4
		`105:Comment:File.Comment`, // comment 5
		`120:Comment:Field.Doc`,    // comment 6
		// y,
		`134:Ident:y`,
		`137:Comment:File.Comment`, // comment 7
		// z int
		`151:Ident:z`,
		`153:Ident:int`,
		`157:Comment:Field.Comment`, // comment 8
		`172:Comment:File.Comment`,  // comment 9
		`185:FieldList:Closing`,
		`186:File:End`)
}

func Test_WalkPos_Channels(t *testing.T) {
	f := loadTest(t,
		`package foo`,
		`func Foo(src <-chan int, dst chan<- int, notUsed chan int) {`,
		`	dst <- src`,
		`}`)
	checkWalkPos(t, f,
		`1:File:Start`,
		`1:File:Package`,
		`9:Ident:foo`,
		`13:FuncDecl:Func`,
		`18:Ident:Foo`,
		`21:FieldList:Opening`,
		// src <-chan int
		`22:Ident:src`,
		`26:ChanType:Begin`,
		`26:ChanType:Arrow`,
		`33:Ident:int`,
		// dst chan<- int
		`38:Ident:dst`,
		`42:ChanType:Begin`,
		`46:ChanType:Arrow`,
		`49:Ident:int`,
		// notUsed chan int
		`54:Ident:notUsed`,
		`62:ChanType:Begin`,
		`67:Ident:int`,
		`70:FieldList:Closing`,
		`72:BlockStmt:LeftBrace`,
		// dst <- src
		`75:Ident:dst`,
		`79:SendStmt:Arrow`,
		`82:Ident:src`,
		`86:BlockStmt:RightBrace`,
		`87:File:End`)
}

func checkWalkPos(t testing.TB, f *File, expLines ...string) {
	t.Helper()
	lines := []string{}
	var prior int
	for pt := range WalkPos(f.File) {
		if !pt.Pos.IsValid() {
			t.Errorf("invalid position returned for %s", pt.String())
		}
		if prior > int(*pt.Pos) {
			t.Errorf("out-of-order positions returned for %s", pt.String())
		}
		lines = append(lines, pt.String())
		prior = int(*pt.Pos)
	}
	if diff := cmp.Diff(expLines, lines); len(diff) > 0 {
		t.Errorf("the line for WalkPos didn't match expected lines:\n%s", diff)
	}
}
