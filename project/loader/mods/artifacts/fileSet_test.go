package artifacts

import (
	"go/token"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_FileSet_Widths_Simple(t *testing.T) {
	fs := loadTest(t,
		`package test`, // 1 package, 9 test
		``,
		`import "fmt"`, // 15 import, 22 "fmt"
		``,
		`func main() {`,               // 29 func, 34 main, 38 (, 39 ), 41 {
		`	fmt.Println("Hello World")`, // 44 fmt., 48 Println, 55 (, 56 "Hello World", 69 )
		`}`,                           // 71 }, 72 [eof]
	).FileSet

	checkFileSetWidths(t, fs, 1, 8, []int{8})        // 1 package
	checkFileSetWidths(t, fs, 9, 6, []int{5, 1, 0})  // 9 test
	checkFileSetWidths(t, fs, 15, 7, []int{7})       // 15 import
	checkFileSetWidths(t, fs, 22, 7, []int{6, 1, 0}) // 22 "fmt"
	checkFileSetWidths(t, fs, 29, 5, []int{5})       // 29 func
	checkFileSetWidths(t, fs, 34, 4, []int{4})       // 34 main
	checkFileSetWidths(t, fs, 38, 1, []int{1})       // 38 (
	checkFileSetWidths(t, fs, 39, 2, []int{2})       // 39 )
	checkFileSetWidths(t, fs, 41, 3, []int{2, 1})    // 41 {
	checkFileSetWidths(t, fs, 44, 4, []int{4})       // 44 fmt.
	checkFileSetWidths(t, fs, 48, 7, []int{7})       // 48 Println
	checkFileSetWidths(t, fs, 55, 1, []int{1})       // 55 (
	checkFileSetWidths(t, fs, 56, 13, []int{13})     // 56 "Hello World"
	checkFileSetWidths(t, fs, 69, 2, []int{2, 0})    // 69 )
	checkFileSetWidths(t, fs, 71, 1, []int{1})       // 71 }
	checkFileSetWidths(t, fs, 72, 0, []int{0})       // 72 [eof]
}

// TODO: Add more tests to check more of walkPos

func loadTest(t testing.TB, code ...string) *File {
	t.Helper()
	fs := NewFileSet()
	f, err := Load(fs, `fileSetWidths.go`, strings.Join(code, "\n"))
	if err != nil {
		t.Fatalf(`failed to load test file: %v`, err)
	}
	return f
}

func checkFileSetWidths(t testing.TB, fs *FileSet, pos int, expTotal int, expLines []int) {
	total, lines := fs.Widths(token.Pos(pos))
	if total != expTotal {
		t.Errorf(`pos %d: total width (%d) doesn't match expected total width (%d)`, pos, total, expTotal)
	}
	if diff := cmp.Diff(expLines, lines); len(diff) > 0 {
		t.Errorf("pos %d: the line widths didn't match expected lines:\n%s", pos, diff)
	}
}
