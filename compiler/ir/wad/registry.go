package wad

import (
	"reflect"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

const maxReg = 0x7FFF

type Registry struct {
	m map[string]uint16
	c []func() Waddable
}

func NewRegistry() *Registry {
	return &Registry{
		m: map[string]uint16{},
		c: []func() Waddable{},
	}
}

var DefaultRegistry = NewRegistry()

func Register[T Waddable]() { DefaultRegistry.Register[T]() }

func typeName(t reflect.Type) string { return t.PkgPath() + `.` + t.Name() }

func (r *Registry) Register[T Waddable]() {
	name := typeName(reflect.TypeFor[T]())
	if _, has := r.m[name]; !has {
		index := uint16(len(r.c) + 1)
		if index >= maxReg {
			panic(faults.New(`too many Waddables have been registered`).
				With(`maximum`, maxReg))
		}
		r.m[name] = index
		r.c = append(r.c,
			func() Waddable {
				var zero T
				return zero
			})
	}
}

func (r *Registry) getIndex(w Waddable) uint16 {
	if v := reflect.ValueOf(w); !v.IsNil() {
		name := typeName(v.Type())
		if index, ok := r.m[name]; ok {
			return index
		}
		panic(faults.New(`Unable to write an unregistered Waddable`).
			With(`type`, name))
	}
	return 0
}

func (r *Registry) fromIndex(i uint16) Waddable {
	if i == 0 {
		return nil
	}
	if int(i) > len(r.c) {
		panic(faults.New(`Unable to read Waddable from registry with given index`).
			With(`index`, i))
	}
	return r.c[i-1]()
}
