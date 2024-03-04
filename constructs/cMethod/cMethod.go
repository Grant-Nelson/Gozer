package constructs

import (
	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/constructs/cDirectives"
)

func New() constructs.CMethod {
	return &methodImp{
		name:       `unnamed`,
		directives: cDirectives.New(),
	}
}

type methodImp struct {
	name       string
	directives constructs.CDirectives
}

func (imp *methodImp) Name() string {
	return imp.name
}

func (imp *methodImp) SetName(name string) {
	imp.name = name
}

func (imp *methodImp) Directives() constructs.CDirectives {
	return imp.directives
}
