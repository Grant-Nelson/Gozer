package typeMod

import "strings"

type TypeMod int

const (
	None TypeMod = 0

	Const TypeMod = 1 << iota
	Var
	Final
	Readonly
	Static
	Exported
	Internal
	Public
	Private
)

const max = Private

var names = []string{
	`const`,
	`var`,
	`final`,
	`readonly`,
	`static`,
	`exported`,
	`internal`,
	`public`,
	`private`,
}

func (m TypeMod) Valid() bool {
	return m >= None && m < (max<<1)
}

func (m TypeMod) String() string {
	if !m.Valid() {
		return `UnknownTypeMod`
	}
	parts := []string{}
	for _, n := range names {
		if m <= 0 {
			break
		}
		if m&1 != 0 {
			parts = append(parts, n)
		}
		m >>= 1
	}
	if len(parts) <= 0 {
		return `none`
	}
	return strings.Join(parts, `|`)
}
