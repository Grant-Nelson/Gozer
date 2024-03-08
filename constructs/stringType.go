package constructs

type IStringType interface {
	IType
	_stringType()
}

type stringTypeImp struct{}

func (imp *stringTypeImp) _type()       {}
func (imp *stringTypeImp) _stringType() {}

func (imp *stringTypeImp) String() string { return `string` }

var stringInst = &stringTypeImp{}

func String() IStringType { return stringInst }
