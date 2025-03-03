package log

import (
	"fmt"
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
    size  := "%-" + strconv.Itoa(args[i].(int)) + ftype()
    color := Colors[args[i+1].(int)]
    data  := args[i+2]

    format += fmt.Sprintf(color + size + reset, data)
  }

  return format
}

// Return fmt format symbol based on type.
func ftype() string {
  return "s"
}

// Formatted error log.
func Ferror(args ...any) {
	Error("%s", Format(args...))
}
