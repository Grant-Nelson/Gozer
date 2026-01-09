package expectations

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Test_DocAndComment checks that comments in the file are placed into expected
// AST Doc and Comment nodes. The assumption of where comments are put is used
// during augmentation to ensure that no important comments, like directives,
// are lost or attached to the wrong nodes.
func Test_DocAndComment(t *testing.T) {
	_, f := parseFile(t, lines(
		`// file docs`,
		`/* multi-line`,
		`   file docs */`,
		`package foo`,
		``,
		`// file comment`,
		``,
		`// func doc`,
		`func main(`,
		`	// param doc`,
		`	x int, // param comment`,
		`) {`,
		`	// func inner comment`,
		`}`,
		``,
		`// type decl doc`,
		`type (`,
		`	// type spec doc`,
		`	foo struct {`,
		`		// field doc`,
		`		x int // field comment`,
		``,
		`		// type inner comment 1 group 1`,
		`		// type inner comment 2 group 1`,
		``,
		`		// type inner comment 3 group 2`,
		`	}`,
		`)`,
		``,
		`// file end comment`,
	))
	equalLines(t, getCommentInfo(f), lines(
		`File:`,
		`   Doc: [0](1) "// file docs"`,
		`      [1](14) "/* multi-line\n   file docs */"`,
		`   Comments:`,
		`      0: [0](1) "// file docs"`,
		`         [1](14) "/* multi-line\n   file docs */"`,
		`      1: [2](57) "// file comment"`,
		`      2: [3](74) "// func doc"`,
		`      3: [4](98) "// param doc"`,           // floating (only here)
		`      4: [5](119) "// param comment"`,      // floating
		`      5: [6](141) "// func inner comment"`, // floating
		`      6: [7](166) "// type decl doc"`,
		`      7: [8](191) "// type spec doc"`,
		`      8: [9](224) "// field doc"`,
		`      9: [10](245) "// field comment"`,
		`      10: [11](265) "// type inner comment 1 group 1"`, // floating
		`         [12](299) "// type inner comment 2 group 1"`,  // floating
		`      11: [13](334) "// type inner comment 3 group 2"`, // floating
		`      12: [14](372) "// file end comment"`,             // floating
		`FuncDecl:`,
		`   Doc: [3](74) "// func doc"`,
		`Field:`, // main.x int param (comments don't attach)
		`   Doc: <nil>`,
		`   Comment: <nil>`,
		`GenDecl:`,
		`   Doc: [7](166) "// type decl doc"`,
		`TypeSpec:`,
		`   Doc: [8](191) "// type spec doc"`,
		`   Comment: <nil>`,
		`Field:`, // foo.x int field (comments will attach)
		`   Doc: [9](224) "// field doc"`,
		`   Comment: [10](245) "// field comment"`,
	))
}

func lines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func parseFile(t testing.TB, src string) (*token.FileSet, *ast.File) {
	t.Helper()
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, `test.go`, src, parser.ParseComments)
	if err != nil {
		t.Fatalf(`error parsing test file: %v`, err)
	}
	return fSet, f
}

const defaultIndent = `   `

func addCommentGroupInfo(buf *strings.Builder, ptr map[string]int, indent, name string, cg *ast.CommentGroup) {
	if cg == nil {
		fmt.Fprintf(buf, "%s%s: <nil>\n", indent, name)
		return
	}
	if len(cg.List) <= 0 {
		fmt.Fprintf(buf, "%s%s: <empty>\n", indent, name)
		return
	}

	fmt.Fprintf(buf, "%s%s: ", indent, name)
	subIndent := ``
	for _, c := range cg.List {
		p := fmt.Sprintf("%p", c)
		num, exists := ptr[p]
		if !exists {
			num = len(ptr)
			ptr[p] = num
		}
		fmt.Fprintf(buf, "%s[%d](%d) %q\n", subIndent, num, int(c.Slash), c.Text)
		subIndent = indent + defaultIndent
	}
}

func getCommentInfo(f *ast.File) string {
	buf := &strings.Builder{}
	ptr := map[string]int{}
	ast.Inspect(f, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.File:
			fmt.Fprintln(buf, "File:")
			addCommentGroupInfo(buf, ptr, defaultIndent, `Doc`, n.Doc)
			fmt.Fprintln(buf, defaultIndent+"Comments:")
			for i, cg := range n.Comments {
				addCommentGroupInfo(buf, ptr, defaultIndent+defaultIndent, strconv.Itoa(i), cg)
			}

		case *ast.Field:
			fmt.Fprintln(buf, `Field:`)
			addCommentGroupInfo(buf, ptr, defaultIndent, `Doc`, n.Doc)
			addCommentGroupInfo(buf, ptr, defaultIndent, `Comment`, n.Comment)

		case *ast.GenDecl:
			fmt.Fprintln(buf, `GenDecl:`)
			addCommentGroupInfo(buf, ptr, defaultIndent, `Doc`, n.Doc)

		case *ast.ImportSpec:
			fmt.Fprintln(buf, `ImportSpec:`)
			addCommentGroupInfo(buf, ptr, defaultIndent, `Doc`, n.Doc)
			addCommentGroupInfo(buf, ptr, defaultIndent, `Comment`, n.Comment)

		case *ast.ValueSpec:
			fmt.Fprintln(buf, `ValueSpec:`)
			addCommentGroupInfo(buf, ptr, defaultIndent, `Doc`, n.Doc)
			addCommentGroupInfo(buf, ptr, defaultIndent, `Comment`, n.Comment)

		case *ast.TypeSpec:
			fmt.Fprintln(buf, `TypeSpec:`)
			addCommentGroupInfo(buf, ptr, defaultIndent, `Doc`, n.Doc)
			addCommentGroupInfo(buf, ptr, defaultIndent, `Comment`, n.Comment)

		case *ast.FuncDecl:
			fmt.Fprintln(buf, `FuncDecl:`)
			addCommentGroupInfo(buf, ptr, defaultIndent, `Doc`, n.Doc)
		}
		return true
	})
	return strings.TrimSpace(buf.String())
}

func equalLines(t testing.TB, got string, exp string) {
	t.Helper()
	if diff := cmp.Diff(strings.Split(got, "\n"), strings.Split(exp, "\n")); len(diff) > 0 {
		t.Logf("Got:\n%s\n", got)
		t.Errorf("resulting lines didn't match expected lines:\n%s", diff)
	}
}
