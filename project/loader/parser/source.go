package parser

import "strings"

// SourceConverter takes the given path for a directory (e.g. package) or file
// and converts the path to reference a new location.
//
// The path may be accompanied with the data for a file if the path
// is not a real location in the file system. This allows files be read from
// virtual or embedded file systems, or from an override buffer.
// The data may be present even if the file path is a real location in the
// file system to allow for preloading, caching, or sharing files.
//
// If the path is to a location that should be skipped, the [ok] result
// should be set to false. This allows for augmentation files, caches,
// optional files, etc to be skipped if the given path doesn't contain
// such a file to be converted to.
type SourceConverter func(path string, data any) (ok bool, newPath string, newData any, err error)

// Then will append the given next source converters to be run after
// the current converter. The next converters are run in the order they
// are given until one returns not ok or an error, or they have all been run.
func (s SourceConverter) Then(next ...SourceConverter) SourceConverter {
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

// PathRebase performs a simple string-wise match of a prefix path (oldBase)
// and replaces it with a new prefix path (newBase).
// If the oldBase id not prefixed then it will return false.
// If oldBase is empty, this will simply always concatenate the newBase.
// If the given path is not prefixed by the oldBase, the file will be skipped.
func PathRebase(oldBase, newBase string) SourceConverter {
	if len(oldBase) <= 0 {
		return func(path string, data any) (bool, string, any, error) {
			return true, newBase + path, data, nil
		}
	}
	return func(path string, data any) (bool, string, any, error) {
		if suffix, ok := strings.CutPrefix(path, oldBase); ok {
			return true, newBase + suffix, data, nil
		}
		return false, path, data, nil
	}
}
