package models

import (
	"reflect"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
)

func addVal(data map[string]any, name string, value any) {
	if !reflect.ValueOf(value).IsZero() {
		data[name] = value
	}
}

func addList[T any](data map[string]any, name string, list collections.ReadonlyList[T]) {
	if !list.Empty() {
		data[name] = list.ToSlice()
	}
}

func addIds[T Identifier](data map[string]any, name string, list collections.ReadonlyList[T]) {
	if !list.Empty() {
		data[name] = enumerator.Select(list.Enumerate(),
			func(i T) uint64 { return i.Id() }).
			ToSlice()
	}
}
