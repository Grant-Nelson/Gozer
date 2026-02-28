package irc

import (
	"fmt"

	"github.com/Grant-Nelson/Gozer/project/enums/basicType"
)

type (
	Type interface {

		// String gets a simple representation for this type.
		String() string

		// type is an empty method used to compile time type check that
		// only types duck-type to this interface.
		typ()
	}

	BasicType struct {
		Kind basicType.BasicType
	}

	SliceType struct {
		Elem Type
	}

	ArrayType struct {
		Size int
		Elem Type
	}

	PointerType struct {
		Elem Type
	}

	MapType struct {
		Key   Type
		Value Type
	}

	ChanType struct {
		Recv bool // indicates the channel can receive values
		Send bool // indicates the channel can send a value
		Elem Type // the type of values this channel contains
	}

	TupleType struct {
		Elems []Type
	}

	PackageType struct {
		Path string
	}

	// TODO: StructType
	// TODO: InterfaceType
	// TODO: FuncType
)

var (
	_ Type = (*BasicType)(nil)
	_ Type = (*SliceType)(nil)
	_ Type = (*ArrayType)(nil)
	_ Type = (*PointerType)(nil)
	_ Type = (*MapType)(nil)
	_ Type = (*TupleType)(nil)
	_ Type = (*PackageType)(nil)
)

func (t *BasicType) String() string   { return t.Kind.String() }
func (t *SliceType) String() string   { return fmt.Sprintf(`[]%s`, t.Elem) }
func (t *ArrayType) String() string   { return fmt.Sprintf(`[%d]%s`, t.Size, t.Elem) }
func (t *PointerType) String() string { return fmt.Sprintf(`*%s`, t.Elem) }
func (t *MapType) String() string     { return fmt.Sprintf(`map[%s]%s`, t.Key, t.Value) }
func (t *TupleType) String() string   { return fmt.Sprintf(`(%s)`, csvString(t.Elems)) }
func (t *PackageType) String() string { return fmt.Sprintf(`pkg(%s)`, t.Path) }

func (*BasicType) typ()   {}
func (*SliceType) typ()   {}
func (*ArrayType) typ()   {}
func (*PointerType) typ() {}
func (*MapType) typ()     {}
func (*TupleType) typ()   {}
func (*PackageType) typ() {}
