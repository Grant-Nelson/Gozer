package cObject

import (
	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/constructs/cDirectives"
	"github.com/Snow-Gremlin/Gozer/constructs/cMethodSet"
)

func New() constructs.CObject {
	return &objectImp{
		name:       `unnamed`,
		directives: cDirectives.New(),
		methods:    cMethodSet.New(),
	}
}

type objectImp struct {
	name       string
	directives constructs.CDirectives
	methods    constructs.CMethodSet
}

func (imp *objectImp) Name() string {
	return imp.name
}

func (imp *objectImp) SetName(name string) {
	imp.name = name
}

func (imp *objectImp) Directives() constructs.CDirectives {
	return imp.directives
}

func (imp *objectImp) Methods() constructs.CMethodSet {
	return imp.methods
}
