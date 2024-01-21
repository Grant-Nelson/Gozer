package main

import (
	"os"

	"github.com/Snow-Gremlin/Gozer/tools/tool"
)

func main() {
	if err := tool.New().Run(os.Args[1:]); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
