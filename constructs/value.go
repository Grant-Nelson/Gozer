package constructs

type IValue interface {
	INamed
	Constant() bool
	SetConstant(constant bool)
	Type() IType
	SetType(t IType)
	Assignment() IExpression
	SetAssignment(exp IExpression)
	_valueConstruct()
}

type valueImp struct {
	namedImp
	constant bool
	t        IType
	exp      IExpression
}

func (imp *valueImp) _valueConstruct() {}

func (imp *valueImp) Constant() bool            { return imp.constant }
func (imp *valueImp) SetConstant(constant bool) { imp.constant = constant }

func (imp *valueImp) Type() IType     { return imp.t }
func (imp *valueImp) SetType(t IType) { imp.t = t }

func (imp *valueImp) Assignment() IExpression       { return imp.exp }
func (imp *valueImp) SetAssignment(exp IExpression) { imp.exp = exp }

func NewVariable(name string, t IType) IValue {
	return &valueImp{
		namedImp: newName(name),
		constant: false,
		t:        t,
		exp:      nil,
	}
}

func NewConstant(name string, t IType) IValue {
	v := NewVariable(name, t)
	v.SetConstant(true)
	return v
}
