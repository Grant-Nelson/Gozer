package cPackage

import (
	"github.com/Snow-Gremlin/Gozer/constructs"
	"github.com/Snow-Gremlin/Gozer/constructs/cMethodSet"
	"github.com/Snow-Gremlin/Gozer/constructs/cPackageSet"
)

func New() constructs.CPackage {
	return &packageImp{
		name:    ``,
		path:    ``,
		imports: cPackageSet.New(),
		methods: cMethodSet.New(),
	}
}

type packageImp struct {
	name    string
	path    string
	imports constructs.CPackageSet
	methods constructs.CMethodSet
}

func (imp *packageImp) Name() string {
	return imp.name
}

func (imp *packageImp) SetName(name string) {
	imp.name = name
}

func (imp *packageImp) Path() string {
	return imp.path
}

func (imp *packageImp) SetPath(path string) {
	imp.path = path
}

func (imp *packageImp) Imports() constructs.CPackageSet {
	return imp.imports
}

func (imp *packageImp) Methods() constructs.CMethodSet {
	return imp.methods
}
