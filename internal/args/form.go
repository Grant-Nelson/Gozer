package args

import (
	"reflect"
	"strings"
	"unicode"
)

const (
	separator  = `,`
	nameSep    = `|`
	dash       = `-`
	whiteSpace = " \t\n\r"
	tagName    = `arg`
	required   = `required`
	optional   = `optional`
	helpShort  = `h`
	help       = `help`
	skip       = `skip`
	group      = `group`
	flag       = `flag`
	pos        = `pos`
	indent     = "\t"
)

type (
	groupForm struct {
		Names       []string
		Field       reflect.StructField
		Form        *form
		Description string
	}

	flagForm struct {
		Names       []string
		Field       reflect.StructField
		Description string
		Required    bool
	}

	posForm struct {
		Name        string
		Field       reflect.StructField
		Description string
	}

	form struct {
		Help      []reflect.StructField
		Groups    map[string]*groupForm
		Flags     map[string]*flagForm
		AllGroups []*groupForm
		AllFlags  []*flagForm
		Pos       []*posForm
		PosOpAt   int
		Variadic  bool
	}

	formBuilder struct {
		allForms map[reflect.Type]*form
	}
)

func getForm(t reflect.Type) *form {
	b := &formBuilder{
		allForms: map[reflect.Type]*form{},
	}
	return b.ReadStruct(t)
}

func (b formBuilder) takeFirst(tag string) (string, string) {
	parts := strings.SplitN(tag, separator, 2)
	if len(parts) == 1 {
		return strings.TrimSpace(parts[0]), ``
	}
	return strings.TrimSpace(parts[0]), parts[1]
}

func (b formBuilder) isStructType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct
}

func (b formBuilder) isBasicType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	}
	return false
}

func (b formBuilder) isValidName(name string) bool {
	return len(name) > 0 &&
		!strings.ContainsAny(name, whiteSpace) &&
		!strings.HasPrefix(name, dash) &&
		!unicode.IsDigit(rune(name[0]))
}

func (b *formBuilder) ReadStruct(st reflect.Type) *form {
	if f, ok := b.allForms[st]; ok {
		return f
	}

	if st.Kind() == reflect.Ptr {
		st = st.Elem()
	}
	if st.Kind() != reflect.Struct {
		panic(ErrExpectedStructType.with(`%v`, st.Kind()))
	}

	f := &form{
		Groups: map[string]*groupForm{},
		Flags:  map[string]*flagForm{},
	}
	for i := range st.NumField() {
		b.ReadField(f, st.Field(i))
	}
	return f
}

func (b *formBuilder) ReadField(f *form, field reflect.StructField) {
	if !field.IsExported() {
		return
	}

	need := ``
	tag := field.Tag.Get(tagName)
	switch tag {
	case ``, skip:
		return
	case help:
		if field.Type.Kind() != reflect.String {
			panic(ErrInvalidHelpTag.with(`%v`, field.Type.Kind()))
		}
		f.Help = append(f.Help, field)
		return
	}

	first, tag := b.takeFirst(tag)
	switch first {
	case required, optional:
		need = first
		first, tag = b.takeFirst(tag)
	}

	switch first {
	case group:
		b.ReadGroup(f, need, field, tag)
	case flag:
		b.ReadFlag(f, need, field, tag)
	case pos:
		b.ReadPos(f, need, field, tag)
	default:
		panic(ErrUnknownTag.with(`%q`, first))
	}
}

func (b *formBuilder) ReadGroup(f *form, need string, field reflect.StructField, tag string) {
	if !b.isStructType(field.Type) {
		panic(ErrGroupTagWrongType.with(`%v`, field.Type))
	}
	if need == required {
		panic(ErrGroupTagRequired.with(`%q`, tag))
	}
	namePart, description := b.takeFirst(tag)
	if len(namePart) == 0 {
		namePart = field.Name
	}
	names := strings.Split(namePart, nameSep)
	for i, name := range names {
		names[i] = strings.TrimSpace(name)
	}
	gf := &groupForm{
		Names:       names,
		Field:       field,
		Form:        b.ReadStruct(field.Type),
		Description: description,
	}

	f.AllGroups = append(f.AllGroups, gf)
	for _, name := range names {
		if !b.isValidName(name) {
			panic(ErrInvalidGroupName.with(`%q`, name))
		}
		if _, ok := f.Groups[name]; ok {
			panic(ErrGroupAlreadyExists.with(`%q`, name))
		}
		f.Groups[name] = gf
	}
}

func (b *formBuilder) ReadFlag(f *form, need string, field reflect.StructField, tag string) {
	if !b.isBasicType(field.Type) {
		panic(ErrFlagTagWrongType.with(`%v`, field.Type))
	}
	namePart, description := b.takeFirst(tag)
	if len(namePart) == 0 {
		namePart = field.Name
	}
	names := strings.Split(namePart, nameSep)
	for i, name := range names {
		names[i] = strings.TrimSpace(name)
	}
	ff := &flagForm{
		Names:       names,
		Field:       field,
		Description: description,
		Required:    need == required,
	}

	f.AllFlags = append(f.AllFlags, ff)
	for _, name := range names {
		if !b.isValidName(name) {
			panic(ErrInvalidFlagName.with(`%q`, name))
		}
		if name == help || name == helpShort {
			panic(ErrFlagNameReserved.with(`%q`, name))
		}
		if _, ok := f.Flags[name]; ok {
			panic(ErrFlagAlreadyExists.with(`%q`, name))
		}
		f.Flags[name] = ff
	}
}

func (b *formBuilder) ReadPos(f *form, need string, field reflect.StructField, tag string) {
	variadic := field.Type.Kind() == reflect.Slice
	t := field.Type
	if variadic {
		t = t.Elem()
		if t.Kind() == reflect.Pointer {
			panic(ErrPosTagWrongType.with(`%v`, field.Type))
		}
	}
	if !b.isBasicType(t) {
		panic(ErrPosTagWrongType.with(`%v`, field.Type))
	}
	name, description := b.takeFirst(tag)
	if len(name) == 0 {
		name = field.Name
	}
	if !b.isValidName(name) || strings.Contains(name, nameSep) {
		panic(ErrInvalidPosName.with(`%q`, name))
	}
	for _, pos := range f.Pos {
		if pos.Name == name {
			panic(ErrPosAlreadyExists.with(`%q`, name))
		}
	}
	if f.Variadic {
		panic(ErrVarPosNotLast.with(`%q`, name))
	}
	req := (variadic && need == required) ||
		(!variadic && need != optional)
	if req {
		if variadic {
			panic(ErrVarPosRequired.with(`%q`, name))
		}
		if len(f.Pos) > f.PosOpAt {
			panic(ErrPosRequiredAfterOp.with(`%q`, name))
		}
	}

	pf := &posForm{
		Name:        name,
		Field:       field,
		Description: description,
	}
	if req {
		f.PosOpAt = len(f.Pos) + 1
	}
	f.Variadic = variadic
	f.Pos = append(f.Pos, pf)
}
