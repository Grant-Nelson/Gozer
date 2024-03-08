package readers

import "github.com/Snow-Gremlin/Gozer/constructs"

type Config struct {
	MainPackageDir string
}

type Reader interface {
	Name() string
	Read(cfg *Config) (constructs.IProject, error)
}
