package irt

type Type interface {
	Underlying() Type
	String() string
}
