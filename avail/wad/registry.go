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

func RegisterNonPointer[T Waddable]() {
	regSingleton.RegisterNonPointer[T]()
}

const maxReg = 0xFFFF

type regEntry interface {
	read(r *Reader, count int) []Waddable
	write(w *Writer, wad Waddable)
}

type registry struct {
	lock    *sync.RWMutex
	lookup  map[string]uint16
	entries []regEntry
}

func newRegistry() *registry {
	return &registry{
		lock:   new(sync.RWMutex),
		lookup: make(map[string]uint16),
	}
}

func typeName(t reflect.Type) string { return t.PkgPath() + `.` + t.Name() }

func (reg *registry) cloneEntries() []regEntry {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	return slices.Clone(reg.entries)
}

func (reg *registry) cloneEntriesWithLookup() (map[string]uint16, []regEntry) {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	return maps.Clone(reg.lookup), slices.Clone(reg.entries)
}

func (reg *registry) registerEntry[T any](entry regEntry) {
	name := typeName(reflect.TypeFor[T]())

	reg.lock.Lock()
	defer reg.lock.Unlock()

	if _, has := reg.lookup[name]; has {
		panic(faults.New(`Waddable has already been registered`).
			With(`id`, name))
	}

	index := len(reg.entries)
	if index >= maxReg {
		panic(faults.New(`too many Waddables have been registered`).
			With(`maximum`, maxReg))
	}

	reg.lookup[name] = uint16(index + 1)
	reg.entries = append(reg.entries, entry)
}

type ptrRegEntry[P interface {
	*T
	Waddable
}, T any] struct{}

func (ptrRegEntry[P, T]) read(r *Reader, count int) []Waddable {
	s := make([]T, count)
	p := make([]Waddable, count)
	for i, e := range s {
		wad := P(&e)
		wad.UnWad(r)
		p[i] = wad
	}
	return p
}

func (ptrRegEntry[P, T]) write(w *Writer, wad Waddable) {
	wad.Wad(w)
}

func (reg *registry) Register[P interface {
	*T
	Waddable
}, T any]() {
	reg.registerEntry[P](ptrRegEntry[P, T]{})
}

type nonPtrRegEntry[T Waddable] struct{}

func (nonPtrRegEntry[T]) read(r *Reader, count int) []Waddable {
	s := make([]T, count)
	p := make([]Waddable, count)
	for i, e := range s {
		wad := Waddable(e)
		wad.UnWad(r)
		p[i] = wad
	}
	return p
}

func (nonPtrRegEntry[T]) write(w *Writer, wad Waddable) {
	wad.Wad(w)
}

func (reg *registry) RegisterNonPointer[T Waddable]() {
	reg.registerEntry[T](nonPtrRegEntry[T]{})
}
