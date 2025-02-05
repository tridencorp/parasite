package log

import (
	"fmt"
	"io"
	"os"
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
		InfoWriters:  append(defaultWriters, conf.ErrorWriters...),
		DebugWriters: append(defaultWriters, conf.ErrorWriters...),
		TraceWriters: append(defaultWriters, conf.ErrorWriters...),
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

	for _, writer := range config.InfoWriters {
		writer.WriteString(log)
	}
}
