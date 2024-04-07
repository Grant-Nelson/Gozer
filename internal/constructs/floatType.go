package constructs

type IFloatType interface {
	IType
	Size() int
	_floatType()
}

type floatTypeImp struct {
	name string
	size int
}

func (imp *floatTypeImp) _type()      {}
func (imp *floatTypeImp) _floatType() {}

func (imp *floatTypeImp) Size() int      { return imp.size }
func (imp *floatTypeImp) String() string { return imp.name }

var (
	float32Inst = &floatTypeImp{name: `float32`, size: 32}
	float64Inst = &floatTypeImp{name: `float64`, size: 64}
)

func Float32() IFloatType { return float32Inst }
func Float64() IFloatType { return float64Inst }
