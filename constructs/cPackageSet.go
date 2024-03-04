package constructs

import "github.com/Snow-Gremlin/goToolbox/collections"

type CPackageSet interface {
	collections.Set[CPackage]
	TryGetByPath(path string) (CPackage, bool)
}
