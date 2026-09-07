package types

// TODO: Comment
type Basic interface {
	Type
	IsBasic()
}

// TODO: Comment
type Bool struct{}

func (Bool) IsType()      {}
func (Bool) IsBasic()     {}
func (Bool) BitSize() int { return 1 }

// TODO: Comment
type Integer interface {
	Basic
	BitwiseOps
	OrderableOps
	IsInteger()
}

// TODO: Comment
type SignedInteger interface {
	Integer
	IsSigned()
}

// TODO: Comment
type UnsignedInteger interface {
	Integer
	IsUnsigned()
}

// TODO: Comment
type Int struct{}

var _ SignedInteger = Int{}

func (Int) IsType()           {}
func (Int) IsBasic()          {}
func (Int) IsInteger()        {}
func (Int) IsSigned()         {}
func (Int) HasAddOps()        {}
func (Int) HasArithmeticOps() {}
func (Int) HasBitwiseOps()    {}
func (Int) HasComparableOps() {}
func (Int) HasOrderableOps()  {}

// TODO: Comment
type Int8 struct{}

var (
	_ SignedInteger = Int8{}
	_ FixedSize     = Int8{}
)

func (Int8) IsType()           {}
func (Int8) IsBasic()          {}
func (Int8) IsInteger()        {}
func (Int8) IsSigned()         {}
func (Int8) HasAddOps()        {}
func (Int8) HasArithmeticOps() {}
func (Int8) HasBitwiseOps()    {}
func (Int8) HasComparableOps() {}
func (Int8) HasOrderableOps()  {}
func (Int8) BitSize() int      { return 8 }

// TODO: Comment
type Int16 struct{}

var (
	_ SignedInteger = Int16{}
	_ FixedSize     = Int16{}
)

func (Int16) IsType()           {}
func (Int16) IsBasic()          {}
func (Int16) IsInteger()        {}
func (Int16) IsSigned()         {}
func (Int16) HasAddOps()        {}
func (Int16) HasArithmeticOps() {}
func (Int16) HasBitwiseOps()    {}
func (Int16) HasComparableOps() {}
func (Int16) HasOrderableOps()  {}
func (Int16) BitSize() int      { return 16 }

// TODO: Comment
type Int32 struct{}

var (
	_ SignedInteger = Int32{}
	_ FixedSize     = Int32{}
)

func (Int32) IsType()           {}
func (Int32) IsBasic()          {}
func (Int32) IsInteger()        {}
func (Int32) IsSigned()         {}
func (Int32) HasAddOps()        {}
func (Int32) HasArithmeticOps() {}
func (Int32) HasBitwiseOps()    {}
func (Int32) HasComparableOps() {}
func (Int32) HasOrderableOps()  {}
func (Int32) BitSize() int      { return 32 }

// TODO: Comment
type Int64 struct{}

var (
	_ SignedInteger = Int64{}
	_ FixedSize     = Int64{}
)

func (Int64) IsType()           {}
func (Int64) IsBasic()          {}
func (Int64) IsInteger()        {}
func (Int64) IsSigned()         {}
func (Int64) HasAddOps()        {}
func (Int64) HasArithmeticOps() {}
func (Int64) HasBitwiseOps()    {}
func (Int64) HasComparableOps() {}
func (Int64) HasOrderableOps()  {}
func (Int64) BitSize() int      { return 64 }

// TODO: Comment
type Uint struct{}

var _ UnsignedInteger = Uint{}

func (Uint) IsType()           {}
func (Uint) IsBasic()          {}
func (Uint) IsInteger()        {}
func (Uint) IsUnsigned()       {}
func (Uint) HasAddOps()        {}
func (Uint) HasArithmeticOps() {}
func (Uint) HasBitwiseOps()    {}
func (Uint) HasComparableOps() {}
func (Uint) HasOrderableOps()  {}

// TODO: Comment
type Uint8 struct{}

var (
	_ UnsignedInteger = Uint8{}
	_ FixedSize       = Uint8{}
)

func (Uint8) IsType()           {}
func (Uint8) IsBasic()          {}
func (Uint8) IsInteger()        {}
func (Uint8) IsUnsigned()       {}
func (Uint8) HasAddOps()        {}
func (Uint8) HasArithmeticOps() {}
func (Uint8) HasBitwiseOps()    {}
func (Uint8) HasComparableOps() {}
func (Uint8) HasOrderableOps()  {}
func (Uint8) BitSize() int      { return 8 }

// TODO: Comment
type Uint16 struct{}

var (
	_ UnsignedInteger = Uint16{}
	_ FixedSize       = Uint16{}
)

func (Uint16) IsType()           {}
func (Uint16) IsBasic()          {}
func (Uint16) IsInteger()        {}
func (Uint16) IsUnsigned()       {}
func (Uint16) HasAddOps()        {}
func (Uint16) HasArithmeticOps() {}
func (Uint16) HasBitwiseOps()    {}
func (Uint16) HasComparableOps() {}
func (Uint16) HasOrderableOps()  {}
func (Uint16) BitSize() int      { return 16 }

// TODO: Comment
type Uint32 struct{}

var (
	_ UnsignedInteger = Uint32{}
	_ FixedSize       = Uint32{}
)

func (Uint32) IsType()           {}
func (Uint32) IsBasic()          {}
func (Uint32) IsInteger()        {}
func (Uint32) IsUnsigned()       {}
func (Uint32) HasAddOps()        {}
func (Uint32) HasArithmeticOps() {}
func (Uint32) HasBitwiseOps()    {}
func (Uint32) HasComparableOps() {}
func (Uint32) HasOrderableOps()  {}
func (Uint32) BitSize() int      { return 32 }

// TODO: Comment
type Uint64 struct{}

var (
	_ UnsignedInteger = Uint64{}
	_ FixedSize       = Uint64{}
)

func (Uint64) IsType()           {}
func (Uint64) IsBasic()          {}
func (Uint64) IsInteger()        {}
func (Uint64) IsUnsigned()       {}
func (Uint64) HasAddOps()        {}
func (Uint64) HasArithmeticOps() {}
func (Uint64) HasBitwiseOps()    {}
func (Uint64) HasComparableOps() {}
func (Uint64) HasOrderableOps()  {}
func (Uint64) BitSize() int      { return 64 }

// TODO: Comment
type UintPtr struct{}

var _ UnsignedInteger = Uint{}

func (UintPtr) IsType()           {}
func (UintPtr) IsBasic()          {}
func (UintPtr) IsInteger()        {}
func (UintPtr) IsUnsigned()       {}
func (UintPtr) HasAddOps()        {}
func (UintPtr) HasArithmeticOps() {}
func (UintPtr) HasBitwiseOps()    {}
func (UintPtr) HasComparableOps() {}
func (UintPtr) HasOrderableOps()  {}

// TODO: Comment
type String struct{}

// TODO: Comment
type Float interface {
	Basic
	ArithmeticOps
	OrderableOps
	IsFloat()
}

// TODO: Comment
type Float32 struct{}

var (
	_ Float      = Float32{}
	_ FixedSize  = Float32{}
	_ ComplexOps = Float32{}
)

func (Float32) IsType()              {}
func (Float32) IsBasic()             {}
func (Float32) IsFloat()             {}
func (Float32) HasAddOps()           {}
func (Float32) HasArithmeticOps()    {}
func (Float32) HasComparableOps()    {}
func (Float32) HasOrderableOps()     {}
func (Float32) HasComplexOps()       {}
func (Float32) ComplexType() Complex { return Complex64{} }
func (Float32) BitSize() int         { return 32 }

// TODO: Comment
type Float64 struct{}

var (
	_ Float      = Float64{}
	_ FixedSize  = Float64{}
	_ ComplexOps = Float64{}
)

func (Float64) IsType()              {}
func (Float64) IsBasic()             {}
func (Float64) IsFloat()             {}
func (Float64) HasAddOps()           {}
func (Float64) HasArithmeticOps()    {}
func (Float64) HasComparableOps()    {}
func (Float64) HasOrderableOps()     {}
func (Float64) HasComplexOps()       {}
func (Float64) ComplexType() Complex { return Complex128{} }
func (Float64) BitSize() int         { return 64 }

// TODO: Comment
type Numeric struct{}

var _ Float = Numeric{}

func (Numeric) IsType()           {}
func (Numeric) IsBasic()          {}
func (Numeric) IsFloat()          {}
func (Numeric) HasAddOps()        {}
func (Numeric) HasArithmeticOps() {}
func (Numeric) HasComparableOps() {}
func (Numeric) HasOrderableOps()  {}

// TODO: Comment
type Complex interface {
}

// TODO: Comment
type Complex64 struct{}

// TODO: Comment
type Complex128 struct{}
