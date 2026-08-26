package typeOp

import "strings"

type Op uint64

const (
	None     Op = 0
	Make     Op = 1 << iota // make(x, len) for (slice, map, channel)
	Make3                   // make3(x, len, cap) for (slice)
	New                     // new(), new(x)
	Add                     // x+y, x+=y for (numeric, string)
	Arith                   // -x, x-y, x*y, x/y, x-=y, x*=y, x/=y, x++, x--
	Mod                     // x%y, x%=y
	Bitwise                 // x&y, x|y, ^x, x^y, x&^y, x<<y, x>>y, x&=y, x|=y, x^=y, x&^=y, x<<=y, x>>=y
	Deref                   // *x
	Ref                     // &x
	Comp                    // x==y, x!=y
	Order                   // x<y, x<=y, x>y, x>=y
	Invoke                  // x(..)
	Len                     // len(x)
	Cap                     // cap(x)
	Index                   // x[y]
	RefIndex                // &x[y]
	Slice                   // s[x:y]
	Slice3                  // s[x:y:z]
	Copy                    // copy(x, y)
	Append                  // append(x, y)
	Clear                   // clear(x)
	Combine                 // combine(x, y), complex(x, y)
	High                    // high(x), imag(x)
	Low                     // low(x), real(x)
	Send                    // x->y
	Recv                    // <-x
	Range                   // for ... range x
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

	add(Make, `make`)
	add(Make3, `make3`)
	add(New, `new`)
	add(Add, `add`)
	add(Arith, `arith`)
	add(Mod, `mod`)
	add(Bitwise, `bitwise`)
	add(Deref, `deref`)
	add(Ref, `ref`)
	add(Comp, `comp`)
	add(Order, `order`)
	add(Invoke, `invoke`)
	add(Len, `len`)
	add(Cap, `cap`)
	add(Index, `index`)
	add(RefIndex, `refIndex`)
	add(Slice, `slice`)
	add(Slice3, `slice3`)
	add(Copy, `copy`)
	add(Append, `append`)
	add(Clear, `clear`)
	add(Combine, `combine`)
	add(High, `high`)
	add(Low, `low`)
	add(Send, `send`)
	add(Recv, `recv`)
	add(Range, `range`)

	return strings.Join(parts, `|`)
}
