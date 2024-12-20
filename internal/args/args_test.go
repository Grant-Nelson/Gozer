package args

import (
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
	equal(t, s.Verbose, toPtr(false))
	equal(t, s.Input, toPtr(`in.txt`))
	equal(t, s.Output, nil)

	// TODO: Finish
}

func TestFlags_Types(t *testing.T) {
	// TODO: Fill out
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

func toPtr[T any](v T) *T {
	return &v
}

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

func isEqual[T comparable](t *testing.T, got, want T) bool {
	t.Helper()
	if reflect.TypeFor[T]().Kind() == reflect.Pointer {
		vg := reflect.ValueOf(got)
		vw := reflect.ValueOf(want)
		if vg.IsNil() {
			if !vw.IsNil() {

				return
			}
			return
		}

	}
	return got != want
}

func formatValue[T any](t *testing.T, value T) string {

}

func equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if reflect.TypeFor[T]().Kind() == reflect.Pointer {
		vg := reflect.ValueOf(got)
		vw := reflect.ValueOf(want)
		if vg.IsNil() {
			if !vw.IsNil() {

				return
			}
			return
		}

	}
	if got != want {
		if reflect.TypeFor[T]().Kind() == reflect.String {
			t.Errorf("got %q, want %q", any(got).(string), any(want).(string))
		} else {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}
