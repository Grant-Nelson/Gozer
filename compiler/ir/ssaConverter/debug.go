package ssaConverter

import (
	"fmt"
	"go/types"
	"os"

	"golang.org/x/tools/go/ssa"
)

// Quiet not-used errors
var _ = printFunctionDebug
var _ = printMemberTree

func printFunctionDebug(f *ssa.Function) {
	if _, err := f.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
}

func printMemberTree(m ssa.Member) {
	fmt.Printf("Member: %s (%T)\n", m.Name(), m)
	switch v := m.(type) {
	case *ssa.Function:
		printFunctionTree(v, `  `, map[*ssa.Function]bool{})
	case *ssa.Global:
		fmt.Printf("  └─ Global: %s : %s\n", v.Name(), v.Type())
	case *ssa.Type:
		printTypeTree(v, `  `)
	case *ssa.NamedConst:
		fmt.Printf("  └─ NamedConst: %s = %s\n", v.Name(), v.Value.String())
	default:
		fmt.Printf("  └─ %v\n", v)
	}
}

func printTypeTree(v *ssa.Type, indent string) {
	ms := types.NewMethodSet(v.Type())
	ptrMS := types.NewMethodSet(types.NewPointer(v.Type()))
	fmt.Printf("%s└─ Type: %s (type=%s, value-methods=%d, pointer-methods=%d)\n",
		indent, v.Name(), v.Type(), ms.Len(), ptrMS.Len())
	for sel := range ms.Methods() {
		fmt.Printf("%s│   • val method: %s\n", indent, sel.Obj().Name())
	}
	for sel := range ptrMS.Methods() {
		fmt.Printf("%s│   • ptr method: %s\n", indent, sel.Obj().Name())
	}
}

func printFunctionTree(fn *ssa.Function, indent string, visited map[*ssa.Function]bool) {
	if fn == nil {
		fmt.Printf("%s└─ <nil *ssa.Function>\n", indent)
		return
	}
	if visited[fn] {
		fmt.Printf("%s└─ Function %s (already visited)\n", indent, fn.String())
		return
	}
	visited[fn] = true

	objInfo := `<none>`
	if tf, ok := fn.Object().(*types.Func); ok && tf != nil {
		objInfo = fmt.Sprintf(`*types.Func (%s)`, tf.FullName())
	} else if obj := fn.Object(); obj != nil {
		objInfo = fmt.Sprintf(`%T (%s)`, obj, obj.Name())
	}

	fmt.Printf("%s└─ Function %s | sig=%s | obj=%s\n", indent, fn.Name(), fn.Signature, objInfo)

	childIndent := indent + `    `
	for _, p := range fn.Params {
		fmt.Printf("%s• Param: %s : %s\n", childIndent, p.Name(), p.Type())
	}
	for _, fv := range fn.FreeVars {
		fmt.Printf("%s• FreeVar: %s : %s\n", childIndent, fv.Name(), fv.Type())
	}

	for i, b := range fn.Blocks {
		if b == nil {
			continue
		}
		connector := `├─`
		edge := `│   `
		if i == len(fn.Blocks)-1 && len(fn.AnonFuncs) == 0 {
			connector = `└─`
			edge = `    `
		}
		fmt.Printf("%s%s Block #%d (%s) [preds=%d, succs=%d]\n",
			childIndent, connector, b.Index, b.Comment, len(b.Preds), len(b.Succs))
		instrIndent := childIndent + edge

		if vars := blockUsedInputs(b); len(vars) > 0 {
			fmt.Printf(`%s│   • input: `, instrIndent)
			for j, v := range vars {
				if j > 0 {
					fmt.Print(`, `)
				}
				fmt.Printf(`[%T] %s : %s`, v, v.Name(), v.Type())
			}
			fmt.Println()
		}

		for j, instr := range b.Instrs {
			connector := `├─`
			if j == len(b.Instrs)-1 && len(fn.AnonFuncs) == 0 {
				connector = `└─`
			}
			valTag := ``
			if val, ok := instr.(ssa.Value); ok {
				valTag = fmt.Sprintf(` [Value: %s : %s]`, val.Name(), val.Type())
			}
			fmt.Printf("%s%s [%T] %s%s\n", instrIndent, connector, instr, instr.String(), valTag)
		}
	}

	for j, anon := range fn.AnonFuncs {
		connector := `├─`
		edge := `│   `
		if j == len(fn.AnonFuncs)-1 {
			connector = `└─`
			edge = `    `
		}
		fmt.Printf("%s%s AnonFunc #%d:\n", childIndent, connector, j)
		printFunctionTree(anon, childIndent+edge, visited)
	}
}

func blockUsedInputs(b *ssa.BasicBlock) (usedInBlock []ssa.Value) {
	seen := map[ssa.Value]bool{}
	for _, instr := range b.Instrs {
		if v, ok := instr.(ssa.Value); ok {
			seen[v] = true
		}
		for _, op := range instr.Operands(nil) {
			if op != nil && *op != nil && !seen[*op] {
				if _, constant := (*op).(*ssa.Const); constant {
					continue
				}
				seen[*op] = true
				usedInBlock = append(usedInBlock, *op)
			}
		}
	}
	return usedInBlock
}
