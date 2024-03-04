package constructs

type CPackage interface {
	Named
	Path() string
	SetPath(path string)
	Imports() CPackageSet
	Methods() CMethodSet
}
