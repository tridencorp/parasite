package log

import (
	"fmt"
	"os"
)

// This is only basic skeleton. We are working on it, it's ongoing proccess :D
// This logger will be in a separate repository and will be used across all our services.

var file *os.File

// @TODO: Make it atomic so we can only call it once.
func Setup(path string) error {
	tmp, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	file = tmp
	return nil
}

// Write to file and stdout.
func Info(format string, args ...any) {
	data := fmt.Sprintf(format + "\n", args...)
	file.Write([]byte(data))

	fmt.Print(data)
}
