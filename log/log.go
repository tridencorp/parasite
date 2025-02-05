package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// This is only basic skeleton. We are working on it, it's ongoing proccess :D
// This logger will be in a separate repository and will be used across all our services.

var file *os.File
var config Config

type Config struct {
	ErrorWriters []io.StringWriter
	InfoWriters  []io.StringWriter
	DebugWriters []io.StringWriter
	TraceWriters []io.StringWriter
}

// Standard colors.
const (
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	black   = "\033[30m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
	reset   = "\033[0m"
)

// Quick implementation of io.StringWriter for stdout.
type StdoutWriter struct {}

func (stdout StdoutWriter) WriteString(log string) (int, error) {
	return fmt.Println(log)
}

// @TODO: Make it atomic so we can only call it once.
func Configure(conf *Config) error {
	// Default writers for stdout and file(disabled for now).
	defaultWriters := []io.StringWriter{StdoutWriter{}}

	// We are setting package config.
	config = Config{
		// Setting writers for each log level.
		ErrorWriters: append(defaultWriters, conf.ErrorWriters...),
		InfoWriters:  append(defaultWriters, conf.InfoWriters...),
		DebugWriters: append(defaultWriters, conf.DebugWriters...),
		TraceWriters: append(defaultWriters, conf.TraceWriters...),
	}

	// file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// Write to file and stdout.
func Info(format string, args ...any) {
	log := fmt.Sprintf(format, args...)

	prefix := "INFO " + prefix()

	for _, writer := range config.InfoWriters {
		writer.WriteString(magenta + prefix + log + reset)
	}
}

func prefix() string {
	_, file, line, _ := runtime.Caller(2)

	file = filepath.Base(file)
	prefix := fmt.Sprintf("(%s:%d) ", file, line)

	return prefix
}
