package faults

import (
	"fmt"
	"sort"
	"strings"
)

type Fault struct {
	msg   string
	inner []error
	data  map[string]any
}

func New(msg string, inner ...error) *Fault {
	return &Fault{
		msg:   msg,
		inner: inner,
	}
}

func From(r any) *Fault {
	switch t := r.(type) {
	case *Fault:
		return t
	case *Group:
		return New(`several errors`, t.err...)
	case error:
		return New(t.Error(), t)
	case string:
		return New(t)
	default:
		return New(fmt.Sprintf(`error for %v`, t))
	}
}

func Recover(pe *error) {
	if r := recover(); r != nil {
		*pe = From(r)
	}
}

func (f *Fault) Unwrap() []error {
	if f == nil {
		return nil
	}
	return f.inner
}

func (f *Fault) Data() map[string]any {
	if f == nil {
		return nil
	}
	return f.data
}

func (f *Fault) With(key string, value any) *Fault {
	if f == nil {
		return nil
	}
	if len(f.data) <= 0 {
		f.data = map[string]any{}
	}
	switch t := value.(type) {
	case error:
		f.inner = append(f.inner, t)
		f.data[key] = t
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64, complex64, complex128,
		string:
		f.data[key] = t
	default:
		f.data[key] = fmt.Sprint(t)
	}
	return f
}

func (f *Fault) Error() string {
	if f == nil {
		return `nil`
	}
	parts := make([]string, 0, len(f.data)+1)
	parts = append(parts, f.msg)
	for k, v := range f.data {
		parts = append(parts, fmt.Sprintf("\t%s: %v", k, v))
	}
	sort.Strings(parts[1:])
	return strings.Join(parts, "\n")
}

func (f *Fault) String() string {
	return f.Error()
}
