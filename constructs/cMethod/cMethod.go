package cMethod

import (
	"github.com/Snow-Gremlin/Gozer/constructs"
)

func New() constructs.CMethod {
	return &methodImp{
		name: `unnamed`,
	}
}

type methodImp struct {
	name string
}

func (imp *methodImp) Name() string {
	return imp.name
}

func (imp *methodImp) SetName(name string) {
	imp.name = name
}
