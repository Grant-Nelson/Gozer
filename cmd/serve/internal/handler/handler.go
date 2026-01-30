package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Grant-Nelson/Gozer/cmd/serve/internal/fetcher"
)

func New(verbose bool, fetcher fetcher.Fetcher) http.Handler {
	return &handlerImp{
		verbose: verbose,
		fetcher: fetcher,
	}
}

type handlerImp struct {
	verbose bool
	fetcher fetcher.Fetcher
}

func (h *handlerImp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := fmt.Sprintf(`Server only allows GET but got %s`, r.Method)
		http.Error(w, err, http.StatusNotImplemented)
		return
	}

	path := r.URL.Path
	if h.verbose {
		log.Printf(`Request for %q`, path)
	}

	content, err := h.fetcher.Fetch(path)
	if err != nil {
		if h.verbose {
			log.Printf(`Failed to get %q: %v`, path, err)
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if _, err = w.Write(content); err != nil {
		log.Printf(`Error writing response for %q: %v`, path, err)
	}
}
