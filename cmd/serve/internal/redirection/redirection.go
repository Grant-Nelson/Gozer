package redirection

import (
	"path"
	"strings"

	"github.com/Grant-Nelson/Gozer/cmd/serve/internal/fetcher"
)

func New(verbose bool, mapping map[string]string, next fetcher.Fetcher) fetcher.Fetcher {
	return &redirectionImp{
		verbose: verbose,
		mapping: mapping,
		next:    next,
	}
}

type redirectionImp struct {
	verbose bool
	mapping map[string]string
	next    fetcher.Fetcher
}

func (r *redirectionImp) Fetch(p string) ([]byte, error) {
	if !strings.HasPrefix(p, `/`) {
		p = `/` + p
	}
	p = path.Clean(p)
	if rePath, found := r.mapping[p]; found {
		p = rePath
	}
	return r.next.Fetch(p)
}
