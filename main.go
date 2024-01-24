package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Snow-Gremlin/Gozer/tools"
	"github.com/Snow-Gremlin/Gozer/tools/abstract"
	"github.com/Snow-Gremlin/Gozer/tools/convert"
	"github.com/Snow-Gremlin/Gozer/tools/help"
)

func addAllTools(ctx *tools.Context) {
	ctx.AddTool(help.New())
	ctx.AddTool(convert.New())
	ctx.AddTool(abstract.New())
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic: %v\n", r)
			os.Exit(1)
		}
	}()

	ctx := tools.NewContext()
	addAllTools(ctx)

	if len(ctx.Args()) <= 1 {
		fmt.Printf("Please provide the tool you wish to use:\n\t%s\n",
			strings.Join(ctx.ToolNames(), "\n\t"))
		os.Exit(1)
	}

	toolName := ctx.Args()[1]
	tool := ctx.GetTool(toolName)
	if tool == nil {
		fmt.Printf("Invalid tool name %s.\nPlease provide the tool you wish to use:\n\t%s\n",
			toolName, strings.Join(ctx.ToolNames(), "\n\t"))
		os.Exit(1)
	}

	exitValue, err := tool.Run(ctx)
	if err != nil {
		fmt.Println(`Error:` + err.Error())
		os.Exit(1)
	}
	os.Exit(exitValue)
}
