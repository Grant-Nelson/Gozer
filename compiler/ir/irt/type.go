package irt

// TODO: Finish moving types over to IR.

type Type interface {
	Underlying() Type
	String() string
}
