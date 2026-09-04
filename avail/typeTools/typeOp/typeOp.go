package typeOp

import "strings"

type Op uint64

const (
	None       Op = 0         // still has zero() and new()
	Add        Op = 1 << iota // x+y, x+=y for (numeric, string)
	Arith                     // -x, x-y, x*y, x/y, x-=y, x*=y, x/=y, x++, x--
	Bitwise                   // x&y, x|y, ^x, x^y, x&^y, x<<y, x>>y, x&=y, x|=y, x^=y, x&^=y, x<<=y, x>>=y
	ByteSlice                 // []byte(x), copyTo(x, y)
	Cap                       // cap(x)
	Clear                     // clear(x)
	Comparable                // x==y, x!=y
	Complex                   // complex(x,y)
	Deref                     // *x
	GetIndex                  // z=x[y]
	GetIndex2                 // z,ok=x[y]
	IsNil                     // x==nil
	Len                       // len(x)
	Make                      // make(x, len) for (slice, map, channel)
	Make3                     // make3(x, len, cap) for (slice)
	Mod                       // x%y, x%=y
	Orderable                 // x<y, x<=y, x>y, x>=y
	Range                     // for _=range x, for y=range x
	Range2                    // for _,_=range x, for y,_=range x, for y,z=range x
	RealImag                  // real(x), imag(x)
	Recv                      // <-x
	Ref                       // &x, new(x)
	RefIndex                  // &x[y]
	Send                      // x->y
	SetIndex                  // x[y]=z
	Slice                     // s[x:y] for (slice, array, string)
	Slice3                    // s[x:y:z] for (slice, array)
)

func (op Op) All(other Op) bool { return op&other == other }
func (op Op) Any(other Op) bool { return op&other != None }

func (op Op) String() string {
	if op == None {
		return "none"
	}

	parts := []string{}
	add := func(other Op, name string) {
		if op.All(other) {
			parts = append(parts, name)
		}
	}

	add(Add, `Add`)
	add(Arith, `Arith`)
	add(Bitwise, `Bitwise`)
	add(ByteSlice, `ByteSlice`)
	add(Cap, `Cap`)
	add(Clear, `Clear`)
	add(Comparable, `Comparable`)
	add(Complex, `Complex`)
	add(Deref, `Deref`)
	add(GetIndex, `GetIndex`)
	add(GetIndex2, `GetIndex2`)
	add(IsNil, `IsNil`)
	add(Len, `Len`)
	add(Make, `Make`)
	add(Make3, `Make3`)
	add(Mod, `Mod`)
	add(Orderable, `Orderable`)
	add(Range, `Range`)
	add(Range2, `Range2`)
	add(RealImag, `RealImag`)
	add(Recv, `Recv`)
	add(Ref, `Ref`)
	add(RefIndex, `RefIndex`)
	add(Send, `Send`)
	add(SetIndex, `SetIndex`)
	add(Slice, `Slice`)
	add(Slice3, `Slice3`)

	return strings.Join(parts, `|`)
}
