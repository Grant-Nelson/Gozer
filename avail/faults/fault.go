package faults

import (
	"fmt"
	"maps"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
)

// disableRecovers is used to disable the [Recover] methods.
// This is useful when trying to find an error and the stack trace
// keeps being consumed by a recover.
const disableRecovers = false

type Fault struct {
	msg    string
	inner  []error
	data   map[string]any
	file   string
	lineNo int
}

var _ error = (*Fault)(nil)

func New(msg string, inner ...error) *Fault {
	f := &Fault{
		msg:   msg,
		inner: inner,
	}
	f.SetCurLoc()
	return f
}

func From(r any) *Fault {
	switch t := r.(type) {
	case *Fault:
		return t
	case *ErrGroup:
		return New(`error group`, t)
	case error:
		return New(t.Error(), t)
	case string:
		return New(t)
	default:
		return New(fmt.Sprintf(`error for %v`, t))
	}
}

func Recover(pe *error) {
	if disableRecovers {
		return
	}
	if pe != nil {
		if r := recover(); r != nil {
			*pe = From(r)
		}
	}
}

func (f *Fault) HasLoc() bool {
	return f != nil && len(f.file) > 0
}

func (f *Fault) Loc() (file string, lineNo int) {
	if f == nil {
		return ``, 0
	}
	return f.file, f.lineNo
}

func (f *Fault) SetLoc(file string, lineNo int) *Fault {
	if f != nil {
		f.file = file
		f.lineNo = lineNo
	}
	return f
}

// SetCurLoc will set the current location to this fault.
//
// This will attempt to find the calling location by filtering based
// on path suffixes so some paths will not show up.
// This will return true if a location is found, otherwise it will
// give up after 10 attempts, will set an empty file path and return false.
func (f *Fault) SetCurLoc() bool {
	if f == nil {
		return false
	}
	f.file = ``
	f.lineNo = 0

	const maxAttempts = 10
	for i := range maxAttempts {
		if _, file, line, ok := runtime.Caller(i); ok {
			dir := filepath.ToSlash(filepath.Dir(file))
			if strings.HasSuffix(dir, `/avail/faults`) ||
				strings.HasSuffix(dir, `/runtime`) ||
				strings.HasSuffix(dir, `/avail/assert`) {
				continue
			}

			f.file = file
			f.lineNo = line
			return true
		}
	}
	return false
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
	case fmt.Stringer:
		f.data[key] = t.String()
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

func (f *Fault) WithNonZero[T comparable](key string, value T) *Fault {
	var zero T
	if value != zero {
		return f.With(key, value)
	}
	return f
}

func (f *Fault) WithF(key, format string, args ...any) *Fault {
	return f.With(key, fmt.Sprintf(format, args...))
}

func (f *Fault) WithTypeFor[T any](key string) *Fault {
	var zero T
	return f.WithTypeOf(key, zero)
}

func (f *Fault) WithTypeOf[T any](key string, value T) *Fault {
	return f.WithF(key, `%T`, value)
}

func (f *Fault) Error() string {
	if f == nil {
		return `nil`
	}

	const indent = "\t"

	buf := &strings.Builder{}
	write := func(s string) {
		if _, err := buf.WriteString(s); err != nil {
			panic(err)
		}
	}

	if len(f.msg) > 0 {
		write(f.msg)
	} else {
		write(`nonspecific error`)
	}
	if len(f.file) > 0 {
		write(fmt.Sprintf(` @ %s`, f.file))
		if f.lineNo >= 0 {
			write(fmt.Sprintf(`:%d`, f.lineNo))
		}
	}

	dataError := map[error]bool{}
	keys := slices.Collect(maps.Keys(f.data))
	sort.Strings(keys)
	keyWidth := 0
	for _, k := range keys {
		keyWidth = max(keyWidth, len(k))
	}
	spacing := strings.Repeat(` `, keyWidth+2)
	for _, k := range keys {
		v := f.data[k]
		if e, ok := v.(error); ok {
			dataError[e] = true
		}

		val := fmt.Sprint(v)
		val = strings.ReplaceAll(val, "\n", "\n"+indent+spacing)
		key := k + `: `
		key += strings.Repeat(" ", keyWidth-len(key))
		write("\n" + indent + key + val)
	}

	for _, in := range f.inner {
		if !dataError[in] {
			write("\n" + indent + in.Error())
		}
	}
	return buf.String()
}

func (f *Fault) String() string {
	return f.Error()
}
