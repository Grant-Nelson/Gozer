package source

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Converter takes the given path for a directory (e.g. package) or file
// and converts the path to reference a new location.
//
// The path may be accompanied with the data for a file if the path
// is not a real location in the file system. This allows files be read from
// virtual or embedded file systems, or from an override buffer.
//
// The data may be present even if the file path is a real location in the
// file system to allow for preloading, caching, or sharing files.
// The data must be nil, a string, a []byte, or an io.Reader.
//
// If the path is to a location that should be skipped, the [ok] result
// should be set to false. This allows for augmentation files, caches,
// optional files, etc to be skipped if the given path doesn't contain
// such a file to be converted to.
type Converter func(path string, data any) (ok bool, pathOut string, dataOut any, err error)

// Then will append the given next source converters to be run after
// the current converter. The next converters are run in the order they
// are given until one returns not ok or an error, or they have all been run.
func (s Converter) Then(next ...Converter) Converter {
	if len(next) <= 0 {
		return s
	}
	return func(path string, data any) (bool, string, any, error) {
		ok, path, data, err := s(path, data)
		for _, n := range next {
			if !ok || err != nil {
				return ok, path, data, err
			}
			ok, path, data, err = n(path, data)
		}
		return ok, path, data, err
	}
}

// Noop is a source converter that just returns the path and data
// that is passed into it.
func Noop(path string, data any) (bool, string, any, error) {
	return true, path, data, nil
}

// Clean is a source converter that cleans the path that is passed into it.
func Clean(path string, data any) (bool, string, any, error) {
	return true, filepath.Clean(path), data, nil
}

// Abs is a source converter that gets the absolute path for
// the path that is passed into it.
func Abs(path string, data any) (bool, string, any, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, path, data, err
	}
	return true, absPath, data, nil
}

// PathRebase performs a simple string-wise match of a prefix path (oldBase)
// and replaces it with a new prefix path (newBase).
//
// If the oldBase is not prefixed then it will return false.
// If oldBase is empty, this will simply always concatenate the newBase.
// If the given path is not prefixed by the oldBase, the file will be skipped.
func PathRebase(oldBase, newBase string) Converter {
	if len(oldBase) <= 0 {
		return func(path string, data any) (bool, string, any, error) {
			return true, filepath.Join(newBase, path), data, nil
		}
	}

	return func(path string, data any) (bool, string, any, error) {
		if suffix, ok := strings.CutPrefix(path, oldBase); ok {
			return true, filepath.Join(newBase, suffix), data, nil
		}
		return false, path, data, nil
	}
}

// Override will return the value in the given maps as the data or new path
// if the given path is equal to a key in a map.
//
// Either map may be nil to not override the value.
// The new path is set after looking up the data such that any data override
// should be keyed with the given path, not the overridden path.
// The paths are string matched so the `pathOverrides` keys should be
// cleaned or normalized the same way as the paths being overridden.
func Override(pathOverrides map[string]string, dataOverrides map[string][]byte) Converter {
	return func(path string, data any) (bool, string, any, error) {
		if over, has := dataOverrides[path]; has {
			data = over
		}
		if over, has := pathOverrides[path]; has {
			path = over
		}
		return true, path, data, nil
	}
}

// Skipper will return false for [ok] when the given path is one of the
// paths to skip. The paths are string matched so the `skipPaths` should be
// cleaned or normalized the same way as the paths being skipped.
func Skipper(skipPaths ...string) Converter {
	skip := map[string]bool{}
	for _, s := range skipPaths {
		skip[s] = true
	}
	return func(path string, data any) (bool, string, any, error) {
		return !skip[path], path, data, nil
	}
}

// ToReader will get a reader from the given path and data.
// If the data is nil, the file at the given path will be read
// from the OS file system.
func ToReader(path string, data any) (io.Reader, error) {
	switch d := data.(type) {
	case nil:
		fd, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(fd), nil
	case string:
		return strings.NewReader(d), nil
	case []byte:
		return bytes.NewReader(d), nil
	case io.ReadCloser:
		// This assumes that if there is a Close method, then it should be
		// be read fully and closed before returning a reader of the data.
		fd, err := io.ReadAll(d)
		if err != nil {
			return nil, err
		}
		if err := d.Close(); err != nil {
			return nil, err
		}
		return bytes.NewReader(fd), nil
	case io.Reader:
		// This assumes that since there is no Close method, then it can
		// be used as is without pre-reading it into memory.
		return d, nil
	default:
		return nil, fmt.Errorf(`unexpected data type: %T`, d)
	}
}
