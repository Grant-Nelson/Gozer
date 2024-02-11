package readers

import "github.com/Snow-Gremlin/Gozer/constructs"

type Config struct {
	MainProject string
}

type Reader interface {
	Name() string
	Read(cfg *Config) (*constructs.CProject, error)
}
