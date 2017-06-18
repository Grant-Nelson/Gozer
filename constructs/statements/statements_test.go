package statements

import (
	"testing"

	"github.com/grant-nelson/Gozer/common"
	"github.com/grant-nelson/Gozer/constructs/expressions"
	"github.com/grant-nelson/Gozer/constructs/types"
)

func TestBlock(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*BlockStat)(nil).String(), "nil")
	blk := Block()
	CheckStat(t, blk, "{}")

	id := expressions.Identifier("name", types.String())
	fType := types.Function().AddParam("who", types.String())
	blk.Statements = []Statement{
		expressions.Definition(id,
			expressions.Literal(`"World"`, types.String())),
		expressions.Call(fType,
			expressions.Identifier("SayHello", fType),
			[]expressions.Expression{id})}
	CheckStat(t, blk,
		`{`,
		`  string name = "World"`,
		`  SayHello(name)`,
		`}`)
}

func TestBranch(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*BranchStat)(nil).String(), "nil")
	CheckStat(t, Branch(true), "break")
	CheckStat(t, Branch(false), "continue")
}

func TestDecIncOp(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*DecIncOpStat)(nil).String(), "nil")
	id := expressions.Identifier("i", types.Int16())
	CheckStat(t, IncDecOp(id, true), "i++")
	CheckStat(t, IncDecOp(id, false), "i--")
}

func TestFor(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*ForStat)(nil).String(), "nil")

	stat := For(nil, nil, nil, nil)
	CheckStat(t, stat, "for(; ; ) nil")

	id := expressions.Identifier("a", types.Int())
	stat.Init = expressions.Definition(id, expressions.Literal("0", types.Int()))
	stat.Cond = expressions.BinaryOp(id, expressions.Literal("10", types.Int()), expressions.LessThanOp, types.Bool())
	stat.Post = IncDecOp(id, true)
	CheckStat(t, stat, "for(int a = 0; (a < 10); a++) nil")

	fType := types.Function().AddParam("a", types.String())
	stat.Body = expressions.Call(fType,
		expressions.Identifier("print", fType),
		[]expressions.Expression{id})
	CheckStat(t, stat, "for(int a = 0; (a < 10); a++) print(a)")
}

func TestIf(tt *testing.T) {
	t := common.NewTester(tt)
	t.CheckStr((*IfStat)(nil).String(), "nil")

	stat := If(nil, nil, nil)
	CheckStat(t, stat, "if nil nil")

	id := expressions.Identifier("name", types.String())
	stat.Cond = expressions.BinaryOp(id, expressions.Literal(`"World"`, types.String()), expressions.EqualOp, types.Bool())
	CheckStat(t, stat, `if (name == "World") nil`)

	fType := types.Function().AddParam("a", types.String()).SetEllipse(true)
	stat.Body = expressions.Call(fType,
		expressions.Identifier("print", fType),
		[]expressions.Expression{
			expressions.Literal(`"Hello "`, types.String()),
			id})
	CheckStat(t, stat, `if (name == "World") print("Hello ", name)`)

	stat.Else = expressions.Call(fType,
		expressions.Identifier("print", fType),
		[]expressions.Expression{
			expressions.Literal(`"Goodnight "`, types.String()),
			id})
	CheckStat(t, stat, `if (name == "World") print("Hello ", name) else print("Goodnight ", name)`)
}

//============================================================================

// CheckStat checks that the statement's string matches the given string.
func CheckStat(t *common.Tester, s Statement, exp ...string) {
	t.CheckStr(ToString(s), exp...)
}
