package models

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
)

func addName(data map[string]any, name string) {
	if len(name) > 0 {
		data[`name`] = name
	}
}

func addIndex(data map[string]any, name string, index uint64) {
	if index > 0 {
		data[name] = index
	}
}

func addIndexed(data map[string]any, name string, indexed IndexedModel) {
	if indexed != nil {
		addIndex(data, name, indexed.Index())
	}
}

func addData[T any](data map[string]any, name string, list collections.List[T]) {
	if !list.Empty() {
		data[name] = list.ToSlice()
	}
}

func setIndices[T IndexedModel](list collections.List[T]) {
	var index uint64
	list.Enumerate().Foreach(func(i T) {
		index++ // increment first to be one based
		i.setIndex(index)
	})
}

func addIndices[T IndexedModel](data map[string]any, name string, list collections.List[T]) {
	if !list.Empty() {
		data[name] = enumerator.Select(list.Enumerate(),
			func(i T) uint64 { return i.Index() }).
			ToSlice()
	}
}

func setTypeIndices[T TypeModel](typeIndex *uint64, list collections.List[T]) {
	list.Enumerate().Foreach(func(i T) {
		*typeIndex++ // increment first to be one based
		i.setTypeIndex(*typeIndex)
	})
}

func addTypeIndices[T TypeModel](data map[string]any, name string, list collections.List[T]) {
	if !list.Empty() {
		data[name] = enumerator.Select(list.Enumerate(),
			func(i T) uint64 { return i.TypeIndex() }).
			ToSlice()
	}
}
