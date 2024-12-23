package args

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func parseArgs(f *form, val reflect.Value, args []string, stdOut, stdErr io.Writer) bool {
	if len(args) < 1 {
		panic(ErrTooFewArgs)
	}
	if stdOut == nil {
		stdOut = os.Stdout
	}
	if stdErr == nil {
		stdErr = os.Stderr
	}
	p := &parser{
		cmdPath:    args[0],
		form:       f,
		foundFlags: map[*flagForm]bool{},
		usedTool:   map[*toolForm]bool{},
		atPos:      0,
		val:        val,
		args:       args[1:],
		firstVar:   true,
		stdOut:     stdOut,
		stdErr:     stdErr,
	}
	return p.Parse()
}

type parser struct {
	cmdPath    string
	form       *form
	foundFlags map[*flagForm]bool
	usedTool   map[*toolForm]bool
	atPos      int
	val        reflect.Value
	args       []string
	firstVar   bool
	stdOut     io.Writer
	stdErr     io.Writer
}

func (p *parser) printf(format string, a ...any) {
	if _, err := fmt.Fprintf(p.stdOut, format+"\n", a...); err != nil {
		panic(ErrWriteOutFailure.with(`%w`, err))
	}
}

func (p *parser) errorf(format string, a ...any) {
	if _, err := fmt.Fprintf(p.stdErr, format+"\n", a...); err != nil {
		panic(ErrWriteErrFailure.with(`%w`, err))
	}
}

func (p *parser) takeArg() string {
	arg := p.args[0]
	p.args = p.args[1:]
	return arg
}

func (p *parser) isFlag(arg string) bool {
	return len(arg) > 1 &&
		strings.HasPrefix(arg, dash) &&
		!unicode.IsDigit(rune(arg[1]))
}

func (p *parser) printHelpHint() {
	p.errorf(`Use %q to print help.`, p.cmdPath+` -h`)
}

func (p *parser) printHelp() {
	p.printf(`Usage of %s:`, p.cmdPath)
	for _, helpField := range p.form.Help {
		helpText := p.val.FieldByIndex(helpField.Index).String()
		if len(helpText) > 0 {
			p.printf(helpText)
		}
	}
	p.printHelpForTools()
	p.printHelpForFlags()
	p.printHelpForPos()
}

func (p *parser) printHelpForTools() {
	if len(p.form.Tools) == 0 {
		return
	}

	p.printf(`Tools:`)
	for _, tool := range p.form.AllTools {
		p.printf(indent+`%s`, strings.Join(tool.Names, nameSep))
		if len(tool.Description) > 0 {
			p.printf(indent+indent+`%s`, tool.Description)
		}
	}
}

func (p *parser) defaultValue(val reflect.Value) string {
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return ``
		}
		val = val.Elem()
	}

	if val.Kind() == reflect.Slice {
		parts := make([]string, val.Len())
		for i := range val.Len() {
			parts[i] = p.defaultValue(val.Index(i))
		}
		return `[` + strings.Join(parts, `, `) + `]`
	}

	if val.Kind() == reflect.String {
		return fmt.Sprintf(`%q`, val.String())
	}
	return fmt.Sprintf(`%v`, val.Interface())
}

func (p *parser) printHelpForFlags() {
	p.printf(`Flags:`)
	p.printf(indent + helpShort + nameSep + help)
	p.printf(indent + indent + `Shows help for the current tool`)
	for _, flag := range p.form.AllFlags {
		extra := ``
		if flag.Required {
			extra = ` (` + required + `)`
		} else {
			val := p.val.FieldByIndex(flag.Field.Index)
			if defVal := p.defaultValue(val); len(defVal) > 0 {
				extra = ` = ` + defVal
			}
		}
		t := flag.Field.Type
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		p.printf(indent+`%s %s%s`, strings.Join(flag.Names, nameSep), t.String(), extra)
		if len(flag.Description) > 0 {
			p.printf(indent+indent+`%s`, flag.Description)
		}
	}
}

func (p *parser) printHelpForPos() {
	if len(p.form.Pos) == 0 {
		return
	}

	p.printf(`Positional Arguments:`)
	for i, pos := range p.form.Pos {
		extra := ``
		if i < p.form.PosOpAt {
			extra = ` (` + required + `)`
		} else {
			val := p.val.FieldByIndex(pos.Field.Index)
			if defVal := p.defaultValue(val); len(defVal) > 0 {
				extra = ` = ` + defVal
			}
		}
		t := pos.Field.Type
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		p.printf(indent+`%s %s%s`, pos.Name, t.String(), extra)
		if len(pos.Description) > 0 {
			p.printf(indent+indent+`%s`, pos.Description)
		}
	}
}

func (p *parser) satisfiedFlags() bool {
	missingFlags := []string{}
	for _, flag := range p.form.AllFlags {
		if flag.Required && !p.foundFlags[flag] {
			missingFlags = append(missingFlags, strings.Join(flag.Names, nameSep))
		}
	}
	switch len(missingFlags) {
	case 0:
		return true
	case 1:
		p.errorf(`Missing required flag: %s`, missingFlags[0])
		return false
	default:
		p.errorf(`Missing required flags: %s`, strings.Join(missingFlags, `, `))
		return false
	}
}

func (p *parser) satisfiedPos() bool {
	if p.atPos >= p.form.PosOpAt {
		return true
	}
	missingPos := []string{}
	for _, pos := range p.form.Pos[p.atPos:p.form.PosOpAt] {
		missingPos = append(missingPos, pos.Name)
	}
	switch len(missingPos) {
	case 1:
		p.errorf(`Missing required positional: %s`, missingPos[0])
		return false
	default:
		p.errorf(`Missing required positionals: %s`, strings.Join(missingPos, `, `))
		return false
	}
}

func (p *parser) getField(field reflect.StructField) reflect.Value {
	return p.val.FieldByIndex(field.Index)
}

func (p *parser) clearUnset() {
	for _, flag := range p.form.AllFlags {
		if !p.foundFlags[flag] {
			p.clearField(flag.Field)
		}
	}
	for _, tool := range p.form.AllTools {
		if !p.usedTool[tool] {
			p.clearField(tool.Field)
		}
	}
	if posCount := len(p.form.Pos); p.atPos < posCount &&
		(!p.form.Variadic || p.atPos != posCount-1 || p.firstVar) {
		for _, pos := range p.form.Pos[p.atPos:] {
			p.clearField(pos.Field)
		}
	}
}

func (p *parser) clearField(field reflect.StructField) {
	val := p.getField(field)
	switch val.Kind() {
	case reflect.Ptr, reflect.Slice:
		val.Set(reflect.Zero(val.Type()))
	}
}

func (p *parser) Parse() bool {
	for len(p.args) > 0 {
		if !p.parseArg() {
			return false
		}
	}
	if p.satisfiedFlags() && p.satisfiedPos() {
		p.clearUnset()
		return true
	}
	p.printHelpHint()
	return false
}

func (p *parser) parseArg() bool {
	arg := p.takeArg()
	if p.isFlag(arg) {
		flagName, _ := strings.CutPrefix(arg, dash)
		return p.parseFlag(flagName)
	}
	if tool, ok := p.form.Tools[arg]; ok {
		return p.parseTool(arg, tool)
	}
	return p.parsePos(arg)
}

func (p *parser) parseTool(name string, tool *toolForm) bool {
	p.usedTool[tool] = true

	if !p.satisfiedFlags() || !p.satisfiedPos() {
		p.errorf(`Must fill requirements prior to calling tool %s.`, name)
		p.printHelpHint()
		return false
	}
	p.clearUnset()

	val := p.getField(tool.Field)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}

	p.cmdPath += ` ` + name
	p.form = tool.Form
	p.foundFlags = map[*flagForm]bool{}
	p.usedTool = map[*toolForm]bool{}
	p.atPos = 0
	p.val = val
	return true
}

func (p *parser) parseFlag(name string) bool {
	if name == help || name == helpShort {
		p.printHelp()
		return false
	}

	flag, ok := p.form.Flags[name]
	if !ok {
		p.errorf(`Unknown flag %q.`, name)
		p.printHelpHint()
		return false
	}
	if p.foundFlags[flag] {
		p.errorf(`%q flag already set.`, name)
		p.printHelpHint()
		return false
	}
	p.foundFlags[flag] = true

	if p.isFieldBool(flag.Field.Type) {
		val := p.getField(flag.Field)
		if len(p.args) > 0 {
			if b, ok := p.isBool(p.args[0]); ok {
				p.takeArg()
				p.setBool(val, b)
				return true
			}
		}
		p.setBool(val, true)
		return true
	}

	if len(p.args) == 0 {
		p.errorf(`%q flag requires a value.`, name)
		p.printHelpHint()
		return false
	}

	value := p.takeArg()
	if p.isFlag(value) {
		flagName, _ := strings.CutPrefix(value, dash)
		if flagName == help || flagName == helpShort {
			p.printHelp()
			return false
		}

		p.errorf(`%q flag requires a value.`, name)
		if flag.Field.Type.Kind() == reflect.String {
			p.errorf(`If the intended string value starts with a dash, escape the value: -%s %q`, name, value)
		}
		p.printHelpHint()
		return false
	}

	val := p.getField(flag.Field)
	return p.setValue(val, name, value)
}

func (p *parser) parsePos(value string) bool {
	posCount := len(p.form.Pos)
	if p.atPos >= posCount {
		p.errorf(`Unexpected positional argument: %s`, value)
		p.printHelpHint()
		return false
	}

	pos := p.form.Pos[p.atPos]
	val := p.getField(pos.Field)
	if p.form.Variadic && p.atPos == posCount-1 {
		elem := reflect.New(val.Type().Elem())
		if !p.setValue(elem, pos.Name, value) {
			return false
		}
		if p.firstVar {
			p.firstVar = false
			val.Set(reflect.New(val.Type()).Elem())
		}
		val.Set(reflect.Append(val, elem.Elem()))
		return true
	}

	p.atPos++
	return p.setValue(val, pos.Name, value)
}

func (p *parser) isFieldBool(t reflect.Type) bool {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Kind() == reflect.Bool
}

func (p *parser) isBool(value string) (bool, bool) {
	if unquoted, err := strconv.Unquote(value); err == nil {
		value = unquoted
	}
	b, err := strconv.ParseBool(value)
	return b, err == nil
}

func (p *parser) setBool(val reflect.Value, value bool) {
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}
	val.SetBool(value)
}

func (p *parser) setValue(val reflect.Value, name string, value string) bool {
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}

	if unquoted, err := strconv.Unquote(value); err == nil {
		value = unquoted
	}

	switch val.Kind() {
	case reflect.Bool:
		return p.setBoolValue(val, name, value)
	case reflect.Int:
		return p.setIntValue(val, name, value, 0, math.MinInt, math.MaxInt)
	case reflect.Int8:
		return p.setIntValue(val, name, value, 8, math.MinInt8, math.MaxInt8)
	case reflect.Int16:
		return p.setIntValue(val, name, value, 16, math.MinInt16, math.MaxInt16)
	case reflect.Int32:
		return p.setIntValue(val, name, value, 32, math.MinInt32, math.MaxInt32)
	case reflect.Int64:
		return p.setIntValue(val, name, value, 64, math.MinInt64, math.MaxInt64)
	case reflect.Uint:
		return p.setUintValue(val, name, value, 0, math.MaxUint)
	case reflect.Uint8:
		return p.setUintValue(val, name, value, 8, math.MaxUint8)
	case reflect.Uint16:
		return p.setUintValue(val, name, value, 16, math.MaxUint16)
	case reflect.Uint32:
		return p.setUintValue(val, name, value, 32, math.MaxUint32)
	case reflect.Uint64:
		return p.setUintValue(val, name, value, 64, math.MaxUint64)
	case reflect.Float32:
		return p.setFloatValue(val, name, value, 32)
	case reflect.Float64:
		return p.setFloatValue(val, name, value, 64)
	case reflect.String:
		val.SetString(value)
		return true
	}
	p.errorf(`Argument parser error: Unsupported type: %v`, val.Kind())
	return false
}

func (p *parser) setBoolValue(val reflect.Value, name string, value string) bool {
	b, err := strconv.ParseBool(value)
	if err != nil {
		p.errorf(`Invalid boolean value for %s: %s`, name, value)
		p.printHelpHint()
		return false
	}
	val.SetBool(b)
	return true
}

func (p *parser) setIntValue(val reflect.Value, name string, value string, bitSize int, min, max int64) bool {
	n, err := strconv.ParseInt(value, 0, bitSize)
	if err != nil {
		if errors.Is(err, strconv.ErrRange) {
			p.errorf(`Integer value for %s is out of the range %d to %d: %s`, name, min, max, value)
			p.printHelpHint()
			return false
		}
		p.errorf(`Invalid integer value for %s: %s`, name, value)
		p.printHelpHint()
		return false
	}
	val.SetInt(n)
	return true
}

func (p *parser) setUintValue(val reflect.Value, name string, value string, bitSize int, max uint64) bool {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, `-`) {
		p.errorf(`Unsigned integer value for %s is out of the range 0 to %d: %s`, name, max, value)
		p.printHelpHint()
		return false
	}
	value, _ = strings.CutPrefix(value, `+`)
	n, err := strconv.ParseUint(value, 0, bitSize)
	if err != nil {
		if errors.Is(err, strconv.ErrRange) {
			p.errorf(`Unsigned integer value for %s is out of the range 0 to %d: %s`, name, max, value)
			p.printHelpHint()
			return false
		}
		p.errorf(`Invalid unsigned integer value for %s: %s`, name, value)
		p.printHelpHint()
		return false
	}
	val.SetUint(n)
	return true
}

func (p *parser) setFloatValue(val reflect.Value, name string, value string, bitSize int) bool {
	n, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		if errors.Is(err, strconv.ErrRange) {
			p.errorf(`Unsigned integer value for %s is out of range for a float %d: %s`, name, bitSize, value)
			p.printHelpHint()
			return false
		}
		p.errorf(`Invalid float value for %s: %s`, name, value)
		p.printHelpHint()
		return false
	}
	val.SetFloat(n)
	return true
}
