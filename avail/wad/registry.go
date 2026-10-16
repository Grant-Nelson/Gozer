package wad

import (
	"reflect"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

var DefaultRegistry = NewRegistry()

// TODO: FIX
//func Register[T Waddable]() { DefaultRegistry.Register[T]() }

const maxReg = 0xFFFF

type (
	Registry struct {
		lookup  map[string]uint16
		entries []*regEntry
	}

	regEntry struct {
		nextBatch    int
		preallocated []Waddable
		factory      factory
	}

	factory interface {
		Create(count int) []Waddable
	}
)

func NewRegistry() *Registry {
	return &Registry{
		lookup:  map[string]uint16{},
		entries: []*regEntry{},
	}
}

func typeName(t reflect.Type) string { return t.PkgPath() + `.` + t.Name() }

func (re *regEntry) getNext() Waddable {
	if count := len(re.preallocated); count > 0 {
		next := re.preallocated[0]
		if count > 1 {
			re.preallocated = re.preallocated[1:]
		} else if re.nextBatch > 0 {
			re.preallocated = re.factory.Create(re.nextBatch)
			re.nextBatch = 0
		} else {
			re.preallocated = nil
		}
		return next
	}
	return re.factory.Create(1)[0]
}

func (re *regEntry) preallocate(count int) {
	if len(re.preallocated) > 0 {
		re.nextBatch += count
		return
	}
	re.preallocated = re.factory.Create(count)
}

func (r *Registry) toIndex(w Waddable) uint16 {
	if v := reflect.ValueOf(w); !v.IsNil() {
		name := typeName(v.Type())
		if index, ok := r.lookup[name]; ok {
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
	if int(i) > len(r.entries) {
		panic(faults.New(`Unable to read Waddable from registry with given index`).
			With(`index`, i))
	}
	return r.entries[i-1].getNext()
}

func (r *Registry) registerFactory[T Waddable](f factory) {
	name := typeName(reflect.TypeFor[T]())
	if _, has := r.lookup[name]; !has {
		index := len(r.entries)
		if index >= maxReg {
			panic(faults.New(`too many Waddables have been registered`).
				With(`maximum`, maxReg))
		}
		r.lookup[name] = uint16(index + 1)
		r.entries = append(r.entries, &regEntry{
			factory: f,
		})
	}
}
