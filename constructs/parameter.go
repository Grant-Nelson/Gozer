package constructs

type IParameter interface {
	INamed
	Type() IType
}

type parameterImp struct {
	namedImp
	t IType
}

func (imp *parameterImp) Type() IType { return imp.t }

func NewParameter(name string, t IType) IParameter {
	return &parameterImp{
		namedImp: newName(name),
		t:        t,
	}
}
