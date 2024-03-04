package constructs

type CProject interface {
	Named
	Packages() CPackageSet
}
