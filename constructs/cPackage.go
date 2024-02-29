package constructs

type CPackage struct {
	Name    string
	Path    string
	Imports *CPackageSet
}

func NewPackage(path string) *CPackage {
	return &CPackage{
		Name:    ``,
		Path:    path,
		Imports: NewPackageSet(),
	}
}
