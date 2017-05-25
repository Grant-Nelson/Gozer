package statements

import (
	"strings"
	"testing"

	"github.com/grant-nelson/Gozer/constructs/expressions"
	"github.com/grant-nelson/Gozer/constructs/types"
)

func TestBlock(t *testing.T) {
	checkString(t, ((*BlockStat)(nil)).String(), "nil")
	blk := Block()
	checkStat(t, blk, "{}")

	id := expressions.Identifier("name", types.String())
	fType := types.Function().AddParam("who", types.String())
	blk.Statements = []Statement{
		expressions.Definition(id,
			expressions.Literal(`"World"`, types.String())),
		expressions.Call(fType,
			expressions.Identifier("SayHello", fType),
			[]expressions.Expression{id})}
	checkStat(t, blk,
		`{`,
		`  string name = "World"`,
		`  SayHello(name)`,
		`}`)
}

func TestBranch(t *testing.T) {
	checkString(t, ((*BranchStat)(nil)).String(), "nil")
	checkStat(t, Branch(true), "break")
	checkStat(t, Branch(false), "continue")
}

func TestDecIncOp(t *testing.T) {
	checkString(t, ((*DecIncOpStat)(nil)).String(), "nil")
	id := expressions.Identifier("i", types.Int16())
	checkStat(t, IncDecOp(id, true), "i++")
	checkStat(t, IncDecOp(id, false), "i--")
}

func TestFor(t *testing.T) {
	checkString(t, ((*ForStat)(nil)).String(), "nil")

	stat := For(nil, nil, nil, nil)
	checkStat(t, stat, "for(; ; ) nil")

	id := expressions.Identifier("a", types.Int())
	stat.Init = expressions.Definition(id, expressions.Literal("0", types.Int()))
	stat.Cond = expressions.BinaryOp(id, expressions.Literal("10", types.Int()), expressions.LessThanOp, types.Bool())
	stat.Post = IncDecOp(id, true)
	checkStat(t, stat, "for(int a = 0; (a < 10); a++) nil")

	fType := types.Function().AddParam("a", types.String())
	stat.Body = expressions.Call(fType,
		expressions.Identifier("print", fType),
		[]expressions.Expression{id})
	checkStat(t, stat, "for(int a = 0; (a < 10); a++) print(a)")
}

func TestIf(t *testing.T) {
	checkString(t, ((*IfStat)(nil)).String(), "nil")

	stat := If(nil, nil, nil)
	checkStat(t, stat, "if nil nil")

	id := expressions.Identifier("name", types.String())
	stat.Cond = expressions.BinaryOp(id, expressions.Literal(`"World"`, types.String()), expressions.EqualOp, types.Bool())
	checkStat(t, stat, `if (name == "World") nil`)

	fType := types.Function().AddParam("a", types.String()).SetEllipse(true)
	stat.Body = expressions.Call(fType,
		expressions.Identifier("print", fType),
		[]expressions.Expression{
			expressions.Literal(`"Hello "`, types.String()),
			id})
	checkStat(t, stat, `if (name == "World") print("Hello ", name)`)

	stat.Else = expressions.Call(fType,
		expressions.Identifier("print", fType),
		[]expressions.Expression{
			expressions.Literal(`"Goodnight "`, types.String()),
			id})
	checkStat(t, stat, `if (name == "World") print("Hello ", name) else print("Goodnight ", name)`)
}

//============================================================================

// checkString checks the the given string matches the given expected lines.
// The lines will be joined with newlines.
func checkString(t *testing.T, result string, exp ...string) {
	expStr := strings.Join(exp, "\n")
	if result != expStr {
		t.Fatal("Unexpected construct string:",
			"\n   Expected: ", expStr,
			"\n   Gotten:   ", result)
	}
}

// checkStat checks that the statement's string matches the given string.
func checkStat(t *testing.T, s Statement, exp ...string) {
	checkString(t, ToString(s), exp...)
}
