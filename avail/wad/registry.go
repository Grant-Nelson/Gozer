package wad

import (
	"maps"
	"reflect"
	"slices"
	"sync"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

var regSingleton = newRegistry()

func Register[P interface {
	*T
	Waddable
}, T any]() {
	regSingleton.Register[P, T]()
}

const maxReg = 0xFFFF

type ReadFn func(*Reader) any
type WriteFn func(*Writer, any)

type registry struct {
	lock    *sync.RWMutex
	lookup  map[string]uint16
	readers []ReadFn
	writers []WriteFn
}

func newRegistry() *registry {
	return &registry{
		lock:   new(sync.RWMutex),
		lookup: make(map[string]uint16),
	}
}

func typeName(t reflect.Type) string {
	return t.PkgPath() + `.` + t.Name()
}

func (reg *registry) cloneReaders() []ReadFn {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	return slices.Clone(reg.readers)
}

func (reg *registry) cloneWriters() (map[string]uint16, []WriteFn) {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	return maps.Clone(reg.lookup), slices.Clone(reg.writers)
}

func (reg *registry) RegisterCustom[T any](reader ReadFn, writer WriteFn) {
	name := typeName(reflect.TypeFor[T]())

	reg.lock.Lock()
	defer reg.lock.Unlock()

	if _, has := reg.lookup[name]; has {
		panic(faults.New(`Waddable has already been registered`).
			With(`id`, name))
	}

	index := len(reg.readers)
	if index >= maxReg {
		panic(faults.New(`too many Waddables have been registered`).
			With(`maximum`, maxReg))
	}

	reg.lookup[name] = uint16(index + 1)
	reg.readers = append(reg.readers, reader)
	reg.writers = append(reg.writers, writer)
}

func (reg *registry) Register[P interface {
	*T
	Waddable
}, T any]() {
	reg.RegisterCustom[P](
		func(r *Reader) any {
			p := P(new(T))
			p.UnWad(r)
			return p
		},
		func(w *Writer, value any) {
			value.(P).Wad(w)
		})

	reg.RegisterCustom[T](
		func(r *Reader) any {
			p := P(new(T))
			p.UnWad(r)
			return *p
		},
		func(w *Writer, value any) {
			e := value.(T)
			p := P(&e)
			p.Wad(w)
		})
}
