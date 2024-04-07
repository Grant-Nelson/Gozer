package constructs

type IPointerType interface {
	IType
	Inner() IType
	Unsafe() bool
	_pointerType()
}

type pointerTypeImp struct {
	inner IType
}

func (imp *pointerTypeImp) _type()        {}
func (imp *pointerTypeImp) _pointerType() {}

func (imp *pointerTypeImp) Inner() IType { return imp.inner }
func (imp *pointerTypeImp) Unsafe() bool { return imp.inner == nil }

func (imp *pointerTypeImp) String() string {
	if imp.Unsafe() {
		return `unsafe`
	}
	return `*` + imp.inner.String()
}

var unsafePointerInst = &pointerTypeImp{inner: nil}

func UnsafePointer() IPointerType { return unsafePointerInst }
func Pointer(inner IType) IPointerType {
	if inner == nil {
		return UnsafePointer()
	}
	return &pointerTypeImp{inner: inner}
}
