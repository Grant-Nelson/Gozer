package writers

import "github.com/Snow-Gremlin/Gozer/constructs"

type Config struct {
	WriterName string
	OutPath    string
}

type Writer interface {
	Name() string
	Write(cfg *Config, proj *constructs.CProject) error
}
