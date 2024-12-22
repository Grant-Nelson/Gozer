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
		`	verbose|v bool = false`,
		`		Print verbose output`,
		`	input|i string (required)`,
		`		Path to input file`,
		`	output|o string = "out.txt"`,
		`		Path to write the output to`)

	parseHelp(t, newS(), `cat -i -h`,
		`Usage of cat:`,
		`Flags:`,
		`	h|help`,
		`		Shows help for the current tool`,
		`	verbose|v bool = false`,
		`		Print verbose output`,
		`	input|i string (required)`,
		`		Path to input file`,
		`	output|o string = "out.txt"`,
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
		`	v bool = true`,
		`		Print verbose output`,
		`	i string = "in.txt"`,
		`		Path to input file`,
		`	o string`,
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

	s = &S{}
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
	parsePass(t, s, `cat -Bool "true"`)
	equal(t, s.Bool, true)

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
	parsePanic(t, s1, `cat`, ErrFlagTagWrongType.with(`chan bool`))

	s2 := &struct {
		M map[string]int `arg:"flag,,"`
	}{}
	parsePanic(t, s2, `cat`, ErrFlagTagWrongType.with(`map[string]int`))

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
	parsePanic(t, s7, `cat`, ErrFlagTagWrongType.with(`[]int`))
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
		`	input|in|i string (required)`,
		`		Path to input file`,
		`	Output string (required)`,
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
	type S1 struct {
		Verbose bool   `arg:"flag,v,Print verbose output"`
		Input   string `arg:"pos,input,Path to input file"`
		Output  string `arg:"optional,pos,,Path to write the output to"`
	}

	s1 := &S1{}
	parsePass(t, s1, `cat food.json`)
	equal(t, s1.Verbose, false)
	equal(t, s1.Input, `food.json`)
	equal(t, s1.Output, ``)

	s1 = &S1{}
	parsePass(t, s1, `cat -v food.json`)
	equal(t, s1.Verbose, true)
	equal(t, s1.Input, `food.json`)
	equal(t, s1.Output, ``)

	s1 = &S1{}
	parsePass(t, s1, `cat food.json -v zoom`)
	equal(t, s1.Verbose, true)
	equal(t, s1.Input, `food.json`)
	equal(t, s1.Output, `zoom`)

	parseFail(t, &S1{}, `cat -v`,
		`Missing required positional: input`,
		`Use "cat -h" to print help.`)

	parseHelp(t, &S1{}, `cat -h`,
		`Usage of cat:`,
		`Flags:`,
		`	h|help`,
		`		Shows help for the current tool`,
		`	v bool = false`,
		`		Print verbose output`,
		`Positional Arguments:`,
		`	input string (required)`,
		`		Path to input file`,
		`	Output string = ""`,
		`		Path to write the output to`)

	type S2 struct {
		Verbose bool `arg:"flag,v,Print verbose output"`
		X       int  `arg:"pos,input,X value"`
		Y       int  `arg:"optional,pos,,Y value"`
	}
	s2 := &S2{}
	parsePass(t, s2, `cat 145 0xFF`)
	equal(t, s2.Verbose, false)
	equal(t, s2.X, 145)
	equal(t, s2.Y, 255)

	s2 = &S2{}
	parseFail(t, s2, `cat dog`,
		`Invalid integer value for input: dog`,
		`Use "cat -h" to print help.`)

	type S3 struct {
		Verbose *bool `arg:"optional,pos,,talkative"`
		X       *int  `arg:"optional,pos,,X value"`
		Y       *int  `arg:"optional,pos,,Y value"`
	}
	s3 := &S3{}
	parsePass(t, s3, `cat true -42 "+56"`)
	equal(t, s3.Verbose, toPtr(true))
	equal(t, s3.X, toPtr(-42))
	equal(t, s3.Y, toPtr(56))

	s3 = &S3{}
	parseHelp(t, s3, `cat -h`,
		`Usage of cat:`,
		`Flags:`,
		`	h|help`,
		`		Shows help for the current tool`,
		`Positional Arguments:`,
		`	Verbose bool`,
		`		talkative`,
		`	X int`,
		`		X value`,
		`	Y int`,
		`		Y value`)

	s3 = &S3{
		Verbose: toPtr(true),
		X:       toPtr(12),
		Y:       toPtr(34),
	}
	parsePass(t, s3, `cat`)
	equal(t, s3.Verbose, nil)
	equal(t, s3.X, nil)
	equal(t, s3.Y, nil)

	s3 = &S3{}
	parsePass(t, s3, `cat 1`)
	equal(t, s3.Verbose, toPtr(true))
	equal(t, s3.X, nil)
	equal(t, s3.Y, nil)

	s3 = &S3{}
	parsePass(t, s3, `cat "true"`)
	equal(t, s3.Verbose, toPtr(true))
	equal(t, s3.X, nil)
	equal(t, s3.Y, nil)

	parseFail(t, &S3{}, `cat dog`,
		`Invalid boolean value for Verbose: dog`,
		`Use "cat -h" to print help.`)

	type S4 struct {
		X int `arg:"pos,,X value"`
		Y int `arg:"required,pos,,Y value"`
	}
	s4 := &S4{}
	parsePass(t, s4, `cat 11 22`)
	equal(t, s4.X, 11)
	equal(t, s4.Y, 22)

	parseFail(t, &S4{}, `cat`,
		`Missing required positionals: X, Y`,
		`Use "cat -h" to print help.`)
}

func TestPos_Variadic(t *testing.T) {
	type S1 struct {
		Verbose bool     `arg:"flag,v,talkative"`
		Output  string   `arg:"pos,,Output file"`
		Inputs  []string `arg:"pos,,Zero or more input files"`
	}
	s1 := &S1{}
	parsePass(t, s1, `cat out.txt`)
	equal(t, s1.Output, `out.txt`)
	equal(t, s1.Inputs, []string(nil))

	s1 = &S1{}
	parsePass(t, s1, `cat out.txt a.txt b.txt`)
	equal(t, s1.Output, `out.txt`)
	equal(t, s1.Inputs, []string{`a.txt`, `b.txt`})

	s1 = &S1{
		Inputs: []string{`x`, `y`},
	}
	parsePass(t, s1, `cat out.txt a b c`)
	equal(t, s1.Output, `out.txt`)
	equal(t, s1.Inputs, []string{`a`, `b`, `c`})

	s1 = &S1{
		Inputs: []string{`x`, `y`},
	}
	parsePass(t, s1, `cat out.txt `)
	equal(t, s1.Output, `out.txt`)
	equal(t, s1.Inputs, []string(nil))

	s1 = &S1{}
	parsePass(t, s1, `cat out.txt a -v b`)
	equal(t, s1.Output, `out.txt`)
	equal(t, s1.Inputs, []string{`a`, `b`})

	s1 = &S1{}
	parsePass(t, s1, `cat out.txt a -v true b`)
	equal(t, s1.Output, `out.txt`)
	equal(t, s1.Inputs, []string{`a`, `b`})

	s1 = &S1{
		Output: `out.json`,
		Inputs: []string{`in1`, `in2`},
	}
	parseHelp(t, s1, `cat -h`,
		`Usage of cat:`,
		`Flags:`,
		`	h|help`,
		`		Shows help for the current tool`,
		`	v bool = false`,
		`		talkative`,
		`Positional Arguments:`,
		`	Output string (required)`,
		`		Output file`,
		`	Inputs []string = ["in1", "in2"]`,
		`		Zero or more input files`)

	s2 := &struct {
		Verbose bool   `arg:"flag,v,talkative"`
		Output  string `arg:"pos,,Output file"`
		Inputs  []int  `arg:"pos,,Zero or more input files"`
	}{}
	parseFail(t, s2, `cat o.json dog`,
		`Invalid integer value for Inputs: dog`,
		`Use "cat -h" to print help.`)
}

func TestPos_BadForm(t *testing.T) {
	s1 := &struct {
		Inputs []string `arg:"pos,i|input,Zero or more input files"`
	}{}
	parsePanic(t, s1, `cat`, ErrInvalidPosName.with(`"i|input"`))

	s2 := &struct {
		Inputs []string `arg:"pos,input files,Zero or more input files"`
	}{}
	parsePanic(t, s2, `cat`, ErrInvalidPosName.with(`"input files"`))

	s3 := &struct {
		Inputs []*string `arg:"pos,,Zero or more input files"`
	}{}
	parsePanic(t, s3, `cat`, ErrPosTagWrongType.with(`[]*string`))

	s4 := &struct {
		Inputs *[]string `arg:"pos,,Zero or more input files"`
	}{}
	parsePanic(t, s4, `cat`, ErrPosTagWrongType.with(`*[]string`))

	s5 := &struct {
		Inputs []string `arg:"required,pos,,"`
	}{}
	parsePanic(t, s5, `cat`, ErrVarPosRequired.with(`"Inputs"`))

	s6 := &struct {
		Input  string `arg:"pos,X,"`
		Output string `arg:"pos,X,"`
	}{}
	parsePanic(t, s6, `cat`, ErrPosAlreadyExists.with(`"X"`))

	s7 := &struct {
		Inputs []string `arg:"pos,,"`
		Output string   `arg:"pos,,"`
	}{}
	parsePanic(t, s7, `cat`, ErrVarPosNotLast.with(`"Output"`))

	s8 := &struct {
		Inputs string `arg:"optional,pos,,"`
		Output string `arg:"pos,,"`
	}{}
	parsePanic(t, s8, `cat`, ErrPosRequiredAfterOp.with(`"Output"`))
}

func toPtr[T any](v T) *T { return &v }

func splitArgs(args string) []string {
	parts := strings.Fields(args)
	for i, part := range parts {
		// For simplicity, since Fields doesn't honor quotes,
		// use underscores in quotes to represents spaces.
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
		case reflect.Array, reflect.Slice:
			if val.IsNil() {
				return `[]<nil>`
			}
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

func equal[T any](t *testing.T, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %s, want %s", formatValue(got), formatValue(want))
	}
}
