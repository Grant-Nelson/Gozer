package golang

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/testers/check"

	"github.com/Snow-Gremlin/Gozer/constructs/cPackage"
	"github.com/Snow-Gremlin/Gozer/constructs/cProject"
)

func checkConvertSource(t *testing.T, source ...string) func(exp ...string) {
	p := cPackage.New()
	p.SetName(`testPackage`)
	p.SetPath(`testPath`)

	fSet := token.NewFileSet()

	proj := cProject.New()
	proj.SetName(`testProject`)
	proj.Packages().Add(p)

	con := newConverter(p, nil, fSet, proj)

	code := strings.Join(source, "\n")
	f, err := parser.ParseFile(con.fSet, `testFile`, code, parser.ParseComments)
	check.NoError(t).Require(err)
	con.addFileNode(f)

	return func(exp ...string) {
		e := strings.Join(exp, "\n")
		a := fmt.Sprint(p)
		check.Equal(t, e).Assert(a)
	}
}

func Test_Convert_Struct(t *testing.T) {
	checkConvertSource(t,
		`package test`,
		`type Point struct{ x, y int }`,
	)(
		``,
	)
}
