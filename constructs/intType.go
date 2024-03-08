package constructs

type IIntType interface {
	IType
	Signed() bool
	Size() int
	_intType()
}

type intTypeImp struct {
	name   string
	signed bool
	size   int
}

func (imp *intTypeImp) _type()    {}
func (imp *intTypeImp) _intType() {}

func (imp *intTypeImp) Signed() bool   { return imp.signed }
func (imp *intTypeImp) Size() int      { return imp.size }
func (imp *intTypeImp) String() string { return imp.name }

var (
	intInst    = &intTypeImp{name: `int`, signed: true, size: 64}
	int8Inst   = &intTypeImp{name: `int8`, signed: true, size: 8}
	int16Inst  = &intTypeImp{name: `int16`, signed: true, size: 16}
	int32Inst  = &intTypeImp{name: `int32`, signed: true, size: 32}
	int64Inst  = &intTypeImp{name: `int64`, signed: true, size: 64}
	uintInst   = &intTypeImp{name: `uint`, signed: false, size: 64}
	uint8Inst  = &intTypeImp{name: `uint8`, signed: false, size: 8}
	uint16Inst = &intTypeImp{name: `uint16`, signed: false, size: 16}
	uint32Inst = &intTypeImp{name: `uint32`, signed: false, size: 32}
	uint64Inst = &intTypeImp{name: `uint64`, signed: false, size: 64}
)

func Int() IIntType    { return intInst }
func Int8() IIntType   { return int8Inst }
func Int16() IIntType  { return int16Inst }
func Int32() IIntType  { return int32Inst }
func Int64() IIntType  { return int64Inst }
func Uint() IIntType   { return uintInst }
func Uint8() IIntType  { return uint8Inst }
func Uint16() IIntType { return uint16Inst }
func Uint32() IIntType { return uint32Inst }
func Uint64() IIntType { return uint64Inst }
