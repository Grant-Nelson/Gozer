package filePos

import (
	"go/token"
	"strings"
	"testing"

	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
	"github.com/google/go-cmp/cmp"
)

func Test_FileSet_Widths_Simple(t *testing.T) {
	f := loadTest(t,
		`package test`, // 1 package, 9 test
		``,
		`import "fmt"`, // 15 import, 22 "fmt"
		``,
		`func main() {`,               // 29 func, 34 main, 38 (, 39 ), 41 {
		`	fmt.Println("Hello World")`, // 44 fmt., 48 Println, 55 (, 56 "Hello World", 69 )
		`}`,                           // 71 }, 72 [eof]
	)
	fp := New(f.TempFileSet())
	fp.RegisterFile(f)

	checkFileSetWidths(t, fp, 1, 8, []int{8})        // 1 package
	checkFileSetWidths(t, fp, 9, 6, []int{5, 1, 0})  // 9 test
	checkFileSetWidths(t, fp, 15, 7, []int{7})       // 15 import
	checkFileSetWidths(t, fp, 22, 7, []int{6, 1, 0}) // 22 "fmt"
	checkFileSetWidths(t, fp, 29, 5, []int{5})       // 29 func
	checkFileSetWidths(t, fp, 34, 4, []int{4})       // 34 main
	checkFileSetWidths(t, fp, 38, 1, []int{1})       // 38 (
	checkFileSetWidths(t, fp, 39, 2, []int{2})       // 39 )
	checkFileSetWidths(t, fp, 41, 3, []int{2, 1})    // 41 {
	checkFileSetWidths(t, fp, 44, 4, []int{4})       // 44 fmt.
	checkFileSetWidths(t, fp, 48, 7, []int{7})       // 48 Println
	checkFileSetWidths(t, fp, 55, 1, []int{1})       // 55 (
	checkFileSetWidths(t, fp, 56, 13, []int{13})     // 56 "Hello World"
	checkFileSetWidths(t, fp, 69, 2, []int{2, 0})    // 69 )
	checkFileSetWidths(t, fp, 71, 1, []int{1})       // 71 }
	checkFileSetWidths(t, fp, 72, 0, []int{0})       // 72 [eof]

	checkFileSetNeighbors(t, fp, 1, 1, 9)
	checkFileSetNeighbors(t, fp, 9, 1, 15)
	checkFileSetNeighbors(t, fp, 15, 9, 22)
	checkFileSetNeighbors(t, fp, 22, 15, 29)
	checkFileSetNeighbors(t, fp, 69, 56, 71)
	checkFileSetNeighbors(t, fp, 71, 69, 72)
	checkFileSetNeighbors(t, fp, 72, 71, 72)
}

// TODO: Add more tests to check more of walkPos

func loadTest(t testing.TB, code ...string) *artifacts.File {
	t.Helper()
	fs := token.NewFileSet()
	f, err := artifacts.Load(fs, `test.go`, strings.Join(code, "\n"))
	if err != nil {
		t.Fatalf(`failed to load test file: %v`, err)
	}
	return f
}

func checkFileSetWidths(t testing.TB, fp *FilePos, pos int, expTotal int, expLines []int) {
	t.Helper()
	total, lines := fp.Widths(token.Pos(pos))
	if total != expTotal {
		t.Errorf(`pos %d: total width (%d) doesn't match expected total width (%d)`, pos, total, expTotal)
	}
	if diff := cmp.Diff(expLines, lines); len(diff) > 0 {
		t.Errorf("pos %d: the line widths didn't match expected lines:\n%s", pos, diff)
	}
}

func checkFileSetNeighbors(t testing.TB, fp *FilePos, pos int, expPrev, expNext int) {
	t.Helper()
	if prev := fp.FindPrevious(token.Pos(pos)); prev != token.Pos(expPrev) {
		t.Errorf("pos %d: the previous found node was not expected:\n\texpected:%d\n\tgotten:%d", pos, expPrev, prev)
	}
	if next := fp.FindNext(token.Pos(pos)); next != token.Pos(expNext) {
		t.Errorf("pos %d: the next found node was not expected:\n\texpected:%d\n\tgotten:%d", pos, expNext, next)
	}
}
