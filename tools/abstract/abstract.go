package abstract

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Snow-Gremlin/goToolbox/argers/args"

	"github.com/Snow-Gremlin/Gozer/reader"
	"github.com/Snow-Gremlin/Gozer/tools"
	"github.com/Snow-Gremlin/Gozer/tools/abstract/models"
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
	verbose := false
	minimize := false
	input := `.`
	output := ``

	args.New().
		Flag(&verbose, `v`, `verbose`).
		Flag(&minimize, `m`, `min`).
		NamedStr(&input, `i`, `input`).
		NamedStr(&output, `o`, `output`).
		Process(ctx.Args()[2:])

	proj, err := reader.Read(&reader.Config{
		Verbose: verbose,
		Path:    input,
	})
	if err != nil {
		return 1, err
	}

	model := abstract(proj)

	data, err := jsonMarshal(minimize, model)
	if err != nil {
		return 1, err
	}

	err = writeJson(output, data)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func jsonMarshal(minimize bool, model models.ProjectModel) ([]byte, error) {
	if minimize {
		return json.Marshal(model)
	}
	return json.MarshalIndent(model, ``, `  `)
}

func writeJson(path string, data []byte) error {
	if len(path) > 0 {
		return os.WriteFile(path, data, 0666)
	}

	_, err := fmt.Println(string(data))
	return err
}
