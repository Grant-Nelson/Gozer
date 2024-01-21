package constructs

type CPackage struct {
	Name string
}

func NewPackage(name string) *CPackage {
	return &CPackage{
		Name: name,
	}
}
