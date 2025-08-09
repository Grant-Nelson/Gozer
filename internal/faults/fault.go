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

func From(err error) *Fault {
	return New(err.Error(), err)
}

func (f *Fault) Unwrap() []error {
	return f.inner
}

func (f *Fault) Data() map[string]any {
	return f.data
}

func (f *Fault) With(key string, value any) *Fault {
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
	parts := make([]string, 0, len(f.data)+1)
	parts = append(parts, f.msg)
	for k, v := range f.data {
		parts = append(parts, fmt.Sprintf(`\t%s: %v`, k, v))
	}
	sort.Strings(parts[1:])
	return strings.Join(parts, "\n")
}

func (f *Fault) String() string {
	return f.Error()
}
