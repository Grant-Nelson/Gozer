package constructs

import "github.com/Snow-Gremlin/goToolbox/collections"

type CObjectSet interface {
	collections.Set[CObject]
	TryGetByName(name string) (CObject, bool)
}
