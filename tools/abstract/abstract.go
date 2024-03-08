package abstract

import (
	"fmt"

	"github.com/Snow-Gremlin/goToolbox/argers/args"

	"github.com/Snow-Gremlin/Gozer/readers"
	"github.com/Snow-Gremlin/Gozer/readers/golang"
	"github.com/Snow-Gremlin/Gozer/tools"
)

// New creates a new abstract tool.
func New() tools.Tool { return &toolImp{} }

type toolImp struct{}

func (t *toolImp) Name() string { return `abstract` }

func (t *toolImp) Aliases() []string { return []string{`ab`} }

func (t *toolImp) Summary() string {
	return `Abstract creates a JSON output containing design information.`
}

func (t *toolImp) Description() string {
	return `Abstract reads and analyses a Go project. It collects information ` +
		`about the structs, functions, etc that indicate how they work together. ` +
		`The information includes things like which fields are accessed by ` +
		`which functions and what receivers are used for functions.` +
		`With this information a design recovery algorithm can be run to ` +
		`reverse engineer the Go project into an estimate Object-Oriented ` +
		`membership map called a participation matrix.`
}

func (t *toolImp) Run(ctx tools.Context) (int, error) {
	input := `.`
	output := `.`

	args.New().
		NamedStr(&input, `i`, `input`).
		NamedStr(&output, `o`, `output`).
		Process(ctx.Args()[2:])

	proj, err := golang.New().Read(&readers.Config{
		MainPackageDir: input,
	})
	if err != nil {
		return 1, err
	}

	fmt.Printf(">> result: %+v\n", proj) // TODO: REMOVE

	return 0, nil
}
