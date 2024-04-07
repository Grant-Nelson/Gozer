package constructs

type IBoolType interface {
	IType
	_boolType()
}

type boolTypeImp struct{}

func (imp *boolTypeImp) _type()     {}
func (imp *boolTypeImp) _boolType() {}

func (imp *boolTypeImp) String() string { return `bool` }

var boolInst = &boolTypeImp{}

func Bool() IBoolType { return boolInst }
