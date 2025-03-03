package log

import (
	"fmt"
	"reflect"
	"strings"
)

// Require tuples with size 2:
//  * size:color
//  * data
//
// ex: Format("20:red", data)
func Format(args ...any) string {
	reset  := "\033[0m"
	format := ""

	for i:=0; i < len(args); i+=2 {
		data  := args[i+1]
		kind := reflect.TypeOf(data).Kind()

		// Split size:color argument.
		subargs := strings.Split(args[i].(string), ":")

		size  := "%-" + subargs[0] + ftype(kind)
		color := Colors[subargs[1]]

		if kind == reflect.Slice || kind  == reflect.Array || kind == reflect.Struct {
			symbol := "%" + ftype(kind)
			format += fmt.Sprintf(color + symbol + reset, data)
			continue
		}

		format += fmt.Sprintf(color + size + reset, data)
	}

	return format
}

// Return fmt format symbol based on type.
func ftype(kind reflect.Kind) string {
	switch kind {
		case reflect.String:
			return "s"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return "d"
		case reflect.Float32, reflect.Float64:
			return "f"
		default:
			return "v"
	}
}

// Formatted error log.
func Ferror(args ...any) {
	Error("%s", Format(args...))
}

// Formatted info log.
func Finfo(args ...any) {
	Info("%s", Format(args...))
}

// Formatted debug log.
func Fdebug(args ...any) {
	Debug("%s", Format(args...))
}

// Formatted trace log.
func Ftrace(args ...any) {
	Trace("%s", Format(args...))
}
