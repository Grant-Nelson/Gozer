package cType

import "github.com/Snow-Gremlin/Gozer/constructs"

func New() constructs.CType {
	return &typeImp{}
}

type typeImp struct {
}
