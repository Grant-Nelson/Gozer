package constructs

type ITypeParameter interface {
	INamed
	Constraint() IConstraint
}

type typeParameterImp struct {
	namedImp
	constraint IConstraint
}

func (imp *typeParameterImp) Constraint() IConstraint { return imp.constraint }

func NewTypeParameter(name string, constraint IConstraint) ITypeParameter {
	return &typeParameterImp{
		namedImp:   newName(name),
		constraint: constraint,
	}
}
