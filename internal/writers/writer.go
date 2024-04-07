package writers

import "github.com/Snow-Gremlin/Gozer/internal/constructs"

type Config struct {
	OutPath string
}

type Writer interface {
	Name() string
	Write(cfg *Config, proj constructs.IProject) error
}
