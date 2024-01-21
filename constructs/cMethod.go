package constructs

type CMethod struct {
	Name     string
	Receiver *CObject
}

func NewMethod(name string, receiver *CObject) *CMethod {
	return &CMethod{
		Name:     name,
		Receiver: receiver,
	}
}
