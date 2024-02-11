package help

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Snow-Gremlin/Gozer/tools"
)

// New creates a new help tool.
func New() tools.Tool { return &toolImp{} }

type toolImp struct{}

func (t *toolImp) Name() string { return `help` }

func (t *toolImp) Aliases() []string { return []string{`h`} }

func (t *toolImp) Summary() string {
	return `Shows tool summaries. Add a tool name after ` +
		`"help" to get details for that tool.`
}

func (t *toolImp) Description() string {
	return `Gozer the Traveller. He will come in one of the pre-chosen forms. ` +
		`During the rectification of the Vuldronaii, the Traveller came as a ` +
		`large, moving Torb! Then, during the third reconciliation of the ` +
		`last of the Meketrex Supplicants, they chose a new form for him...` +
		`That of a giant Sloar! Many Shubs and Zulls knew what it was to be ` +
		`roasted in the depths of the Sloar that day, I can tell you.`
}

func (t *toolImp) Run(ctx tools.Context) (int, error) {
	argN := len(ctx.Args())
	switch {
	case argN < 3:
		return t.showSummaries(ctx)

	case argN == 3:
		return t.showToolDetails(ctx)

	default:
		fmt.Println(`Unexpected arguments for help. Please provide only one tool to get help for.`)
		return 1, nil
	}
}

func (t *toolImp) showSummaries(ctx tools.Context) (int, error) {
	count := ctx.Tools().Count()
	parts := make([]string, count)
	for i := count - 1; i >= 0; i-- {
		t := ctx.Tools().At(i)

		names := append([]string{t.Name()}, t.Aliases()...)
		parts[i] = fmt.Sprintf("\t%s:\n\t%s\n", strings.Join(names, `, `), t.Summary())
	}
	sort.Strings(parts)

	_, _ = fmt.Print("Gozer has the following tools available:\n", strings.Join(parts, ``))
	return 0, nil
}

func (t *toolImp) showToolDetails(ctx tools.Context) (int, error) {
	toolName := ctx.Args()[2]
	tool := ctx.Tools().Get(toolName)

	if tool != nil {
		_, _ = fmt.Printf("%s:\n%s\n", toolName, tool.Description())
		return 0, nil
	}

	fmt.Printf("No tool by the name %q exists.\nPlease provide the "+
		"tool you wish to get help for:\n\t%s\n",
		toolName, strings.Join(ctx.Tools().Names(), "\n\t"))
	return 1, nil
}
