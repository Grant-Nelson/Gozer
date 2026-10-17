package wad

import (
	"reflect"
	"sync"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

var DefaultRegistry = NewRegistry()

const maxReg = 0xFFFF

type (
	Registry struct {
		lock    sync.RWMutex
		lookup  map[string]uint16
		entries []*registerEntry
	}

	registerEntry struct {
		useCount     int
		nextBatch    int
		preallocated []Waddable
		factory      factory
	}

	factory func(count int) []Waddable
)

func NewRegistry() *Registry {
	return &Registry{
		lookup:  map[string]uint16{},
		entries: []*registerEntry{},
	}
}

func typeName(t reflect.Type) string { return t.PkgPath() + `.` + t.Name() }

func (r *Registry) registerFactory[T Waddable](f factory) {
	name := typeName(reflect.TypeFor[T]())
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
	r.entries = append(r.entries, &registerEntry{
		factory: f,
	})
}

func Register[P interface {
	*T
	Waddable
}, T any]() {
	DefaultRegistry.Register[P, T]()
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

func RegisterNP[T Waddable]() { DefaultRegistry.RegisterNP[T]() }

func (r *Registry) RegisterNP[T Waddable]() {
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
	if v := reflect.ValueOf(w); !v.IsNil() {
		name := typeName(v.Type())
		if index, ok := r.lookup[name]; ok {
			r.entries[index].useCount++
			return index
		}
		panic(faults.New(`Unable to write an unregistered Waddable`).
			With(`type`, name))
	}
	return 0
}

func (r *Registry) fromIndex(i uint16) Waddable {
	if i <= 0 {
		return nil
	}
	if int(i) > len(r.entries) {
		panic(faults.New(`Unable to read Waddable from registry with given index`).
			With(`index`, i))
	}

	e := r.entries[i-1]
	e.useCount++
	return e.getNext()
}

func (r *Registry) Preallocate[T Waddable](count int) {
	name := typeName(reflect.TypeFor[T]())
	if index, ok := r.lookup[name]; ok {
		r.entries[index].preallocate(count)
	}
	panic(faults.New(`can not preallocate an unregistered Waddable`).
		With(`id`, name))
}

func (re *registerEntry) reset() {
	re.useCount = 0
	re.nextBatch = 0
	re.preallocated = nil
}

func (re *registerEntry) preallocate(count int) {
	if count <= 0 {
		return
	}
	if pre := len(re.preallocated); pre > 0 {
		re.nextBatch = max(re.nextBatch, count-pre)
		return
	}
	re.preallocated = re.factory(count)
}

func (re *registerEntry) getNext() Waddable {
	if count := len(re.preallocated); count > 0 {
		next := re.preallocated[0]
		if count > 1 {
			re.preallocated = re.preallocated[1:]
		} else if re.nextBatch > 0 {
			re.preallocated = re.factory(re.nextBatch)
			re.nextBatch = 0
		} else {
			re.preallocated = nil
		}
		return next
	}
	return re.factory(1)[0]
}
