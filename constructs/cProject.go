package constructs

type CProject struct {
	Name     string
	Packages *CPackageSet
}

func NewProject(name string) *CProject {
	return &CProject{
		Name:     name,
		Packages: NewPackageSet(),
	}
}
