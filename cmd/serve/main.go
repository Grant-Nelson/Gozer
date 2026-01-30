package serve

import (
	"log"
	"net/http"

	"github.com/Grant-Nelson/Gozer/avail/args"
	"github.com/Grant-Nelson/Gozer/cmd/serve/internal/builder"
	"github.com/Grant-Nelson/Gozer/cmd/serve/internal/cache"
	"github.com/Grant-Nelson/Gozer/cmd/serve/internal/fileReader"
	"github.com/Grant-Nelson/Gozer/cmd/serve/internal/handler"
	"github.com/Grant-Nelson/Gozer/cmd/serve/internal/redirection"
)

type Config struct {
	Help     string `arg:"help"`
	BasePath string `arg:"optional,pos,BasePath,Base path to serve files from"`
	Verbose  bool   `arg:"optional,flag,v|verbose,Print verbose output"`
	Address  string `arg:"optional,flag,a|addr,Address to serve to"`
	Mini     bool   `arg:"optional,flag,m|mini,Minifies any *.js being built"`
}

func main() {
	cfg := &Config{
		Help: `Serve will serve file(s) at the base path. ` +
			`If a requested file is *.js that doesn't exist on disk, ` +
			`then the server will look for a *.ts with the same name. ` +
			`If a *.ts exists, then it will be compiled for the missing *.js ` +
			`and source mapping will also be generated. If index.html doesn't ` +
			`exist on disk, a very simple one will be generated to load all of ` +
			`the *.js and *.ts scripts.`,
		BasePath: `.`,
		Verbose:  false,
		Address:  `:8080`,
	}
	args.Parse(cfg)

	log.Fatal(http.ListenAndServe(cfg.Address,
		handler.New(cfg.Verbose,
			redirection.New(cfg.Verbose, map[string]string{
				`/`: `/index.html`,
			}, cache.New(
				builder.New(cfg.Verbose, cfg.Mini,
					fileReader.New(cfg.Verbose, cfg.BasePath),
				),
			)),
		),
	))
}
