package log

import (
	"fmt"
	"reflect"
	"strconv"
)

// Require tuples with size 3:
//  * column size
//  * color
//  * data
//
// If last element is not a tuple (single element), display
// it based on it's type - without any resizing.
func Format(args ...any) string {
	reset  := "\033[0m"
	format := ""

	for i:=0; i < len(args); i+=3 {
		data  := args[i+2]
		kind := reflect.TypeOf(data).Kind()

		size  := "%-" + strconv.Itoa(args[i].(int)) + ftype(kind)
		color := Colors[args[i+1].(int)]

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
