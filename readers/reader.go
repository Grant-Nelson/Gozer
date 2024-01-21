package readers

import "github.com/Snow-Gremlin/Gozer/constructs"

type Config struct {
	ReaderName  string
	ProjectPath string
}

type Reader interface {
	Read(cfg *Config) (*constructs.CProject, error)
}
