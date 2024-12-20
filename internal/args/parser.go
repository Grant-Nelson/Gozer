package args

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func parseArgs(f *form, val reflect.Value, args []string, out io.Writer) bool {
	p := &parserNode{
		cmdPath:    args[0],
		form:       f,
		foundFlags: map[*flagForm]bool{},
		usedGroup:  map[*groupForm]bool{},
		atPos:      0,
		val:        val,
		args:       args[1:],
		out:        out,
	}
	return p.Parse()
}

type parserNode struct {
	cmdPath    string
	form       *form
	foundFlags map[*flagForm]bool
	usedGroup  map[*groupForm]bool
	atPos      int
	val        reflect.Value
	args       []string
	out        io.Writer
}

func (p *parserNode) printf(format string, a ...any) {
	if _, err := fmt.Fprintf(p.out, format, a...); err != nil {
		panic(ErrWriteFailure.with(`%w`, err))
	}
}

func (p *parserNode) takeArg() string {
	arg := p.args[0]
	p.args = p.args[1:]
	return arg
}

func (p *parserNode) isFlag(arg string) bool {
	return len(arg) > 1 &&
		strings.HasPrefix(arg, dash) &&
		!unicode.IsDigit(rune(arg[1]))
}

func (p *parserNode) printHelpHint() {
	p.printf("Use %q to print help.\n", p.cmdPath+` -h`)
}

func (p *parserNode) printHelp() {
	p.printf("Usage of %s:\n\n", p.cmdPath)
	for _, helpField := range p.form.Help {
		helpText := p.val.FieldByIndex(helpField.Index).String()
		if len(helpText) > 0 {
			p.printf("%s\n\n", helpText)
		}
	}
	p.printHelpForGroups()
	p.printHelpForFlags()
	p.printHelpForPos()
}

func (p *parserNode) printHelpForGroups() {
	if len(p.form.Groups) == 0 {
		return
	}

	p.printf("Groups:\n")
	for _, group := range p.form.AllGroups {
		p.printf("\t%s\n", strings.Join(group.Names, nameSep))
		if len(group.Description) > 0 {
			p.printf("\t\t%s\n", group.Description)
		}
	}
}

func (p *parserNode) defaultValue(val reflect.Value) string {
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
	return fmt.Sprintf(`%s`, val.Interface())
}

func (p *parserNode) printHelpForFlags() {
	if len(p.form.Flags) == 0 {
		return
	}

	p.printf("Flags:\n")
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
		p.printf("\t%s%s\n", strings.Join(flag.Names, nameSep), extra)
		if len(flag.Description) > 0 {
			p.printf("\t\t%s\n", flag.Description)
		}
	}
}

func (p *parserNode) printHelpForPos() {
	if len(p.form.Pos) == 0 {
		return
	}

	p.printf("Positional Arguments:\n")
	for i, pos := range p.form.Pos {
		name := pos.Name
		if len(name) == 0 {
			name = fmt.Sprintf("arg%d", i+1)
		}
		extra := ``
		if i < p.form.PosOpAt {
			extra = ` (` + required + `)`
		} else {
			val := p.val.FieldByIndex(pos.Field.Index)
			if defVal := p.defaultValue(val); len(defVal) > 0 {
				extra = ` = ` + defVal
			}
		}
		p.printf("\t%s%s\n", name, extra)
		if len(pos.Description) > 0 {
			p.printf("\t\t%s\n", pos.Description)
		}
	}
}

func (p *parserNode) satisfiedFlags() bool {
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
		p.printf("Missing required flag: %s\n", missingFlags[0])
		return false
	default:
		p.printf("Missing required flags: %s\n", strings.Join(missingFlags, `, `))
		return false
	}
}

func (p *parserNode) satisfiedPos() bool {
	if p.atPos >= p.form.PosOpAt {
		return true
	}
	missingPos := []string{}
	for _, pos := range p.form.Pos[p.atPos:p.form.PosOpAt] {
		missingPos = append(missingPos, pos.Name)
	}
	switch len(missingPos) {
	case 0:
		return true
	case 1:
		p.printf("Missing required positional: %s\n", missingPos[0])
		return false
	default:
		p.printf("Missing required positionals: %s\n", strings.Join(missingPos, `, `))
		return false
	}
}

func (p *parserNode) clearUnset() {
	for _, flag := range p.form.AllFlags {
		if !p.foundFlags[flag] {
			p.clearField(flag.Field)
		}
	}
	for _, group := range p.form.AllGroups {
		if !p.usedGroup[group] {
			p.clearField(group.Field)
		}
	}
	if p.atPos < len(p.form.Pos) {
		for _, pos := range p.form.Pos[p.atPos:] {
			p.clearField(pos.Field)
		}
	}
}

func (p *parserNode) clearField(field reflect.StructField) {
	val := p.val.FieldByIndex(field.Index)
	if val.Kind() == reflect.Ptr {
		val.Set(reflect.Zero(val.Type()))
		return
	}
}

func (p *parserNode) Parse() bool {
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

func (p *parserNode) parseArg() bool {
	arg := p.takeArg()
	if p.isFlag(arg) {
		flagName, _ := strings.CutPrefix(arg, dash)
		return p.parseFlag(flagName)
	}
	if group, ok := p.form.Groups[arg]; ok {
		return p.parseGroup(arg, group)
	}
	return p.parsePos(arg)
}

func (p *parserNode) parseGroup(name string, group *groupForm) bool {
	p.usedGroup[group] = true

	if !p.satisfiedFlags() || !p.satisfiedPos() {
		p.printf("Must fill requirements prior to calling group %q.\n", name)
		p.printHelpHint()
		return false
	}
	p.clearUnset()

	val := p.val.FieldByIndex(group.Field.Index)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}

	p.cmdPath += ` ` + name
	p.form = group.Form
	p.foundFlags = map[*flagForm]bool{}
	p.usedGroup = map[*groupForm]bool{}
	p.atPos = 0
	p.val = val
	return true
}

func (p *parserNode) parseFlag(name string) bool {
	flag, ok := p.form.Flags[name]
	if !ok {
		p.printf("Unknown flag: %s\n", name)
		p.printHelpHint()
		return false
	}
	if p.foundFlags[flag] {
		p.printf("Flag %q already set.\n", name)
		p.printHelpHint()
		return false
	}
	p.foundFlags[flag] = true

	if flag.Field.Type.Kind() == reflect.Bool {
		if len(p.args) > 0 {
			switch p.args[0] {
			case trueStr:
				p.takeArg()
				p.val.FieldByIndex(flag.Field.Index).SetBool(true)
				return true
			case falseStr:
				p.takeArg()
				p.val.FieldByIndex(flag.Field.Index).SetBool(false)
				return true
			}
		}
		p.val.FieldByIndex(flag.Field.Index).SetBool(true)
		return true
	}

	if len(p.args) == 0 {
		p.printf("Flag %q requires a value.\n", name)
		p.printHelpHint()
		return false
	}

	value := p.takeArg()
	if p.isFlag(value) {
		if value == help || value == helpShort {
			p.printHelp()
			return false
		}

		p.printf("Flag %q requires a value.\n", name)
		if flag.Field.Type.Kind() == reflect.String {
			p.printf("If the intended string value starts with a dash, escape the value: -%s %q\n", name, value)
		}
		p.printHelpHint()
		return false
	}

	val := p.val.FieldByIndex(flag.Field.Index)
	return p.setValue(val, name, value)
}

func (p *parserNode) parsePos(value string) bool {
	posCount := len(p.form.Pos)
	if p.atPos >= posCount {
		p.printf("Unexpected positional argument: %s\n", value)
		p.printHelpHint()
		return false
	}

	pos := p.form.Pos[p.atPos]
	val := p.val.FieldByIndex(pos.Field.Index)
	if p.atPos == posCount-1 && p.form.Variadic {
		elem := reflect.New(val.Type().Elem())
		if !p.setValue(elem, pos.Name, value) {
			return false
		}
		val.Set(reflect.Append(val, elem))
		return true
	}

	p.atPos++
	return p.setValue(val, pos.Name, value)
}

func (p *parserNode) setValue(val reflect.Value, name string, value string) bool {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val2 := reflect.New(val.Type().Elem())
			val.Set(val2)
			val = val2
		} else {
			val = val.Elem()
		}
	}

	if unquoted, err := strconv.Unquote(value); err == nil {
		value = unquoted
	}

	switch val.Kind() {
	case reflect.Bool:
		switch value {
		case trueStr:
			val.SetBool(true)
			return true
		case falseStr:
			val.SetBool(false)
			return true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			p.printf("Invalid integer value for %s: %s\n", name, value)
			p.printHelpHint()
			return false
		}
		val.SetInt(n)
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			p.printf("Invalid unsigned integer value for %s: %s\n", name, value)
			p.printHelpHint()
			return false
		}
		val.SetUint(n)
		return true
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(value, 64)
		if err != nil {
			p.printf("Invalid float value for %s: %s\n", name, value)
			p.printHelpHint()
			return false
		}
		val.SetFloat(n)
		return true
	case reflect.String:
		val.SetString(value)
		return true
	}
	p.printf("Unsupported type: %v\n", val.Kind())
	p.printHelpHint()
	return false
}
