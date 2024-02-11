package tools

type ArgReader interface {
	Flag(target *bool, short, long string) ArgReader

	NamedBool(target *bool, short, long string) ArgReader
	NamedInt(target *int, short, long string) ArgReader
	NamedStr(target *string, short, long string) ArgReader

	PosBool(target *bool) ArgReader
	PosInt(target *int) ArgReader
	PosStr(target *string) ArgReader

	OptionalBool(target *bool) ArgReader
	OptionalInt(target *int) ArgReader
	OptionalStr(target *string) ArgReader

	VarBool(target *[]bool) ArgReader
	VarInt(target *[]int) ArgReader
	VarStr(target *[]string) ArgReader

	Process(args []string) error
}
