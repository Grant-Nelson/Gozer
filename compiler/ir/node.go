package ir

import "go/token"

type Node interface {
	Pos() token.Pos
}
