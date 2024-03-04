package constructs

import "github.com/Snow-Gremlin/goToolbox/collections"

type CMethodSet interface {
	collections.Set[CMethod]
	TryGetByName(name string) (CMethod, bool)
}
