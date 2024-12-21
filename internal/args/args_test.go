package args

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestEmpty_Args(t *testing.T) {
	s := &struct{}{}
	parsePass(t, s, `cat`)

	parsePanic(t, s, ``, ErrTooFewArgs)
	parseFail(t, s, `cat pants`,
		`Unexpected positional argument: pants`,
		`Use "cat -h" to print help.`)
	parseFail(t, s, `cat -pants`,
		`Unknown flag "pants".`,
		`Use "cat -h" to print help.`)

	parseHelp(t, s, `cat -h`,
		`Usage of cat:`,
		`Flags:`,
		`	h|help`,
		`		Shows help for the current tool`)
}

func TestEmpty_Form(t *testing.T) {
	var s struct{}
	parsePanic(t, s, `cat`, ErrNotStructPointer.with(`got struct`))
	var nilS *struct{}
	parsePanic(t, nilS, `cat`, ErrNilPointer)
	var intS int
	parsePanic(t, intS, `cat`, ErrNotStructPointer.with(`got int`))
	parsePanic(t, &intS, `cat`, ErrNotStructPointer.with(`got pointer to int`))
}

func TestEmpty_Skips(t *testing.T) {
	s := &struct {
		unexported bool
		Skipped    bool `arg:"skip"`
		NoTag      bool
	}{}
	parsePass(t, s, `cat`)
	parseHelp(t, s, `cat -h`,
		`Usage of cat:`,
		`Flags:`,
		`	h|help`,
		`		Shows help for the current tool`)
}

func TestFlags_Args(t *testing.T) {
	type S struct {
		Verbose bool   `arg:"flag,verbose|v,Print verbose output"`
		Input   string `arg:"required,flag,input|i,Path to input file"`
		Output  string `arg:"optional,flag,output|o,Path to write the output to"`
	}
	newS := func() *S {
		return &S{
			Verbose: false,
			Input:   `in.txt`,
			Output:  `out.txt`,
		}
	}

	s := newS()
	parsePass(t, s, `cat -v -i pet.txt -o purr.txt`)
	equal(t, s.Verbose, true)
	equal(t, s.Input, `pet.txt`)
	equal(t, s.Output, `purr.txt`)

	s = newS()
	parsePass(t, s, `cat -v true -i pet.txt`)
	equal(t, s.Verbose, true)
	equal(t, s.Input, `pet.txt`)
	equal(t, s.Output, `out.txt`)

	s = newS()
	parsePass(t, s, `cat -o purr.txt -i pet.txt`)
	equal(t, s.Verbose, false)
	equal(t, s.Input, `pet.txt`)
	equal(t, s.Output, `purr.txt`)

	s = newS()
	parsePass(t, s, `cat -o purr.txt -i pet.txt -v`)
	equal(t, s.Verbose, true)
	equal(t, s.Input, `pet.txt`)
	equal(t, s.Output, `purr.txt`)

	s = newS()
	parsePass(t, s, `cat -o purr.txt -i "-x"`)
	equal(t, s.Verbose, false)
	equal(t, s.Input, `-x`)
	equal(t, s.Output, `purr.txt`)

	s = newS()
	parsePass(t, s, `cat -v false -i "pet_cat" -o "bite_\"em"`)
	equal(t, s.Verbose, false)
	equal(t, s.Input, `pet cat`)
	equal(t, s.Output, `bite "em`)

	parseFail(t, newS(), `cat`,
		`Missing required flag: input|i`,
		`Use "cat -h" to print help.`)
	parseFail(t, newS(), `cat -k`,
		`Unknown flag "k".`,
		`Use "cat -h" to print help.`)
	parseFail(t, newS(), `cat -i`,
		`"i" flag requires a value.`,
		`Use "cat -h" to print help.`)
	parseFail(t, newS(), `cat -i purr.txt -i pet.txt`,
		`"i" flag already set.`,
		`Use "cat -h" to print help.`)
	parseFail(t, newS(), `cat -i -x`,
		`"i" flag requires a value.`,
		`If the intended string value starts with a dash, escape the value: -i "-x"`,
		`Use "cat -h" to print help.`)

	parseHelp(t, newS(), `cat -h`,
		`Usage of cat:`,
		`Flags:`,
		`	h|help`,
		`		Shows help for the current tool`,
		`	verbose|v = false`,
		`		Print verbose output`,
		`	input|i (required)`,
		`		Path to input file`,
		`	output|o = "out.txt"`,
		`		Path to write the output to`)

	parseHelp(t, newS(), `cat -i -h`,
		`Usage of cat:`,
		`Flags:`,
		`	h|help`,
		`		Shows help for the current tool`,
		`	verbose|v = false`,
		`		Print verbose output`,
		`	input|i (required)`,
		`		Path to input file`,
		`	output|o = "out.txt"`,
		`		Path to write the output to`)
}

func TestFlags_Clearable(t *testing.T) {
	type S struct {
		Verbose *bool   `arg:"flag,v,Print verbose output"`
		Input   *string `arg:"flag,i,Path to input file"`
		Output  *string `arg:"flag,o,Path to write the output to"`
	}
	newS := func() *S {
		return &S{
			Verbose: toPtr(true),
			Input:   toPtr(`in.txt`),
			Output:  nil,
		}
	}

	s := newS()
	parseHelp(t, s, `cat -h`,
		`Usage of cat:`,
		`Flags:`,
		`	h|help`,
		`		Shows help for the current tool`,
		`	v = true`,
		`		Print verbose output`,
		`	i = "in.txt"`,
		`		Path to input file`,
		`	o`,
		`		Path to write the output to`)
	equal(t, s.Verbose, toPtr(true))
	equal(t, s.Input, toPtr(`in.txt`))
	equal(t, s.Output, nil)

	s = newS()
	parsePass(t, s, `cat`)
	equal(t, s.Verbose, nil)
	equal(t, s.Input, nil)
	equal(t, s.Output, nil)

	s = newS()
	parsePass(t, s, `cat -v false -i in.json -o out.json`)
	equal(t, s.Verbose, toPtr(false))
	equal(t, s.Input, toPtr(`in.json`))
	equal(t, s.Output, toPtr(`out.json`))

	s = newS()
	parsePass(t, s, `cat -v -o out.json`)
	equal(t, s.Verbose, toPtr(true))
	equal(t, s.Input, nil)
	equal(t, s.Output, toPtr(`out.json`))

	s = newS()
	parsePass(t, s, `cat -v`)
	equal(t, s.Verbose, toPtr(true))
	equal(t, s.Input, nil)
	equal(t, s.Output, nil)
}

func TestFlags_Types(t *testing.T) {
	s := &struct {
		Bool    bool    `arg:"flag,,"`
		Int     int     `arg:"flag,,"`
		Int8    int8    `arg:"flag,,"`
		Int16   int16   `arg:"flag,,"`
		Int32   int32   `arg:"flag,,"`
		Int64   int64   `arg:"flag,,"`
		Uint    uint    `arg:"flag,,"`
		Uint8   uint8   `arg:"flag,,"`
		Uint16  uint16  `arg:"flag,,"`
		Uint32  uint32  `arg:"flag,,"`
		Uint64  uint64  `arg:"flag,,"`
		Float32 float32 `arg:"flag,,"`
		Float64 float64 `arg:"flag,,"`
		String  string  `arg:"flag,,"`
		Byte    byte    `arg:"flag,,"`
		Rune    rune    `arg:"flag,,"`
	}{}

	parsePass(t, s, `cat -Bool true`)
	equal(t, s.Bool, true)
	parsePass(t, s, `cat -Bool F`)
	equal(t, s.Bool, false)
	parsePass(t, s, `cat -Bool 1`)
	equal(t, s.Bool, true)
	parsePass(t, s, `cat -Bool FALSE`)
	equal(t, s.Bool, false)
	parseFail(t, s, `cat -Bool dog`,
		`Unexpected positional argument: dog`,
		`Use "cat -h" to print help.`)

	parsePass(t, s, `cat -Int 123`)
	equal(t, s.Int, 123)
	parsePass(t, s, `cat -Int -321`)
	equal(t, s.Int, -321)

	parseFail(t, s, `cat -Int8 1423`,
		`Integer value for Int8 is out of the range -128 to 127: 1423`,
		`Use "cat -h" to print help.`)
	parseFail(t, s, `cat -Int8 dog`,
		`Invalid integer value for Int8: dog`,
		`Use "cat -h" to print help.`)
	parsePass(t, s, `cat -Int8 42`)
	equal(t, s.Int8, 42)
	parsePass(t, s, `cat -Int8 +42`)
	equal(t, s.Int8, 42)
	parsePass(t, s, `cat -Int8 -0x1F`)
	equal(t, s.Int8, -31)
	parsePass(t, s, `cat -Int8 0b0101`)
	equal(t, s.Int8, 5)

	parsePass(t, s, `cat -Int16 2468`)
	equal(t, s.Int16, 2468)
	parsePass(t, s, `cat -Int32 554`)
	equal(t, s.Int32, 554)
	parsePass(t, s, `cat -Int64 773`)
	equal(t, s.Int64, 773)

	parsePass(t, s, `cat -Uint 123`)
	equal(t, s.Uint, 123)

	parseFail(t, s, `cat -Uint8 1423`,
		`Unsigned integer value for Uint8 is out of the range 0 to 255: 1423`,
		`Use "cat -h" to print help.`)
	parseFail(t, s, `cat -Uint8 -5`,
		`Unsigned integer value for Uint8 is out of the range 0 to 255: -5`,
		`Use "cat -h" to print help.`)
	parseFail(t, s, `cat -Uint8 dog`,
		`Invalid unsigned integer value for Uint8: dog`,
		`Use "cat -h" to print help.`)
	parsePass(t, s, `cat -Uint8 42`)
	equal(t, s.Uint8, 42)
	parsePass(t, s, `cat -Uint8 +42`)
	equal(t, s.Uint8, 42)
	parsePass(t, s, `cat -Uint8 0x1F`)
	equal(t, s.Uint8, 31)
	parsePass(t, s, `cat -Uint8 0b0101`)
	equal(t, s.Uint8, 5)

	parsePass(t, s, `cat -Uint16 2468`)
	equal(t, s.Uint16, 2468)
	parsePass(t, s, `cat -Uint32 554`)
	equal(t, s.Uint32, 554)
	parsePass(t, s, `cat -Uint64 773`)
	equal(t, s.Uint64, 773)

	parsePass(t, s, `cat -Float32 1.0e-2`)
	equal(t, s.Float32, 1.0e-2)
	parsePass(t, s, `cat -Float32 -0.024`)
	equal(t, s.Float32, -0.024)

	parseFail(t, s, `cat -Float64 dog`,
		`Invalid float value for Float64: dog`,
		`Use "cat -h" to print help.`)
	parseFail(t, s, `cat -Float64 1.0.0`,
		`Invalid float value for Float64: 1.0.0`,
		`Use "cat -h" to print help.`)
	parseFail(t, s, `cat -Float64 1.0e1000000`,
		`Unsigned integer value for Float64 is out of range for a float 64: 1.0e1000000`,
		`Use "cat -h" to print help.`)
	parsePass(t, s, `cat -Float64 -1e+12`)
	equal(t, s.Float64, -1e12)

	parsePass(t, s, `cat -String dog`)
	equal(t, s.String, `dog`)
	parsePass(t, s, `cat -String 1.0e2`)
	equal(t, s.String, `1.0e2`)
	parsePass(t, s, `cat -String ""`)
	equal(t, s.String, ``)
	parsePass(t, s, `cat -String "dog_and_cat"`)
	equal(t, s.String, `dog and cat`)

	parsePass(t, s, `cat -Byte 142`)
	equal(t, s.Byte, 142)
	parsePass(t, s, `cat -Rune 242`)
	equal(t, s.Rune, 242)
}

func TestFlags_BadForm(t *testing.T) {
	s1 := &struct {
		C chan bool `arg:"flag,,"`
	}{}
	parsePanic(t, s1, `cat`, ErrFlagTagWrongType.with(`chan`))

	s2 := &struct {
		M map[string]int `arg:"flag,,"`
	}{}
	parsePanic(t, s2, `cat`, ErrFlagTagWrongType.with(`map`))

	s3 := &struct {
		X int `arg:"flag,bad kitty,"`
	}{}
	parsePanic(t, s3, `cat`, ErrInvalidFlagName.with(`"bad kitty"`))

	s4 := &struct {
		X int `arg:"flag,-X,"`
	}{}
	parsePanic(t, s4, `cat`, ErrInvalidFlagName.with(`"-X"`))

	s5 := &struct {
		X int `arg:"flag,4cats,"`
	}{}
	parsePanic(t, s5, `cat`, ErrInvalidFlagName.with(`"4cats"`))

	s6 := &struct {
		X int `arg:"flag,x,"`
		Y int `arg:"flag,x,"`
	}{}
	parsePanic(t, s6, `cat`, ErrFlagAlreadyExists.with(`"x"`))

	s7 := &struct {
		X []int `arg:"flag,,"`
	}{}
	parsePanic(t, s7, `cat`, ErrFlagTagWrongType.with(`slice`))
}

func TestFlags_Form(t *testing.T) {
	s1 := &struct {
		Input  string `arg:"required,flag,input|in|i,Path to input file"`
		Output string `arg:"required,flag,,Path to write the output to"`
	}{}
	parseFail(t, s1, `cat`,
		`Missing required flags: input|in|i, Output`,
		`Use "cat -h" to print help.`)
	parseHelp(t, s1, `cat -h`,
		`Usage of cat:`,
		`Flags:`,
		`	h|help`,
		`		Shows help for the current tool`,
		`	input|in|i (required)`,
		`		Path to input file`,
		`	Output (required)`,
		`		Path to write the output to`)

	s2 := &struct {
		Input  string `arg:"required,flag,f|file,Path to input file"`
		Output string `arg:"required,flag,f|file,Path to write the output to"`
	}{}
	parsePanic(t, s2, `cat`, ErrFlagAlreadyExists.with(`%q`, `f`))

	s3 := &struct {
		Input string `arg:"required,flag,i|h|input,Path to input file"`
	}{}
	parsePanic(t, s3, `cat`, ErrFlagNameReserved.with(`%q`, `h`))

	s4 := &struct {
		Input string `arg:"required,flag,i|help|input,Path to input file"`
	}{}
	parsePanic(t, s4, `cat`, ErrFlagNameReserved.with(`%q`, `help`))
}

func TestPos_Args(t *testing.T) {
	// TODO: Implement
}

func toPtr[T any](v T) *T { return &v }

func splitArgs(args string) []string {
	parts := strings.Fields(args)
	for i, part := range parts {
		parts[i] = strings.Replace(part, `_`, ` `, -1)
	}
	return parts
}

func parsePass(t *testing.T, s any, args string) {
	t.Helper()
	bufOut := &strings.Builder{}
	bufErr := &strings.Builder{}
	result := ParseArgs(s, splitArgs(args), bufOut, bufErr)
	equal(t, result, true)
	equal(t, bufOut.String(), ``)
	equal(t, bufErr.String(), ``)
}

func parseFail(t *testing.T, s any, args string, lines ...string) {
	t.Helper()
	bufOut := &strings.Builder{}
	bufErr := &strings.Builder{}
	result := ParseArgs(s, splitArgs(args), bufOut, bufErr)
	equal(t, result, false)
	equal(t, bufOut.String(), ``)
	equal(t, bufErr.String(), strings.Join(lines, "\n")+"\n")
}

func parseHelp(t *testing.T, s any, args string, lines ...string) {
	t.Helper()
	bufOut := &strings.Builder{}
	bufErr := &strings.Builder{}
	result := ParseArgs(s, splitArgs(args), bufOut, bufErr)
	equal(t, result, false)
	equal(t, bufOut.String(), strings.Join(lines, "\n")+"\n")
	equal(t, bufErr.String(), ``)
}

func parsePanic(t *testing.T, s any, args string, exp error) {
	t.Helper()
	bufOut := &strings.Builder{}
	bufErr := &strings.Builder{}
	defer func() {
		t.Helper()
		if r := recover(); r != nil {
			equal(t, r.(error).Error(), exp.Error())
			equal(t, bufOut.String(), ``)
			equal(t, bufErr.String(), ``)
		}
	}()
	ParseArgs(s, strings.Fields(args), bufOut, bufErr)
	t.Fatal(`expected panic`)
}

func formatValue(value any) string {
	var format func(val reflect.Value) string
	format = func(val reflect.Value) string {
		switch val.Type().Kind() {
		case reflect.Pointer:
			return `*` + format(val.Elem())
		case reflect.String:
			return fmt.Sprintf(`%q`, val.String())
		case reflect.Array:
			elems := make([]string, val.Len())
			for i := range val.Len() {
				elems[i] = format(val.Index(i))
			}
			return `[` + strings.Join(elems, `, `) + `]`
		default:
			return fmt.Sprintf(`%v`, val)
		}
	}
	return format(reflect.ValueOf(value))
}

func equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %s, want %s", formatValue(got), formatValue(want))
	}
}
