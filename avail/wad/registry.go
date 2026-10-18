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
		entries   []*registerEntry
	}

	registerEntry struct {
		useCount     int
		nextBatch    int
		preallocated []Waddable
	}

	factory func(count int) []Waddable
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

	entries := make([]*registerEntry, len(r.factories))
	for i := range entries {
		entries[i] = new(registerEntry)
	}
	return &Registry{
		lock:      new(sync.Mutex),
		lookup:    maps.Clone(r.lookup),
		factories: slices.Clone(r.factories),
		entries:   entries,
	}
}

func (r *Registry) registerFactory[T Waddable](f factory) {
	name := typeName(reflect.TypeFor[T]())

	r.lock.Lock()
	defer r.lock.Unlock()

	if _, has := r.lookup[name]; has {
		panic(faults.New(`Waddable has already been registered`).
			With(`id`, name))
	}

	index := len(r.entries)
	if index >= maxReg {
		panic(faults.New(`too many Waddables have been registered`).
			With(`maximum`, maxReg))
	}

	r.lookup[name] = uint16(index + 1)
	r.factories = append(r.factories, f)
	r.entries = append(r.entries, new(registerEntry))
}

func (r *Registry) Register[P interface {
	*T
	Waddable
}, T any]() {
	r.registerFactory[P](func(count int) []Waddable {
		s := make([]T, count)
		p := make([]Waddable, count)
		for i, e := range s {
			p[i] = P(&e)
		}
		return p
	})
}

func (r *Registry) RegisterNonPointer[T Waddable]() {
	r.registerFactory[T](func(count int) []Waddable {
		s := make([]T, count)
		p := make([]Waddable, count)
		for i, e := range s {
			p[i] = e
		}
		return p
	})
}

func (r *Registry) toIndex(w Waddable) uint16 {
	v := reflect.ValueOf(w)
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return 0
	}
	name := typeName(v.Type())

	r.lock.Lock()
	defer r.lock.Unlock()

	if index, ok := r.lookup[name]; ok {
		r.entries[index].useCount++
		return index
	}
	panic(faults.New(`Unable to write an unregistered Waddable`).
		With(`type`, name))
}

func (r *Registry) fromIndex(i uint16) Waddable {
	if i <= 0 {
		return nil
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if int(i) > len(r.entries) {
		panic(faults.New(`Unable to read Waddable from registry with given index`).
			With(`index`, i))
	}

	e := r.entries[i-1]
	e.useCount++
	if count := len(e.preallocated); count > 0 {
		next := e.preallocated[0]
		if count > 1 {
			e.preallocated = e.preallocated[1:]
		} else if e.nextBatch > 0 {
			e.preallocated = r.factories[i-1](e.nextBatch)
			e.nextBatch = 0
		} else {
			e.preallocated = nil
		}
		return next
	}
	return r.factories[i-1](1)[0]
}

func (r *Registry) getUseCountWithIndex(index uint16) int {
	if int(index) >= len(r.entries) {
		panic(faults.New(`can not get the "use count" with an out-of-bounds index`).
			With(`index`, index).
			With(`length`, len(r.entries)))
	}
	return r.entries[index].useCount
}

func (r *Registry) GetUseCount[T Waddable]() int {
	name := typeName(reflect.TypeFor[T]())

	r.lock.Lock()
	defer r.lock.Unlock()

	if index, ok := r.lookup[name]; ok {
		return r.getUseCountWithIndex(index)
	}
	panic(faults.New(`can not get the "use count" with an unregistered Waddable`).
		With(`id`, name))
}

func (r *Registry) preallocateWithIndex(index uint16, count int) {
	if int(index) >= len(r.entries) {
		panic(faults.New(`can not preallocate types with an out-of-bounds index`).
			With(`index`, index).
			With(`length`, len(r.entries)))
	}
	e := r.entries[index]
	if pre := len(e.preallocated); pre > 0 {
		e.nextBatch = max(e.nextBatch, count-pre)
		return
	}
	e.preallocated = r.factories[index](count)
}

func (r *Registry) Preallocate[T Waddable](count int) {
	if count <= 0 {
		return
	}
	name := typeName(reflect.TypeFor[T]())

	r.lock.Lock()
	defer r.lock.Unlock()

	if index, ok := r.lookup[name]; ok {
		r.preallocateWithIndex(index, count)
		return
	}
	panic(faults.New(`can not preallocate an unregistered Waddable`).
		With(`id`, name))
}
