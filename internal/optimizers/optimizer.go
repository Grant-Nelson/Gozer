package optimizers

import "github.com/Snow-Gremlin/Gozer/internal/constructs"

type Factory func() Optimizer

type Config struct{}

type Optimizer interface {
	Name() string
	Perform(cfg *Config, proj constructs.IProject) error
}
