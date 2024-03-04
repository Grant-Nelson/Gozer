package cDirectives

import (
	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/goToolbox/collections/set"
)

func New() constructs.CDirectives {
	return set.New[string]()
}
