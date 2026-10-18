package wad

import (
	"maps"
	"reflect"
	"slices"
	"sync"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

var defaultRegistry = newRegistry()

func Register[P interface {
	*T
	Waddable
}, T any]() {
	defaultRegistry.Register[P, T]()
}

func RegisterNonPointer[T Waddable]() {
	defaultRegistry.RegisterNonPointer[T]()
}

const maxReg = 0xFFFF

type (
	Registry struct {
		lock      *sync.Mutex
		lookup    map[string]uint16
		factories []factory
	}

	factory func(count int) any
)

func newRegistry() *Registry {
	return &Registry{
		lock:   new(sync.Mutex),
		lookup: make(map[string]uint16),
	}
}

func typeName(t reflect.Type) string { return t.PkgPath() + `.` + t.Name() }

func (r *Registry) clone() *Registry {
	r.lock.Lock()
	defer r.lock.Unlock()
	return &Registry{
		lock:      new(sync.Mutex),
		lookup:    maps.Clone(r.lookup),
		factories: slices.Clone(r.factories),
	}
}

func (r *Registry) registerFactory[T any](f factory) {
	name := typeName(reflect.TypeFor[T]())

	r.lock.Lock()
	defer r.lock.Unlock()

	if _, has := r.lookup[name]; has {
		panic(faults.New(`Waddable has already been registered`).
			With(`id`, name))
	}

	index := len(r.factories)
	if index >= maxReg {
		panic(faults.New(`too many Waddables have been registered`).
			With(`maximum`, maxReg))
	}

	r.lookup[name] = uint16(index + 1)
	r.factories = append(r.factories, f)
}

func (r *Registry) Register[P *T, T any]() {
	r.registerFactory[P](func(count int) any {
		s := make([]T, count)
		p := make([]P, count)
		for i, e := range s {
			p[i] = P(&e)
		}
		return p
	})
}

func (r *Registry) RegisterNonPointer[T any]() {
	r.registerFactory[T](func(count int) any {
		return make([]T, count)
	})
}

func (r *Registry) toIndex[T any](w any) uint16 {
	v := reflect.ValueOf(w)
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return 0
	}
	name := typeName(v.Type())

	r.lock.Lock()
	defer r.lock.Unlock()

	if index, ok := r.lookup[name]; ok {
		return index
	}
	panic(faults.New(`Unable to write an unregistered Waddable`).
		With(`type`, name))
}

func (r *Registry) fromIndex[T any](i uint16, count int) []T {
	if i <= 0 {
		return nil
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if int(i) > len(r.factories) {
		panic(faults.New(`Unable to read Waddable from registry with given index`).
			With(`index`, i))
	}

	return r.factories[i-1](count).([]T)
}
