package constructs

type CObject struct {
	Name string
}

func NewObject(name string) *CObject {
	return &CObject{
		Name: name,
	}
}
