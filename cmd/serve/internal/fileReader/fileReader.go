package fileReader

import (
	"errors"

	"github.com/Grant-Nelson/Gozer/cmd/serve/internal/fetcher"
)

func New(verbose bool, basePath string) fetcher.Fetcher {
	return &fileReaderImp{
		verbose:  verbose,
		basePath: basePath,
	}
}

type fileReaderImp struct {
	verbose  bool
	basePath string
}

func (f *fileReaderImp) Fetch(path string) ([]byte, error) {
	// TODO: Implement
	return nil, errors.New(`not implemented`)
}
