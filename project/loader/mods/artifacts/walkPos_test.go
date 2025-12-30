package artifacts

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_WalkPos_Func(t *testing.T) {
	f := loadTest(t,
		`// comment 1`,
		``,
		`// comment 2`,
		`package foo`,
		`// comment 3`,
		``,
		`// comment 4`)
	checkWalkPos(t, f, `bbbb`)
}

func checkWalkPos(t testing.TB, f *File, expLines ...string) {
	lines := []string{}
	for pt := range WalkPos(f.File) {
		lines = append(lines, fmt.Sprintf("%4d:%T:%s\n", int(*pt.Pos), pt.Node, pt.Name))
	}
	if diff := cmp.Diff(expLines, lines); len(diff) > 0 {
		t.Errorf("the line for WalkPos didn't match expected lines:\n%s", diff)
	}
}
