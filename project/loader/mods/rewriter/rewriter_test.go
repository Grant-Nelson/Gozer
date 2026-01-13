package rewriter

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"testing"

	"github.com/Grant-Nelson/Gozer/avail/iterator"
	"github.com/Grant-Nelson/Gozer/project/loader/parser"
	"github.com/google/go-cmp/cmp"
)

func Test_Rewriter_FileData_Simple(t *testing.T) {
	_, _, _, fd := loadTestFile(t,
		`package test`, // 1 package, 9 test
		``,
		`import "fmt"`, // 15 import, 22 "fmt"
		``,
		`func main() {`,               // 29 func, 34 main, 38 (, 39 ), 41 {
		`	fmt.Println("Hello World")`, // 44 fmt, 48 Println, 55 (, 56 "Hello World", 69 )
		`}`,                           // 71 }, 72 [eof]
	)

	checkPosOrder(t, fd, 1, 9, 15, 22, 29, 34, 38, 39, 41, 44, 48, 55, 56, 69, 71, 72)

	checkFileData(t, fd, 1, ` t: 8;  w: 7;  ln:[ 8      ];  tl:[1      ];  p: 1;  n: 9`) //  1 package
	checkFileData(t, fd, 9, ` t: 6;  w: 4;  ln:[ 5, 1, 0];  tl:[1, 1, 0];  p: 1;  n:15`) //  9 test
	checkFileData(t, fd, 15, `t: 7;  w: 6;  ln:[ 7      ];  tl:[1      ];  p: 9;  n:22`) // 15 import
	checkFileData(t, fd, 22, `t: 7;  w: 5;  ln:[ 6, 1, 0];  tl:[1, 1, 0];  p:15;  n:29`) // 22 "fmt"
	checkFileData(t, fd, 29, `t: 5;  w: 4;  ln:[ 5      ];  tl:[1      ];  p:22;  n:34`) // 29 func
	checkFileData(t, fd, 34, `t: 4;  w: 4;  ln:[ 4      ];  tl:[0      ];  p:29;  n:38`) // 34 main
	checkFileData(t, fd, 38, `t: 1;  w: 1;  ln:[ 1      ];  tl:[0      ];  p:34;  n:39`) // 38 (
	checkFileData(t, fd, 39, `t: 2;  w: 1;  ln:[ 2      ];  tl:[1      ];  p:38;  n:41`) // 39 )
	checkFileData(t, fd, 41, `t: 3;  w: 1;  ln:[ 2, 1   ];  tl:[1, 1   ];  p:39;  n:44`) // 41 {
	checkFileData(t, fd, 44, `t: 4;  w: 3;  ln:[ 4      ];  tl:[1      ];  p:41;  n:48`) // 44 fmt
	checkFileData(t, fd, 48, `t: 7;  w: 7;  ln:[ 7      ];  tl:[0      ];  p:44;  n:55`) // 48 Println
	checkFileData(t, fd, 55, `t: 1;  w: 1;  ln:[ 1      ];  tl:[0      ];  p:48;  n:56`) // 55 (
	checkFileData(t, fd, 56, `t:13;  w:13;  ln:[13      ];  tl:[0      ];  p:55;  n:69`) // 56 "Hello World"
	checkFileData(t, fd, 69, `t: 2;  w: 1;  ln:[ 2, 0   ];  tl:[1, 0   ];  p:56;  n:71`) // 69 )
	checkFileData(t, fd, 71, `t: 1;  w: 1;  ln:[ 1      ];  tl:[0      ];  p:69;  n:72`) // 71 }
	checkFileData(t, fd, 72, `t: 0;  w: 0;  ln:[ 0      ];  tl:[0      ];  p:71;  n:72`) // 72 [eof]
}

func Test_Rewriter_FileData_MultilineConst(t *testing.T) {
	_, _, _, fd := loadTestFile(t,
		`package test`, // 1 package, 9 test
		``,
		"const Text = `Hello", // 15 const, 21 Text, 28 `Hello\nPale Blue `
		"Pale Blue `+",        // 46 +
		"`Dot`",               // 48 `Dot`
		``,
		`/*`, // 55 /*\n  A mote of dust\n  suspended in a sunbeam.\n*/
		`  A mote of dust`,
		`  suspended in a sunbeam.`,
		`*/`, // 103 [eof]
	)

	checkPosOrder(t, fd, 1, 9, 15, 21, 28, 46, 48, 55, 103)

	checkFileData(t, fd, 1, `  t: 8;  w: 7;  ln:[8           ];  tl:[1      ];  p: 1;  n:  9`) //   1 package
	checkFileData(t, fd, 9, `  t: 6;  w: 4;  ln:[5,  1,  0   ];  tl:[1, 1, 0];  p: 1;  n: 15`) //   9 test
	checkFileData(t, fd, 15, ` t: 6;  w: 5;  ln:[6           ];  tl:[1      ];  p: 9;  n: 21`) //  15 const
	checkFileData(t, fd, 21, ` t: 7;  w: 4;  ln:[7           ];  tl:[3      ];  p:15;  n: 28`) //  21 Text
	checkFileData(t, fd, 28, ` t:18;  w:18;  ln:[7, 11       ];  tl:[0      ];  p:21;  n: 46`) //  28 `Hello\nPale Blue `
	checkFileData(t, fd, 46, ` t: 2;  w: 1;  ln:[2,  0       ];  tl:[1, 0   ];  p:28;  n: 48`) //  46 +
	checkFileData(t, fd, 48, ` t: 7;  w: 5;  ln:[6,  1,  0   ];  tl:[1, 1, 0];  p:46;  n: 55`) //  48 `Dot`
	checkFileData(t, fd, 55, ` t:48;  w:48;  ln:[3, 17, 26, 2];  tl:[0      ];  p:48;  n:103`) //  55 /*\n  A mote of dust\n  suspended in a sunbeam.\n*/
	checkFileData(t, fd, 103, `t: 0;  w: 0;  ln:[0           ];  tl:[0      ];  p:55;  n:103`) // 103 [eof]
}

func loadTestFile(t testing.TB, code ...string) (*ast.File, *token.FileSet, *Rewriter, *fileData) {
	t.Helper()
	rw := New(parser.Default)
	fs := token.NewFileSet()
	f, err := rw.Parser(fs, `test.go`, strings.Join(code, "\n"))
	if err != nil {
		t.Fatalf(`failed to load test file: %v`, err)
	}
	ds := rw.getFileData(fs, f)
	return f, fs, rw, ds
}

func joinIntSlice(values []int) string {
	return iterator.Iterate(values...).Join(`,`)
}

func checkPosOrder(t testing.TB, fd *fileData, exp ...int) {
	gotten := fd.PosOrder()
	if diff := cmp.Diff(exp, gotten); len(diff) > 0 {
		t.Errorf("the pos orders didn't match expected values:\n%s", diff)
	}
}

func checkFileData(t testing.TB, fd *fileData, pos int, exp string) {
	t.Helper()
	tp := token.Pos(pos)
	total := fd.Total(tp)
	width := fd.Width(tp)
	lines := fd.Lines(tp)
	tails := fd.Tail(tp)
	prev := fd.Previous(tp)
	next := fd.Next(tp)

	gotten := []string{
		fmt.Sprintf(`t:%d`, total),
		fmt.Sprintf(`w:%d`, width),
		fmt.Sprintf(`ln:[%s]`, joinIntSlice(lines)),
		fmt.Sprintf(`tl:[%s]`, joinIntSlice(tails)),
		fmt.Sprintf(`p:%d`, prev),
		fmt.Sprintf(`n:%d`, next)}

	expLines := strings.Split(exp, `;`)
	for i, ln := range expLines {
		expLines[i] = strings.ReplaceAll(ln, ` `, ``)
	}

	if diff := cmp.Diff(expLines, gotten); len(diff) > 0 {
		t.Errorf("pos %d: the line widths didn't match expected lines:\n%s", pos, diff)
	}
}
