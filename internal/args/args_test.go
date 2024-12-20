package args

import (
	"reflect"
	"strings"
	"testing"
)

func TestEmpty_Ok(t *testing.T) {
	s := &struct{}{}
	parsePass(t, s, `cat`)

	parsePanic(t, s, ``, ErrTooFewArgs)
	parseFail(t, s, `cat pants`,
		`Unexpected positional argument: pants`,
		`Use "cat -h" to print help.`)
	parseFail(t, s, `cat -pants`,
		`Unknown flag: pants`,
		`Use "cat -h" to print help.`)

	parsePanic(t, *s, `cat`, ErrNotStructPointer.with(`got struct`))
	var nilS *struct{}
	parsePanic(t, nilS, `cat`, ErrNilPointer)
	var intS int
	parsePanic(t, intS, `cat`, ErrNotStructPointer.with(`got int`))
	parsePanic(t, &intS, `cat`, ErrNotStructPointer.with(`got pointer to int`))
}

func TestForm_Flags(t *testing.T) {
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

	parseFail(t, newS(), `cat`,
		`Missing required flag: input|i`,
		`Use "cat -h" to print help.`)

}

func parsePass(t *testing.T, s any, args string) {
	t.Helper()
	buf := &strings.Builder{}
	result := ParseArgs(s, strings.Fields(args), buf)
	equal(t, result, true)
	equal(t, buf.String(), ``)
}

func parseFail(t *testing.T, s any, args string, lines ...string) {
	t.Helper()
	buf := &strings.Builder{}
	result := ParseArgs(s, strings.Fields(args), buf)
	equal(t, result, false)
	equal(t, buf.String(), strings.Join(lines, "\n")+"\n")
}

func parsePanic(t *testing.T, s any, args string, exp error) {
	t.Helper()
	defer func() {
		t.Helper()
		if r := recover(); r != nil {
			equal(t, r.(error).Error(), exp.Error())
		}
	}()
	ParseArgs(s, strings.Fields(args), nil)
	t.Fatal(`expected panic`)
}

func equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		if reflect.TypeFor[T]().Kind() == reflect.String {
			t.Errorf("got %q, want %q", any(got).(string), any(want).(string))
		} else {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}
