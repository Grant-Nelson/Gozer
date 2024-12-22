// This package provides a simple way to parse complex CLI arguments for a
// Go application.
//
// It provides methods to populate a struct with values from the command line
// arguments. The exported fields of the struct must be tagged with tags
// that defines how the arguments will be read.
// The tags are used to define groups, flags, and positional arguments.
//
// The tags must follow one of the following patterns:
//
// # Help information
//
// The flags `-h` and `-help` are reserved to print the help message.
// Additional help information can be added to the struct by using the tag
// `arg:"help"` on a string field. The string field must be set to the
// string to show for help.
//
//	type example struct {
//		Help string `arg:"help"`
//	}{
//		Help: "This is the help message to print at the top of the help output",
//	}
//
// If multiple fields are tagged with `arg:"help"`, the fields will be
// concatenated with a newline separator in the order the fields are defined
// in the struct.
//
// The help message is printed automatically when the arguments are invalid
// or when the help flag is used. The help message will be created using
// the struct's fields. The help message will show the help for group being
// used, so it is helpful to define help messages specific to each group.
//
// # Skip a field
//
// Skip a field by using the tag `arg:"skip"`. This is useful for other
// developers to know that the field was intended to not be used.
// A field will also be skipped if it is unexported or if it is not tagged.
//
//	 type example struct {
//		Name string `arg:"skip"`
//	 }
//
// # Define a group
//
// To make a group, use the tag `arg:"group,<name>,<description>"`.
// This is used for grouping arguments similar to Go tools like `go test`.
// Once a group is defined in the arguments, all remaining arguments will
// be handled as part of that group.
//
//   - This is only valid on struct fields or pointers to struct fields.
//     The struct must follow the rules defined in the package.
//   - Groups may be nested. They can also be recursive but it is not
//     recommended for clarity.
//   - A group cannot be required.
//   - The group's name must be unique among other groups in the struct
//     and may not contain spaces, or start with a dash or number.
//     If the name is empty, the field's name is used.
//     Multiple names can be defined with a bar (|) separator, e.g. `t|test`.
//     Group names are case-sensitive.
//   - The group's description may be empty.
//     The description is used for the help output.
//
// The field for a group may be nil if a pointer and a new instance will be
// created for the group if the group is used. If the field is set to an
// instance or not a pointer, the instance will provide the defaults
// for the group and will be set if the group is used. Any group not used
// and is a pointer will be set to nil.
//
//	type testGroup struct {
//		Verbose bool `arg:"flag,v,Print verbose output for tests"`
//	}
//	type example struct {
//		Test *testGroup `arg:"group,t|test,Use to run the test tool"`
//	}
//
// # Define a flag
//
// To define a flag, use the tag `arg:"flag,<name>,<description>"`.
// A flag is used in the command arguments as the flag name preceded by a
// single dash (-) and the value will be the next argument.
// If the field is a bool, then the value maybe left off and the bool will
// be set to true.
//
//   - This is valid on basic type fields, e.g. int, string, bool, etc.
//     The field may be a pointer to a basic type. If the field is a pointer,
//     and the flag is not used, the field will be set to nil.
//   - The flag's name must be unique among other flags in the struct
//     and may not contain spaces, or start with a dash or number.
//     If the name is empty, the field's name is used.
//     Multiple names can be defined with a bar (|) separator, e.g. `i|input`.
//     Group names are case-sensitive.
//   - The flag's description may be empty.
//     The description is used for the help output.
//
// By default a flag is optional and the given value of the field
// will be used as the default value. To make the flag
// required, use the tag `arg:"required,flag,<name>,<description>"`.
//
//	type example struct {
//		Verbose bool `arg:"flag,v,Print verbose output"`
//		Input string `arg:"required,flag,i|input,The input file to use"`
//	}
//
// # Define a positional argument
//
// To define a positional argument (pos), use the tag
// `arg:"pos,<name>,<description>"`.
// The positional arguments are defined in the order they appear
// in the struct. If an argument is not a flag or a group, then it is used as
// a positional argument.
//
//   - This is only valid on basic type fields, e.g. int, string, bool, etc.
//     The field may be a pointer to a basic type. If the field is a pointer,
//     and the optional positional is not given, the field will be set to nil.
//   - The pos's name must be unique among other pos in the struct
//     and may not contain spaces or bars (|), or start with a dash or number.
//     If the name is empty, the field's name is used.
//     This name is used for the help output.
//   - The pos's description may be empty.
//     The description is used for the help output.
//
// By default the positional argument is required, unless variadic that is
// defaulted to optional and may only be optional. To make the positional
// argument optional, use the tag `arg:"optional,pos,<description>"`.
// Once an optional positional argument is used, all following positional
// arguments must also be optional.
//
//	type example struct {
//		Input string `arg:"required,pos,input,The input file to use"`
//		Output string `arg:"optional,pos,output,The optional output file to use"`
//	}
//
// If the last positional argument is a slice, then it is a variadic
// and will consume all remaining positional arguments.
// A variadic positional arguments may not be a pointer nor
// may not have pointer elements. The whole slice will be cleared if no value
// is set or before the first value is set.
//
//	type example struct {
//		Files []string `arg:"pos,files,The files to process"`
//	}
package args

import (
	"io"
	"os"
	"reflect"
)

// Parse parses the command line arguments, i.e. [os.Args], and populates
// the given struct with the values.
//
// The struct must follow the rules defined in the package.
//
// This will panic if the given struct (s) is not valid.
// This returns true if the arguments were parsed successfully,
// otherwise false if the arguments were invalid or help was printed.
func Parse(s any) bool {
	return ParseArgs(s, os.Args, nil, nil)
}

// Parse parses the given arguments and populates the given struct
// with the values.
//
// The given arguments must include the command name.
// The struct must follow the rules defined in the package.
//
// The output writers are used to print the help messages to customers
// if the arguments are invalid or help is requested.
// If `stdOut` is nil then [os.Stdout] will be used.
// If `stdErr` is nil then [os.Stderr] will be used.
//
// This will panic if the given struct (s) is not valid.
// This returns true if the arguments were parsed successfully,
// otherwise false if the arguments were invalid or help was printed.
func ParseArgs(s any, args []string, stdOut, stdErr io.Writer) bool {
	val := getValue(s)
	form := getForm(val.Type())
	return parseArgs(form, val, args, stdOut, stdErr)
}

func getValue(s any) reflect.Value {
	pVal := reflect.ValueOf(s)
	if pVal.Kind() != reflect.Pointer {
		panic(ErrNotStructPointer.with(`got %v`, pVal.Kind()))
	}
	if pVal.IsNil() {
		panic(ErrNilPointer)
	}
	val := pVal.Elem()
	if val.Kind() != reflect.Struct {
		panic(ErrNotStructPointer.with(`got pointer to %v`, val.Kind()))
	}
	return val
}
